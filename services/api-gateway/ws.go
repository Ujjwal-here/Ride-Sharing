package main

import (
	"log"
	"net/http"
	"ride-sharing/shared/contracts"
	"ride-sharing/shared/util"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleDriversWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Printf("Websocket upgrade failed: %v", err)
		return
	}

	defer conn.Close()

	userId := r.URL.Query().Get("userId")

	if userId == "" {
		log.Println("No user ID provided")
		return
	}

	packageSlug := r.URL.Query().Get("packageSlug")

	if packageSlug == "" {
		log.Println("No packageSlug provided")
		return
	}

	type Driver struct {
		Id             string `json:"id"`
		Name           string `json:"name"`
		ProfilePicture string `json:"profile_picture"`
		CarPlate       string `json:"car_plate"`
		PackageSlug    string `json:"package_slug"`
	}

	msg := contracts.WSMessage{
		Type: "driver.cmd.register",
		Data: Driver{
			Id:             userId,
			Name:           "Ujjwal",
			ProfilePicture: util.GetRandomAvatar(1),
			CarPlate:       "ABC123",
			PackageSlug:    packageSlug,
		},
	}

	if err := conn.WriteJSON(msg); err != nil {
		log.Printf("Error sending message: %v", err)
		return
	}

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}

		log.Printf("Received message: %s", message)
	}
}

func handleRidersWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Printf("Websocket upgrade failed: %v", err)
	}

	defer conn.Close()

	userId := r.URL.Query().Get("userId")

	if userId == "" {
		log.Printf("No user ID provided")
		return
	}

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}

		log.Printf("Received message: %s", message)
	}
}
