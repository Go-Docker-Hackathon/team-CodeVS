package main

import (
  "net/http"
  "github.com/go-martini/martini"
  "./run"
  "./debug"

)



func main() {
  m := martini.Classic()

  //file
  //m.Post("/getfile", run.GetFileHandle)
  //m.Post("/Savefile", run.SaveFileHandle)


  //run
  m.Post("/compile", run.CompileHandle)
  m.Post("/run", run.RunHandle)
  //Debug
  m.Post("/debug/send", debug.SendHand)
  m.Get("/debug/socket", debug.MessageHandle)


  http.ListenAndServe("0.0.0.0:8080", m)
}
