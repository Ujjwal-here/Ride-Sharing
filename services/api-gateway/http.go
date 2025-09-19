package main

import (
	"encoding/json"
	"net/http"
	"ride-sharing/shared/contracts"
)

func handleTripPreview(w http.ResponseWriter, r *http.Request) {
	var reqBody previewTripRequest

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Failed to parse JSON data", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	response := contracts.APIResponse{Data: "ok"}
	writeJSON(w, http.StatusCreated, response)
}
