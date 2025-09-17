package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `{"status": "healthy"}`)
}

// PaymentAdd get new payment request
func PaymentAdd(w http.ResponseWriter, r *http.Request) {
	// TODO: move to payments folder, rename to "Add"
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var data map[string]any

	body, _ := io.ReadAll(r.Body)
	err := json.Unmarshal(body, &data)
	fmt.Println(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
