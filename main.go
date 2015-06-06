package main

import (
  "net/http"
  "github.com/go-martini/martini"
  "fmt"
  // "./utils"
  "./run"
  "./debug"

  "log"
)

func SayHello(w http.ResponseWriter, req *http.Request) {
    req.ParseForm()
    fmt.Println("username", req.Form["username"])
    w.Write([]byte("Hello"))
}

func die(err error, w http.ResponseWriter) {
	log.Print("### ", err)
	w.WriteHeader(http.StatusInternalServerError)
}


func main() {
  m := martini.Classic()

  //file
  //m.Post("/getfile", run.GetFileHandle)
  //m.Post("/Savefile", run.SaveFileHandle)


  //run
  m.Post("/compile", run.CompileHandle)
  //m.Post("/run", run.RunHandle)
  //Debug
  m.Post("/debug/send", debug.SendHand)
  m.Get("/debug/socket", debug.MessageHandle)


  http.ListenAndServe("0.0.0.0:8080", m)
}
