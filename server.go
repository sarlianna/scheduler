package main

import (
  "fmt"
  "encoding/json"
  "html/template"
  "net/http"
  "scheduler"
 )

const path = "/api/"
const lenPath = len(path)

func writeJson( w http.ResponseWriter, s *interface{} ) bool {
  
  enc := NewEncoder(w)
  e := enc.Encode(s)
  if e != nil {
    //error template
  }

}

func ApiHandler(w http.ResponseWriter, r *http.Request) {
  
  resourcePath := r.URL.Path[lenPath:]
  
}

func main() {

  http.HandleFunc(path, ApiHandler)
  http.ListenAndServe(":8080", nil)
}
