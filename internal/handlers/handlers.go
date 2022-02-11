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
	Action         string    `json:"action"`
	Message        string    `json:"message"'`
	Hand           game.Hand `json:"hand"`
	Table          string    `json:"table"`
	Deck           bool      `json:"deck"`
	Trash          bool      `json:"trash"`
	ConnectedUsers []string  `json:"connected_users"`
}

type WsPayload struct {
	Action   string              `json:"action"`
	Username string              `json:"username"`
	Message  string              `json:"message"`
	Conn     WebSocketConnection `json:"-"`
}

type ClientData struct {
	Username string
	Hand     game.Hand
}

var (
	clients           = make(map[WebSocketConnection]*ClientData)
	wsChan            = make(chan WsPayload)
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
		Username: "",
		Hand:     make(game.Hand, 0, 18),
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
	response := WsJsonResponse{
		Deck: true,
	}
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
			response.Table = ""
			for client := range clients {
				hand := clients[client].Hand
				for n := 6 - len(hand); n > 0; n-- {
					if deck.Length != 0 {
						hand = game.BringToHand(hand, deck.Get())
					} else {
						break
					}
				}
				response.Hand = hand
				broadcastToClient(client, response)
			}
		case "break":
			response.Action = "hand"
			response.Table = ""
			for client := range clients {
				hand := clients[client].Hand
				for n := 6 - len(hand); n > 0; n-- {
					game.BringToHand(hand, deck.Get())
				}
				response.Hand = hand
				broadcastToClient(client, response)
			}
		case "broadcast":
			response.Action = "broadcast"
			response.Message = fmt.Sprintf("<strong>%s</strong>: %s", e.Username, e.Message)
			broadcastToAll(response)
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
