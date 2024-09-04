package main

import (
	"encoding/json"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
)

type IPRes struct {
	IP      net.IP
	Port    string
	Forward string
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// A simple healthcheck handler

	slog.Info("Receive call on healthcheck handler")

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	io.WriteString(w, `{"alive": true}`)
}

func GetPubIPHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("Receive call on Public IP handler")
	// Split address into ip and port
	ip, port, err := net.SplitHostPort(r.RemoteAddr)

	if err != nil {
		http.Error(w, "Can't get address.", http.StatusBadRequest)
		slog.Error("Error during address splitting", "msg", err.Error())
		return
	}
	// Parse IP
	userIP := net.ParseIP(ip)
	// need to check if userIP is nill
	if userIP == nil {
		http.Error(w, "Unsupported address format.", http.StatusBadRequest)
		slog.Error("Wrong address format")
		return
	}

	// Check if request was forwarded
	forward := r.Header.Get("X-Forwarded-For")

	res := IPRes{
		IP:      userIP,
		Port:    port,
		Forward: forward,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func main() {
	// Setup structured logger
	jsonHandler := slog.NewJSONHandler(os.Stderr, nil)
	logger := slog.New(jsonHandler)
	slog.SetDefault(logger)

	http.HandleFunc("/ip", GetPubIPHandler)
	http.HandleFunc("/healthz", HealthCheckHandler)

	slog.Info("Server is Running ...")
	http.ListenAndServe(":9090", nil)
}
