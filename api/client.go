package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// SendContact envia os dados capturados para a nossa API HTTP local na porta 8080.
func SendContact(name, contact, message string) error {
	payload := map[string]string{
		"name":    name,
		"contact": contact,
		"message": message,
	}

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	endpointURL := "http://127.0.0.1:8080/api/contact"

	resp, err := http.Post(endpointURL, "application/json", bytes.NewBuffer(jsonBytes))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API HTTP retornou erro: %s", resp.Status)
	}

	return nil
}