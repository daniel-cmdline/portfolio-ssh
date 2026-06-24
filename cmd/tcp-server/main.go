package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
	"portfolio-ssh/api"   
	"portfolio-ssh/input" 
	"portfolio-ssh/ui"
)

type Profile struct {
	Name           string
	Age            int
	Role           string
	Education      string
	Stack          []string
	Certifications []string
	Projects       []ui.UIProject 
}

func getMyProfile() Profile {
	p1 := ui.UIProject{
		Name:        "Seguramos",
		Description: "Plataforma full-stack de corretagem de seguros digital corporativa.",
		TechStack:   []string{"React", "Typescript", "Node.js", "PostgreSQL"},
	}

	p2 := ui.UIProject{
		Name:        "Go TUI Portfolio",
		Description: "Servidor SSH concorrente multiplataforma assíncrono e criptografado escrito do zero.",
		TechStack:   []string{"Go", "SSH Protocol", "RFC 4251", "Cryptography", "Linux Kernel"},
	}

	return Profile{
		Name:      "Daniel Caesar Mantilha",
		Age:       38,
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
		"🐙 [2] SYNC GITHUB NODE (LIVE DATA)",
		"💾 [3] ACCESS DATABASE (PROJETOS)",
		"📡 [4] ESTABLISH COMS (ENTRAR EM CONTATO)",
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
				ch.Write([]byte(fmt.Sprintf("│  \033[1;30;106m ➔ %-76s \033[0m│\r\n", item)))
			} else {
				ch.Write([]byte(fmt.Sprintf("│     \033[0;32m%-76s\033[0m │\r\n", item)))
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
					if cursor > 0 { cursor-- } else { cursor = len(menuItems) - 1 }
				}
				if b3 == 66 { // BAIXO
					if cursor < len(menuItems)-1 { cursor++ } else { cursor = 0 }
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
				name := input.ReadLine(ch, reader, "         \033[1;33m│\033[0m  \033[1;32mFROM (NOME):\033[0m ")
				ch.Write([]byte("         \033[1;33m│\033[0m                                                             \033[1;33m│\033[0m\r\n"))
				contact := input.ReadLine(ch, reader, "         \033[1;33m│\033[0m  \033[1;32mUP-LINK (EMAIL/LINK):\033[0m ")
				ch.Write([]byte("         \033[1;33m│\033[0m                                                             \033[1;33m│\033[0m\r\n"))
				ch.Write([]byte("         \033[1;33m├─────────────────────────────────────────────────────────────┤\033[0m\r\n"))
				ch.Write([]byte("         \033[1;33m│\033[0m  \033[1;36mPAYLOAD (MENSAGEM):\033[0m                                         \033[1;33m│\033[0m\r\n"))
				msg := input.ReadLine(ch, reader, "         \033[1;33m│\033[0m  ➔  ")
				ui.DrawCommsEnvelopeBottom(ch)

				ch.Write([]byte("\r\n\r\n  \033[1;31m[!] SEALING ENVELOPE... \033[0m\r\n"))
				time.Sleep(400 * time.Millisecond)
				ch.Write([]byte("  \033[1;36m[*] DISPATCHING PACKETS VIA SECURE HTTP POST... \033[0m"))

				err := api.SendContact(name, contact, msg)
				if err != nil {
					ch.Write([]byte(fmt.Sprintf("\r\n\r\n  \033[1;31m[!] TRANSMISSION FAILURE: %v\033[0m\r\n", err)))
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

func main() {
	// 1. Configura os parâmetros globais do protocolo SSH
	config := &ssh.ServerConfig{
		// Permite conexões anônimas sem checagem de senhas (portfólio público)
		PasswordCallback: func(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
			return nil, nil
		},
		NoClientAuth: true,
	}

	// 2. Carrega a Host Key RSA gerada pelo ssh-keygen
	privateBytes, err := os.ReadFile("id_rsa")
	if err != nil {
		log.Fatalf("Falha crítica: arquivo de chave 'id_rsa' ausente na raiz! Execute: ssh-keygen -t rsa -b 2048 -f id_rsa -N \"\" -> Erro: %v", err)
	}

	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		log.Fatalf("Falha crítica ao decodificar a Host Key privada: %v", err)
	}
	config.AddHostKey(private)

	// 3. Lê dinamicamente a porta definida no ambiente ou assume a 2222 por padrão
	port := os.Getenv("PORT")
	if port == "" {
		port = "2222"
	}

	fmt.Printf("Iniciando Servidor Criptografado SSH na porta %s...\n", port)
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Erro crítico ao abrir socket TCP na porta %s: %v", port, err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		// 4. Goroutine concorrente para isolar o handshake SSH de cada terminal conectado
		go func(c net.Conn) {
			// Executa a negociação criptográfica (Handshake) em cima do cano TCP aceito
			sshConn, chans, reqs, err := ssh.NewServerConn(c, config)
			if err != nil {
				c.Close()
				return
			}
			go ssh.DiscardRequests(reqs)

			// Escuta as requisições de sub-canais gerados pelo cliente SSH
			for newChannel := range chans {
				// Rejeita qualquer canal que não seja do tipo de terminal clássico ("session")
				if newChannel.ChannelType() != "session" {
					newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
					continue
				}

				// Aceita o canal de sessão e recupera a ponte estável de I/O
				ch, requests, err := newChannel.Accept()
				if err != nil {
					continue
				}

				// Passa o canal criptografado para carregar a TUI de forma isolada
				go handleSSHChannel(ch, requests)
			}
			sshConn.Close()
		}(conn)
	}
}
