package handlers

import (
	"fmt"
	"github.com/dragun-igor/fool/internal/game"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sort"
)

type WebSocketConnection struct {
	*websocket.Conn
}

type WsJsonResponse struct {
	Action         string          `json:"action"`
	Message        string          `json:"message"`
	Hand           []game.CardItem `json:"hand"`
	Table          game.Table      `json:"table"`
	Deck           bool            `json:"deck"`
	Trash          bool            `json:"trash"`
	ConnectedUsers []string        `json:"connected_users"`
}

type WsPayload struct {
	Action   string              `json:"action"`
	Username string              `json:"username"`
	Message  string              `json:"message"`
	ID       int                 `json:"id"`
	Card     game.CardItem       `json:"card"`
	Conn     WebSocketConnection `json:"-"`
}

type ClientData struct {
	Username     string
	Hand         map[int]*game.CardItem
	SelectedCard int
}

var sortMap = map[string]int{
	"six":   1,
	"seven": 2,
	"eight": 3,
	"nine":  4,
	"ten":   5,
	"jack":  6,
	"queen": 7,
	"king":  8,
	"ace":   9,
}

var (
	clients = make(map[WebSocketConnection]*ClientData, 6)
	wsChan  = make(chan WsPayload)
	table   = game.Table{
		Pairs: make([]game.Pair, 0, 6),
	}
	upgradeConnection = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
)

func WsEndpoint(w http.ResponseWriter, r *http.Request) {
	ws, err := upgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	conn := WebSocketConnection{ws}
	hand := &ClientData{
		Username:     "",
		Hand:         make(map[int]*game.CardItem, 36),
		SelectedCard: 0,
	}
	clients[conn] = hand

	go ListenForWs(&conn)
}

func ListenForWs(conn *WebSocketConnection) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("err", fmt.Sprintf("%v", r))
		}
	}()

	payload := WsPayload{}

	for {
		err := conn.ReadJSON(&payload)
		if err != nil {
			// do nothing
		} else {
			payload.Conn = *conn
			wsChan <- payload
		}
	}
}

func ListenToWsChannel() {
	deck := game.NewDeck()
	response := WsJsonResponse{}
	for {
		e := <-wsChan
		switch e.Action {
		case "users":
			response = getUserList(response)
			broadcastToAll(response)
		case "add_user":
			clients[e.Conn].Username = e.Username
			response = getUserList(response)
			broadcastToAll(response)
		case "delete_user":
			delete(clients, e.Conn)
			response = getUserList(response)
			broadcastToAll(response)
		case "start_game":
			response.Action = "hand"
			for client := range clients {
				hand := clients[client].Hand
				for i := 0; i < 6; i++ {
					if deck.Length != 0 {
						card := deck.Get()
						hand[card.ID] = &card
					} else {
						break
					}
				}
				handSlice := make([]game.CardItem, 0, len(hand))
				for _, v := range hand {
					handSlice = append(handSlice, *v)
				}
				response.Hand = handSlice
				sort.Slice(handSlice, func(i, j int) bool {
					if handSlice[i].Suit < handSlice[j].Suit {
						return true
					}
					if handSlice[i].Suit > handSlice[j].Suit {
						return false
					}
					if sortMap[handSlice[i].Denomination] < sortMap[handSlice[j].Denomination] {
						return true
					}
					return false
				})
				broadcastToClient(client, response)
			}
		case "select_card":
			response.Action = "hand"
			client := clients[e.Conn]
			if client.SelectedCard > 0 {
				client.Hand[client.SelectedCard].Selected = false
			}
			if e.ID != client.SelectedCard {
				client.Hand[e.ID].Selected = true
				client.SelectedCard = e.ID
			} else {
				client.SelectedCard = 0
			}
			handSlice := make([]game.CardItem, 0, len(client.Hand))
			for _, v := range client.Hand {
				handSlice = append(handSlice, *v)
			}
			response.Hand = handSlice
			sort.Slice(handSlice, func(i, j int) bool {
				if handSlice[i].Suit < handSlice[j].Suit {
					return true
				}
				if handSlice[i].Suit > handSlice[j].Suit {
					return false
				}
				if sortMap[handSlice[i].Denomination] < sortMap[handSlice[j].Denomination] {
					return true
				}
				return false
			})
			broadcastToClient(e.Conn, response)
		case "put_on_table":
			response.Action = "table"
			client := clients[e.Conn]
			if client.SelectedCard <= 0 {
				break
			}
			table.PutCard(*client.Hand[client.SelectedCard])
			response.Table = table
			fmt.Println(table)
			fmt.Println(response.Table)
			broadcastToClient(e.Conn, response)
		case "remove_card":
			response.Action = "hand"
			client := clients[e.Conn]
			delete(client.Hand, client.SelectedCard)
			handSlice := make([]game.CardItem, 0, len(client.Hand))
			for _, v := range client.Hand {
				handSlice = append(handSlice, *v)
			}
			response.Hand = handSlice
			sort.Slice(handSlice, func(i, j int) bool {
				if handSlice[i].Suit < handSlice[j].Suit {
					return true
				}
				if handSlice[i].Suit > handSlice[j].Suit {
					return false
				}
				if sortMap[handSlice[i].Denomination] < sortMap[handSlice[j].Denomination] {
					return true
				}
				return false
			})
			response.Hand = handSlice
			client.SelectedCard = 0
			broadcastToClient(e.Conn, response)
		}
	}
}

func getUserList(response WsJsonResponse) WsJsonResponse {
	var userList []string
	for key := range clients {
		if clients[key].Username != "" {
			userList = append(userList, clients[key].Username)
		}
	}
	sort.Strings(userList)
	response.Action = "list_users"
	response.ConnectedUsers = userList
	return response
}

func broadcastToClient(client WebSocketConnection, response WsJsonResponse) {
	err := client.WriteJSON(response)
	if err != nil {
		log.Println("websocket error")
		_ = client.Close()
		delete(clients, client)
	}
}

func broadcastToAll(response WsJsonResponse) {
	for client := range clients {
		broadcastToClient(client, response)
	}
}

func renderPage(w http.ResponseWriter, r *http.Request, page string) error {
	http.ServeFile(w, r, page)
	return nil
}

// Home renders the home page
func Home(w http.ResponseWriter, r *http.Request) {
	err := renderPage(w, r, "./html/home.html")
	if err != nil {
		log.Println(err)
	}
}
