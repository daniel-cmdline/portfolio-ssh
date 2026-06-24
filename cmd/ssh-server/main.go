package main

import (
	//std libs
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	//imported pckgs
	"github.com/joho/godotenv"
	"golang.org/x/crypto/ssh"
	"portfolio-ssh/functions"
)

func main() {
	// Load env (Popula o processo com os dados do .env antes de tudo)
	_ = godotenv.Load()

	// 1. Configura os parâmetros globais do protocolo SSH
	config := &ssh.ServerConfig{
		PasswordCallback: func(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
			return nil, nil
		},
		NoClientAuth: true,
	}

	// 2. Carrega a Host Key SSH isolada no pacote functions (disco ou efêmera)
	private, err := functions.LoadHostKey()
	if err != nil {
		log.Fatalf("Falha crítica ao preparar a Host Key privada: %v", err)
	}
	config.AddHostKey(private)

	// 3. Lê dinamicamente a porta definida no ambiente ou assume a 2222 por padrão
	port := os.Getenv("PORT")
	if port == "" {
		port = "2222"
	}

	// Abre o socket em todas as interfaces
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
		mux.HandleFunc("/api/contact", functions.HandleContact)
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

				// Roda a TUI cyberpunk isolada do pacote functions
				go functions.HandleSSHChannel(ch, requests)
			}
			sshConn.Close()
		}(conn)
	}
}