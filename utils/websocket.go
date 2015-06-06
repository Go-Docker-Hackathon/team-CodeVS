package utils

//
//
// // terminal WebSocket
// m.Get("/debug/term", websocket.Handler(func(ws *websocket.Conn) {
// 	// copy GDB to WS and WS to GDB
// 	go io.Copy(gdb, ws)
// 	io.Copy(ws, gdb)
// }))
//
// // send command action
// m.Get("/debug/command", func(w http.ResponseWriter, req *http.Request) {
// 	if req.Method != "POST" {
// 		log.Print("### invalid method")
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}
// 	data, err := ioutil.ReadAll(req.Body)
// 	if err != nil {
// 		die(err, w)
// 		return
// 	}
// 	log.Print("<<< ", string(data))
// 	command := []string{}
// 	err = json.Unmarshal(data, &command)
// 	if err != nil {
// 		die(err, w)
// 		return
// 	}
// 	result, err := gdb.Send(command[0], command[1:]...)
// 	if err != nil {
// 		die(err, w)
// 		return
// 	}
// 	reply, err := json.Marshal(result)
// 	if err != nil {
// 		die(err, w)
// 		return
// 	}
// 	io.WriteString(w, string(reply))
// 	log.Print(">>> ", string(reply))
// })
// "encoding/json"
// "github.com/cyrus-and/gdb"
// "golang.org/x/net/websocket"
// "io"
// "io/ioutil"
// "log"
// "net/http"

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
var _, _ = gdb.New(func(gdbMessages map[string]interface{}) {
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

func WebSocketHandle(w http.ResponseWriter, r *http.Request) {
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


// package utils
//
// import (
//   "net/http"
//   "log"
//   "net"
//   "sync"
//   "github.com/gorilla/websocket"
//
// )
//
//
// var ActiveClients = make(map[ClientConn]int)
// var ActiveClientsRWMutex sync.RWMutex
//
// type ClientConn struct {
//   websocket *websocket.Conn
//   clientIP  net.Addr
// }
//
//
//
// func addClient(cc ClientConn) {
// 	ActiveClientsRWMutex.Lock()
// 	ActiveClients[cc] = 0
// 	ActiveClientsRWMutex.Unlock()
// }
//
// func deleteClient(cc ClientConn) {
// 	ActiveClientsRWMutex.Lock()
// 	delete(ActiveClients, cc)
// 	ActiveClientsRWMutex.Unlock()
// }
//
// func broadcastMessage(messageType int, message []byte) {
// 	ActiveClientsRWMutex.RLock()
// 	defer ActiveClientsRWMutex.RUnlock()
//
// 	for client, _ := range ActiveClients {
// 		if err := client.websocket.WriteMessage(messageType, message); err != nil {
// 			return
// 		}
// 	}
// }
//
// func WebSocketHandle(w http.ResponseWriter, r *http.Request) {
//   log.Println(ActiveClients)
//   ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
//   if _, ok := err.(websocket.HandshakeError); ok {
//     http.Error(w, "Not a websocket handshake", 400)
//     return
//   } else if err != nil {
//     log.Println(err)
//     return
//   }
//   client := ws.RemoteAddr()
//   sockCli := ClientConn{ws, client}
//   addClient(sockCli)
//
//   for {
//     log.Println(len(ActiveClients), ActiveClients)
//     messageType, p, err := ws.ReadMessage()
//     log.Println(p)
//     if err != nil {
//       deleteClient(sockCli)
//       log.Println("bye")
//       log.Println(err)
//       return
//     }
//     t := "aaaaa"
//     broadcastMessage(messageType, []byte(t))
//   }
// }
