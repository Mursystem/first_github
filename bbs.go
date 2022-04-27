package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

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
		_, p, err := thisClient.connection.ReadMessage()
		p = append(p, 13)
		if err != nil {
			log.Println(err)
			//log.Println(thisClient.id)
			//remove the disconected user
			thisIndex := -1
			var removeIndex int
			for _, oneClient := range allClientsSlice {
				thisIndex++
				if oneClient.id == thisClient.id {
					//println(thisIndex)
					removeIndex = thisIndex
				}
			}
			allClientsSlice = RemoveIndex(allClientsSlice, removeIndex)
			return
		}
		// print out that message for clarity

		//log.Println(string(p))
		//if err := conn.WriteMessage(messageType, p); err != nil {
		if err := thisClient.connection.WriteMessage(1, []byte("✔")); err != nil {
			log.Println(err)
			return
		}

		for _, oneClient := range allClientsSlice {
			//print all contected client id's
			fmt.Println(oneClient.id)
			oneClient.connection.WriteMessage(1, ([]byte(p)))
		}

	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
	fmt.Fprintf(w, "App running and serving wss")
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
	connectStringRdy := "Welcome, your id is " + strconv.Itoa(newClient.id) + "\n"
	err = newClient.connection.WriteMessage(1, []byte(connectStringRdy))
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
	if aport == "" {
		aport = "5000"
	}
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", aport), nil))

}
