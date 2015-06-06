package run

import (
  "net/http"
  "fmt"
  "os"
  "os/exec"
  "log"
  "path/filepath"
  "strings"


)


func Compile(code string, fileName string ) string {
  path, _ := os.Getwd()
  filePath := filepath.Clean(path + "/tmp/" + fileName)
  fout, _ := os.Create(filePath)
  lines := strings.Split(code, "\\n")

  for _, line := range lines{
    fout.WriteString(line+"\n")
    }

  if err := fout.Close(); nil != err {
    log.Printf("err: %s", err)
    }

  cmd := exec.Command("g++", "-O2", "-o", "test/test1", filePath)
  out, _ := cmd.CombinedOutput()
  return string(out)
  }





func CompileHandle(w http.ResponseWriter, req *http.Request) {
    req.ParseForm()
    fmt.Println("code", req.Form["code"])
    code := req.Form["code"][0]
    fileName := req.Form["filename"][0]
    out := Compile(code, fileName)
    w.Write([]byte(out))

}
