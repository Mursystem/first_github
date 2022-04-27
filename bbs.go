package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

var aport string = os.Getenv("PORT")

// We'll need to define an Upgrader
// this will require a Read and Write buffer size
type clientType struct {
	connection *websocket.Conn
	id         int
}

var allClientsSlice []clientType
var instanceId int = 0
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// define a reader which will listen for
// new messages being sent to our WebSocket
// endpoint
func reader(thisClient clientType) {
	for {
		// read in a message
		mesgeType, p, err := thisClient.connection.ReadMessage()
		if err != nil {
			log.Println(err)
			log.Println(thisClient.id)
			//remove the disconected user
			thisIndex := -1
			var removeIndex int
			for _, oneClient := range allClientsSlice {
				thisIndex++
				if oneClient.id == thisClient.id {
					println(thisIndex)
					removeIndex = thisIndex
				}
			}
			allClientsSlice = RemoveIndex(allClientsSlice, removeIndex)

			return
		}
		// print out that message for clarity
		log.Println(mesgeType)
		log.Println(string(p))
		//if err := conn.WriteMessage(messageType, p); err != nil {
		if err := thisClient.connection.WriteMessage(1, []byte("Recieved")); err != nil {
			log.Println(err)
			return
		}

		for _, oneClient := range allClientsSlice {
			fmt.Println(oneClient.id)
			oneClient.connection.WriteMessage(1, []byte("A change was made by a user"))
		}

	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home Page")
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	// upgrade this connection to a WebSocket
	// connection
	instanceId++
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	//add new user with id and connection
	var newClient clientType
	newClient.connection = ws
	newClient.id = instanceId

	allClientsSlice = append(allClientsSlice, newClient)

	log.Println("Client Connected")

	err = ws.WriteMessage(1, []byte("Hi Client!"))
	if err != nil {
		log.Println(err)
	}
	// listen indefinitely for new messages coming
	// through on our WebSocket connection
	reader(newClient)
}

func setupRoutes() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/ws", wsEndpoint)
}

func RemoveIndex(s []clientType, index int) []clientType {
	ret := make([]clientType, 0)
	ret = append(ret, s[:index]...)
	return append(ret, s[index+1:]...)
}

func main() {
	fmt.Println("Hello World")
	setupRoutes()
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", aport), nil))

}
