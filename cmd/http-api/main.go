package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/resend/resend-go/v2"
)

type ContactPayload struct {
	Name    string `json:"name"`
	Contact string `json:"contact"`
	Message string `json:"message"`
}

func handleContact(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload ContactPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if payload.Name == "" || payload.Contact == "" || payload.Message == "" {
		http.Error(w, "Missing mandatory fields", http.StatusBadRequest)
		return
	}

	apiKey := os.Getenv("RESEND_API_KEY")
	if apiKey == "" {
		log.Println("[💥 ERROR] RESEND_API_KEY environment variable is empty!")
		http.Error(w, "Internal server configuration error", http.StatusInternalServerError)
		return
	}
	client := resend.NewClient(apiKey)

	htmlContent := fmt.Sprintf(`
		<h3>🔋 Novo Lead do Portfólio SSH/TCP</h3>
		<p><strong>Nome do Operador:</strong> %s</p>
		<p><strong>Contato/Up-link:</strong> %s</p>
		<p><strong>Mensagem Disparada:</strong></p>
		<blockquote style="background: #f4f4f4; padding: 10px; border-left: 3px solid #00ffff;">
			%s
		</blockquote>
	`, payload.Name, payload.Contact, payload.Message)

	params := &resend.SendEmailRequest{
		From:    "Portfolio SSH <onboarding@resend.dev>",
		To:      []string{"daniel.cmdline@humanoid.net"},
		Subject: fmt.Sprintf("📡 [COM-LINK] %s enviou uma mensagem!", payload.Name),
		Html:    htmlContent,
	}

	_, err = client.Emails.SendWithContext(context.Background(), params)
	if err != nil {
		log.Printf("[💥 ERROR] Failed to send email via Resend: %v\n", err)
		http.Error(w, fmt.Sprintf("Failed to dispatch email: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("[✓ SUCCESS] Email sent from: %s (%s)\n", payload.Name, payload.Contact)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "dispatched"}`))
}

func main() {
	// Força o godotenv a buscar o arquivo subindo duas pastas a partir do escopo da api, 
	// ou tentando ler a raiz do projeto de forma explícita.
	err := godotenv.Load(".env", "../.env", "../../.env")
	if err != nil {
		log.Println("[⚠️ WARNING] Não foi possível encontrar o arquivo .env, buscando variáveis locais do sistema...")
	}

	http.HandleFunc("/api/contact", handleContact)

	port := ":8080"
	fmt.Printf("Iniciando API HTTP local na porta %s (Aguardando pacotes do Servidor TCP)...\n", port)
	
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Erro ao iniciar servidor HTTP: %v", err)
	}
}
