package main

import (
  "net/http"
  "fmt"
  "encoding/json"
  "net"
)

type IPRes struct {
  IP net.IP
  Port string
  Forward string
}

func ip(w http.ResponseWriter, r *http.Request){
  // Split address into ip and port
  ip, port, err := net.SplitHostPort(r.RemoteAddr)

  if err != nil {
    fmt.Println("Error during address splitting", err.Error())
  }
  // Parse IP
  userIP := net.ParseIP(ip)
  // need to check if userIP is nill
  if userIP == nil {
    fmt.Println("Wrong address format")
  }

  // Check if request was forwarded
  forward := r.Header.Get("X-Forwarded-For")

  res := IPRes{
    IP: userIP,
    Port: port,
    Forward: forward,
  }

  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusOK)
  json.NewEncoder(w).Encode(res)
}

func main() {
  http.HandleFunc("/ip", ip)

  fmt.Println("Server Running...")
  http.ListenAndServe(":9090", nil)
}
