package functions

import (
	"encoding/json"
	"net/http"
	"strings"

	"portfolio-ssh/utils"
)

const contactBodyLimit = 8 << 10

// ContactPayload define a estrutura do JSON para o endpoint de contato.
type ContactPayload struct {
	Name    string `json:"name"`
	Contact string `json:"contact"`
	Message string `json:"message"`
}

// HandleContact é o handler HTTP para o formulário de contato.
func HandleContact(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, contactBodyLimit)
	defer r.Body.Close()

	var payload ContactPayload
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&payload); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	payload.Name = strings.TrimSpace(payload.Name)
	payload.Contact = strings.TrimSpace(payload.Contact)
	payload.Message = strings.TrimSpace(payload.Message)

	if payload.Name == "" || payload.Contact == "" || payload.Message == "" {
		http.Error(w, "Missing mandatory fields", http.StatusBadRequest)
		return
	}

	// FIX: Aponta para a nova função pública dentro de utils
	err := utils.SendEmail(payload.Name, payload.Contact, payload.Message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "dispatched"}`))
}