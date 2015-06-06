package debug



import (
  "github.com/cyrus-and/gdb"
  "encoding/json"
  "net/http"
  "log"
  "net"
  "sync"
  "github.com/gorilla/websocket"

)


var messages = make(chan []byte)
var Gdb, _ = gdb.New(func(gdbMessages map[string]interface{}) {
	gdbMessagesText, err := json.Marshal(gdbMessages)
	if err != nil {
		log.Fatal(err)
	}
	messages <- gdbMessagesText
})




var ActiveClients = make(map[ClientConn]int)
var ActiveClientsRWMutex sync.RWMutex

type ClientConn struct {
  websocket *websocket.Conn
  clientIP  net.Addr
}

func die(err error, w http.ResponseWriter) {
	log.Print("### ", err)
	w.WriteHeader(http.StatusInternalServerError)
}

func addClient(cc ClientConn) {
	ActiveClientsRWMutex.Lock()
	ActiveClients[cc] = 0
	ActiveClientsRWMutex.Unlock()
}

func deleteClient(cc ClientConn) {
	ActiveClientsRWMutex.Lock()
	delete(ActiveClients, cc)
	ActiveClientsRWMutex.Unlock()
}

func broadcastMessage(messageType int, message []byte) {
	ActiveClientsRWMutex.RLock()
	defer ActiveClientsRWMutex.RUnlock()

	for client, _ := range ActiveClients {
		if err := client.websocket.WriteMessage(messageType, message); err != nil {
			return
		}
	}
}

func MessageHandle(w http.ResponseWriter, r *http.Request) {
  log.Println(ActiveClients)
  ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
  if _, ok := err.(websocket.HandshakeError); ok {
    http.Error(w, "Not a websocket handshake", 400)
    return
  } else if err != nil {
    log.Println(err)
    return
  }
  client := ws.RemoteAddr()
  sockCli := ClientConn{ws, client}
  addClient(sockCli)
for message := range messages {
    broadcastMessage(1, message)

}

}

func SendHand(w http.ResponseWriter, req *http.Request) {
    req.ParseForm()
	exec := req.Form["exec"][0]


	result, err := Gdb.Send(exec)
	if err != nil {
		die(err, w)
		return
	}
	reply, err := json.Marshal(result)
	if err != nil {
		die(err, w)
		return
	}

    w.Write([]byte(string(reply)))
}
