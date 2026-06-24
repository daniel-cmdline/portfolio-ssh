package utils

import (
	"context"
	"fmt"
	"html"
	"log"
	"os"
	"strings"
	"time"

	"github.com/resend/resend-go/v2"
)

const (
	maxNameLen    = 80
	maxContactLen = 160
	maxMessageLen = 2000
)

func SendEmail(name, contact, message string) error {
	name, contact, message, err := normalizeContactPayload(name, contact, message)
	if err != nil {
		return err
	}

	apiKey := os.Getenv("RESEND_API_KEY")
	if apiKey == "" {
		log.Println("[💥 ERROR] RESEND_API_KEY environment variable is empty!")
		return fmt.Errorf("internal server configuration error")
	}

	// O pacote oficial exporta como resend.NewClient automaticamente
	client := resend.NewClient(apiKey)

	htmlContent := fmt.Sprintf(`
        <h3>🔋 Novo Lead do Portfólio SSH/TCP</h3>
        <p><strong>Nome do Operador:</strong> %s</p>
        <p><strong>Contato/Up-link:</strong> %s</p>
        <p><strong>Mensagem Disparada:</strong></p>
        <blockquote style="background: #f4f4f4; padding: 10px; border-left: 3px solid #00ffff;">
            %s
        </blockquote>
    `, html.EscapeString(name), html.EscapeString(contact), html.EscapeString(message))

	params := &resend.SendEmailRequest{
		From:    "Portfolio SSH <onboarding@resend.dev>",
		To:      []string{"daniel.cmdline@humanoid.net"},
		Subject: fmt.Sprintf("[COM-LINK] %s enviou uma mensagem!", name),
		Html:    htmlContent,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	sent, err := client.Emails.SendWithContext(ctx, params)
	if err != nil {
		log.Printf("[💥 ERROR] Failed to send email via Resend: %v\n", err)
		return fmt.Errorf("failed to dispatch email: %v", err)
	}

	log.Printf("[✓ SUCCESS] Email sent from: %s (%s) with ID %s\n", name, contact, sent.Id)
	return nil
}

func normalizeContactPayload(name, contact, message string) (string, string, string, error) {
	name = strings.TrimSpace(name)
	contact = strings.TrimSpace(contact)
	message = strings.TrimSpace(message)

	if name == "" || contact == "" || message == "" {
		return "", "", "", fmt.Errorf("missing mandatory fields")
	}
	if len(name) > maxNameLen {
		return "", "", "", fmt.Errorf("name is too long")
	}
	if len(contact) > maxContactLen {
		return "", "", "", fmt.Errorf("contact is too long")
	}
	if len(message) > maxMessageLen {
		return "", "", "", fmt.Errorf("message is too long")
	}

	return name, contact, message, nil
}
