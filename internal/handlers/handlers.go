package handlers

import (
	"fmt"
	"github.com/dragun-igor/fool/internal/game"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sort"
	"time"
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
	Username string
	HandData game.Hand
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
	deck  *game.Deck
	table *game.Table

	clients           = make(map[WebSocketConnection]*ClientData, 6)
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
	clients[conn] = &ClientData{}

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
	response := WsJsonResponse{}
	for {
		var err error = nil
		e := <-wsChan
		switch e.Action {
		case "users":
			response = getUserList(response)
			broadcastToAll(response)
		case "add_user":
			if len(clients) > 6 {
				break
			}
			clients[e.Conn].Username = e.Username
			response = getUserList(response)
			broadcastToAll(response)
		case "delete_user":
			delete(clients, e.Conn)
			response = getUserList(response)
			broadcastToAll(response)
		case "start_game":
			deck, table = game.StartNewGame()
			response.Action = "hand"
			for conn := range clients {
				client := clients[conn]
				client.HandData, err = deck.GetHand()
				if err != nil {
					log.Println(err)
				}
				response.Hand = handMapToSortedSlice(client.HandData.Hand)
				broadcastToClient(conn, response)
			}
		case "select_card":
			client := clients[e.Conn]
			client.HandData, err = table.SelectCard(e.ID, client.HandData)
			if err != nil {
				log.Println(err)
			}
			//client.HandData, err = table.HelperCardCanPut(client.HandData)
			response.Hand = handMapToSortedSlice(client.HandData.Hand)
			response.Table = *table

			response.Action = "hand"
			broadcastToClient(e.Conn, response)

			time.Sleep(10 * time.Millisecond)

			response.Action = "table"
			broadcastToClient(e.Conn, response)
		case "put_on_table":
			client := clients[e.Conn]
			client.HandData, err = table.PutCardOnTable(client.HandData)
			if err != nil {
				log.Println(err)
			}
			//client.HandData, err = table.HelperCardCanPut(client.HandData)
			response.Hand = handMapToSortedSlice(client.HandData.Hand)
			response.Table = *table

			response.Action = "table"
			broadcastToAll(response)

			time.Sleep(10 * time.Millisecond)

			response.Action = "hand"
			broadcastToClient(e.Conn, response)
		}
	}
}

func handMapToSortedSlice(hand map[int]game.CardItem) []game.CardItem {
	res := make([]game.CardItem, 0, len(hand))
	for _, v := range hand {
		res = append(res, v)
	}
	sort.Slice(res, func(i, j int) bool {
		if res[i].TrumpSuit && !res[j].TrumpSuit {
			return true
		}
		if !res[i].TrumpSuit && res[j].TrumpSuit {
			return false
		}
		if res[i].Suit < res[j].Suit {
			return true
		}
		if res[i].Suit > res[j].Suit {
			return false
		}
		if sortMap[res[i].Denomination] < sortMap[res[j].Denomination] {
			return true
		}
		return false
	})
	return res
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
