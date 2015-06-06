package run
import (
  "net/http"
  "fmt"
  "os"
  "os/exec"
  "path/filepath"
  )



func Run(fileName string ) string {
    path, _ := os.Getwd()
    runPath := filepath.Clean(path + "/tmp/compile/" + fileName)
    cmd := exec.Command(runPath)
    fmt.Println(runPath)
    out, _ := cmd.CombinedOutput()
    return string(out)
}



func RunHandle(w http.ResponseWriter, req *http.Request) {
      req.ParseForm()
      fmt.Println("code", req.Form["code"])
      code := req.Form["code"][0]
      fileName := req.Form["filename"][0]
      out := Compile(code, fileName)

      if  len(out) ==0  {
        out = Run(fileName)
        }

      w.Write([]byte(out))

}
