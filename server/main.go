package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/livekit/protocol/auth"
)

var (
	livekitUrl = "wss://tentwenty-bp8gb2jg.livekit.cloud"
	apiKey     = "API4vY5fJ6zxS6e"
	apiSecret  = "GGqV7dBkfi4mtBwK1UD1EvJCLRQCouB7YcDSwyR07MR"
)

func getJoinToken(apiKey, apiSecret, room, identity string) (string, error) {
	at := auth.NewAccessToken(apiKey, apiSecret)
	grant := &auth.VideoGrant{
		RoomJoin: true,
		Room:     room,
	}
	at.AddGrant(grant).SetIdentity(identity).SetValidFor(time.Hour)

	return at.ToJWT()
}

func tokenHandler(w http.ResponseWriter, r *http.Request) {
	room := r.URL.Query().Get("roomName")
	identity := r.URL.Query().Get("identity")

	token, err := getJoinToken(apiKey, apiSecret, room, identity)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		log.Printf("Error generating token: %v\n", err)
		return
	}

	fmt.Println("room: ", room)
	fmt.Println("identity: ", identity)
	fmt.Println("token: ", token)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"accessToken": token})
}

func main() {
	http.HandleFunc("/api/get_token", tokenHandler)
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
