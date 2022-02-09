package handlers

import (
	"fmt"
	"github.com/dragun-igor/fool_card_game/internal/game"
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
	Hand           game.Hand `json:"hand"`
	Table          string    `json:"table"`
	Deck           bool      `json:"deck"`
	Trash          bool      `json:"trash"`
	ConnectedUsers []string  `json:"connected_users"`
}

type WsPayload struct {
	Action   string              `json:"action"`
	Username string              `json:"username"`
	Message  interface{}         `json:"message"`
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
		case "username":
			clients[e.Conn].Username = e.Username
			users := getUserList()
			response.Action = "list_users"
			response.ConnectedUsers = users
			broadcastToAll(response)
		case "left":
			delete(clients, e.Conn)
			users := getUserList()
			response.Action = "list_users"
			response.ConnectedUsers = users
			broadcastToAll(response)
		case "start_game":
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
		}
	}
}

func getUserList() []string {
	var userList []string
	for key := range clients {
		if clients[key].Username != "" {
			userList = append(userList, clients[key].Username)
		}
	}
	sort.Strings(userList)
	return userList
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
