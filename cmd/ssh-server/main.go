package main

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/ssh"
	"portfolio-ssh/input"
	"portfolio-ssh/ui"
	"portfolio-ssh/utils"
)

type Profile struct {
	Name           string
	Age            int
	Role           string
	Education      string
	Stack          []string //nolint:revive
	Certifications []string
	Projects       []ui.UIProject // FIX: Usando o tipo correto do pacote ui
}

const contactBodyLimit = 8 << 10

func getMyProfile() Profile {
	p1 := ui.UIProject{ // FIX: Instanciando como ui.UIProject
		Name:        "Seguramos",
		Description: "Plataforma full-stack de corretagem de seguros digital corporativa.",
		TechStack:   []string{"React", "Typescript", "Node.js", "PostgreSQL"},
	}

	p2 := ui.UIProject{ // FIX: Instanciando como ui.UIProject
		Name:        "Go TUI Portfolio",
		Description: "Servidor SSH concorrente multiplataforma assíncrono e criptografado escrito do zero.",
		TechStack:   []string{"Go", "SSH Protocol", "RFC 4251", "Cryptography", "Linux Kernel"},
	}

	return Profile{
		Name:      "Daniel Caesar Mantilha",
		Age:       34,
		Role:      "Systems & Full-Stack Software Engineer // Network Engineer",
		Education: "Sistemas de Informação (Foco em Engenharia de Software)",
		Stack: []string{
			"Linux/GNU", "Go (Golang)", "C Lang", "Node.js", "TypeScript",
			"PostgreSQL", "HTTP/Websockets", "REST APIs", "React/Next.js",
		},
		Certifications: []string{
			"CCNA (Cisco Certified Network Associate) - ID: Enterprise & Security Core",
			"CAE (Certificate in Advanced English) - University of Cambridge",
		},
		Projects: []ui.UIProject{p1, p2},
	}
}

// handleSSHChannel orquestra a TUI cyberpunk dentro da sessão criptografada do SSH
func handleSSHChannel(ch ssh.Channel, requests <-chan *ssh.Request) {
	defer ch.Close()

	// Trata requisições globais do canal SSH (como manter a conexão viva ou redimensionar janela)
	go func() {
		for req := range requests {
			if req.Type == "shell" || req.Type == "pty-req" {
				req.Reply(true, nil)
			} else {
				req.Reply(false, nil)
			}
		}
	}()

	// Roda a introdução cinematográfica da Matrix assim que o canal SSH abre
	ui.DrawMatrixIntro(ch)
	time.Sleep(1500 * time.Millisecond)

	profile := getMyProfile()
	menuItems := []string{
		"👤 [1] DECRYPT PROFILE (SOBRE MIM)",
		"🐙 [2] SYNC GITHUB (REPOSITORIOS PUBLICOS)",
		"💾 [3] ACCESS DATABASE (PROJETOS EM PRODUÇÃO)",
		"📨 [4] ESTABLISH CONTACT (ENVIA UM E-MAIL VIA ENDPOINT HTTP)",
		"❌ [5] TERMINATE SESSION (SAIR)",
	}
	cursor := 0
	reader := bufio.NewReader(ch)

	for {
		ch.Write([]byte("\033[2J\033[H"))
		ui.DrawCyberBanner(ch)

		ch.Write([]byte("\033[1;32m┌── SELECT DESTINATION PROTOCOL ──────────────────────────────────────────────────┐\r\n"))
		for idx, item := range menuItems {
			if idx == cursor {
				ch.Write([]byte(fmt.Sprintf("  \033[1;30;106m ➔ %-76s \033[0m\r\n", item)))
			} else {
				ch.Write([]byte(fmt.Sprintf("     \033[0;32m%-76s\033[0m \r\n", item)))
			}
		}
		ch.Write([]byte("\033[1;32m└─────────────────────────────────────────────────────────────────────────────────┘\r\n\033[0m"))
		ui.DrawStatusBar(ch, menuItems[cursor][4:])

		b, err := reader.ReadByte()
		if err != nil {
			break
		}

		// Tratamento estável de setas direcionais no ambiente de canais SSH
		if b == 27 {
			b2, _ := reader.ReadByte()
			b3, _ := reader.ReadByte()
			if b2 == 91 {
				if b3 == 65 { // CIMA
					if cursor > 0 {
						cursor--
					} else {
						cursor = len(menuItems) - 1
					}
				}
				if b3 == 66 { // BAIXO
					if cursor < len(menuItems)-1 {
						cursor++
					} else {
						cursor = 0
					}
				}
			}
			continue
		}

		if b == 10 || b == 13 {
			ch.Write([]byte("\033[2J\033[H"))

			switch cursor {
			case 0: // Sobre Mim
				ui.DrawAboutMe(ch, profile.Name, profile.Age, profile.Role, profile.Education, profile.Stack, profile.Certifications)

			case 1: // GitHub Node (Live data)
				ui.DrawCyberBanner(ch)
				ui.FetchAndDrawGitHub(ch, "daniel-cmdline")

			case 2: // Meus Projetos
				ui.DrawProjects(ch, profile.Projects)

			case 3: // Comms Window (Integração Resend)
				ui.DrawCommsEnvelopeTop(ch)
				name := input.ReadLine(ch, reader, "         \033[1;33m│\033[0m  \033[1;32mFROM (NOME):\033[0m ", 0)
				ch.Write([]byte("         \033[1;33m│\033[0m                                                             \033[1;33m│\033[0m\r\n"))

				contact := input.ReadLine(ch, reader, "         \033[1;33m│\033[0m  \033[1;32mUP-LINK (EMAIL/LINK):\033[0m ", 0)
				ch.Write([]byte("         \033[1;33m│\033[0m                                                             \033[1;33m│\033[0m\r\n"))
				ch.Write([]byte("         \033[1;33m├─────────────────────────────────────────────────────────────┤\033[0m\r\n"))
				ch.Write([]byte("         \033[1;33m│\033[0m  \033[1;36mPAYLOAD (MENSAGEM):\033[0m                                         \033[1;33m│\033[0m\r\n"))

				msg := input.ReadLine(ch, reader, "         \033[1;33m│\033[0m  ➔  ", 50)
				ui.DrawCommsEnvelopeBottom(ch)

				ch.Write([]byte("\r\n\r\n  \033[1;31m[!] SEALING ENVELOPE... \033[0m\r\n"))
				time.Sleep(400 * time.Millisecond)
				ch.Write([]byte("  \033[1;36m[*] DISPATCHING PACKETS VIA SECURE HTTP POST... \033[0m"))

				// FIX: Aponta para a nova função pública dentro de utils
				err := utils.SendEmail(name, contact, msg)
				if err != nil {
					ch.Write([]byte(fmt.Sprintf("\r\n\r\n  \033[1;31m[!] TRANSMISSION FAILURE: %v\033[0m\r\n", err.Error())))
				} else {
					ch.Write([]byte("\r\n\r\n  \033[1;32m[+] SUCCESS: MAIL SENT DIRECTLY TO DANIEL'S CORE HANDSET.\033[0m\r\n"))
				}
				ch.Write([]byte("\033[1;30m  ─────────────────────────────────────────────────────────────────────────────────\033[0m\r\n"))

			case 4: // Sair
				ch.Write([]byte("\033[1;31m\r\n[!] TERMINATING CRYPTO SESSION... BYE.\r\n\033[0m"))
				return
			}

			ch.Write([]byte("\r\n\033[5;1;91m➔ PRESS ANY KEY TO RETURN TO CORE OS...\033[0m\r\n"))
			reader.ReadByte()
		}
	}
}

// ContactPayload define a estrutura do JSON para o endpoint de contato.
type ContactPayload struct {
	Name    string `json:"name"`
	Contact string `json:"contact"`
	Message string `json:"message"`
}

// handleContact é o handler HTTP para o formulário de contato.
func handleContact(w http.ResponseWriter, r *http.Request) {
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

func loadHostKey() (ssh.Signer, error) {
	hostKeyPath := os.Getenv("SSH_HOST_KEY_PATH")
	if hostKeyPath == "" {
		hostKeyPath = "id_rsa"
	}

	privateBytes, err := os.ReadFile(hostKeyPath)
	if err == nil {
		return ssh.ParsePrivateKey(privateBytes)
	}
	if !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to read host key %q: %w", hostKeyPath, err)
	}

	log.Printf("Host key %q not found; generating ephemeral SSH host key for this process.", hostKeyPath)
	return generateEphemeralHostKey()
}

func generateEphemeralHostKey() (ssh.Signer, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("failed to generate ephemeral host key: %w", err)
	}

	pemBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	return ssh.ParsePrivateKey(pemBytes)
}

func main() {
	//Load env
	err := godotenv.Load()
	// 1. Configura os parâmetros globais do protocolo SSH
	config := &ssh.ServerConfig{
		PasswordCallback: func(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
			return nil, nil
		},
		NoClientAuth: true,
	}

	// 2. Carrega a Host Key SSH de arquivo ou gera uma chave efêmera para containers.
	private, err := loadHostKey()
	if err != nil {
		log.Fatalf("Falha crítica ao preparar a Host Key privada: %v", err)
	}
	config.AddHostKey(private)

	// 3. Lê dinamicamente a porta definida no ambiente ou assume a 2222 por padrão
	port := os.Getenv("PORT")
	if port == "" {
		port = "2222"
	}

	listenAddr := fmt.Sprintf("0.0.0.0:%s", port)
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalf("Erro ao abrir socket: %v", err)
	}
	defer listener.Close()
	log.Printf("Iniciando Servidor Criptografado SSH em %s...\n", listenAddr)

	// 4. Inicia o servidor HTTP para a API de contato em uma goroutine
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/api/contact", handleContact)
		server := &http.Server{
			Addr:              ":8080",
			Handler:           mux,
			ReadHeaderTimeout: 5 * time.Second,
		}

		log.Println("Iniciando API HTTP interna na porta :8080...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("API HTTP interna falhou: %v", err)
		}
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		// 5. Goroutine concorrente para isolar o handshake SSH de cada terminal conectado
		go func(c net.Conn) {
			sshConn, chans, reqs, err := ssh.NewServerConn(c, config)
			if err != nil {
				c.Close()
				return
			}
			go ssh.DiscardRequests(reqs)

			for newChannel := range chans {
				if newChannel.ChannelType() != "session" {
					newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
					continue
				}

				ch, requests, err := newChannel.Accept()
				if err != nil {
					continue
				}

				go handleSSHChannel(ch, requests)
			}
			sshConn.Close()
		}(conn)
	}
}
