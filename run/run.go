package run
import (
  "net/http"
  "fmt"
  )


func SayHello(w http.ResponseWriter, req *http.Request, t string) {
    fmt.Println(t)
    w.Write([]byte("Hello"))
}
