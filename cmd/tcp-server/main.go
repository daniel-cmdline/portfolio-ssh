package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"

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
		Description: "Plataforma full-stack de corretagem de seguros digital.",
		TechStack:   []string{"React", "Typescript", "Node.js", "PostgreSQL"},
	}

	p2 := ui.UIProject{
		Name:        "Go TUI Portfolio",
		Description: "Servidor TCP concorrente multiplataforma escrito do zero.",
		TechStack:   []string{"Go", "TCP Sockets", "Linux Kernel(FD)"},
	}

	return Profile{
		Name:      "Daniel Caesar Mantilha",
		Age:       38,
		Role:      "Full-Stack Software Engineer & Network Engineer",
		Education: "Sistemas de Informação (Foco em Engenharia de Software)",
		Stack:     []string{"Node.js", "Express", "TypeScript", "React", "Next.js", "Python", "Go", "C", "PostgresSQL", "DrizzleORM"},
		Certifications: []string{
			"CCNA (Cisco Certified Network Associate)",
			"CAE (Certificate in Advanced English) - University of Cambridge",
		},
		Projects: []ui.UIProject{p1, p2},
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// 🔥 AQUI SIM! Roda a intro da Matrix na conexão TCP do usuário!
	ui.DrawMatrixIntro(conn)

	// Dá 1.5 segundos para o cara processar o "BOOM" antes do loop limpar a tela
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
	reader := bufio.NewReader(conn)

	for {
		conn.Write([]byte("\033[2J\033[H"))
		ui.DrawCyberBanner(conn)

		conn.Write([]byte("\033[1;32m┌── SELECT DESTINATION PROTOCOL ──────────────────────────────────────────────────┐\r\n"))
		for idx, item := range menuItems {
			if idx == cursor {
				conn.Write([]byte(fmt.Sprintf("│  \033[1;30;106m ➔ %-76s \033[0m│\r\n", item)))
			} else {
				conn.Write([]byte(fmt.Sprintf("│     \033[0;32m%-76s\033[0m │\r\n", item)))
			}
		}
		conn.Write([]byte("\033[1;32m└─────────────────────────────────────────────────────────────────────────────────┘\r\n\033[0m"))
		ui.DrawStatusBar(conn, menuItems[cursor][4:])

		b, err := reader.ReadByte()
		if err != nil {
			break
		}

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
			conn.Write([]byte("\033[2J\033[H")) 

			switch cursor {
			case 0: // Sobre Mim
				ui.DrawAboutMe(conn, profile.Name, profile.Age, profile.Role, profile.Education, profile.Stack, profile.Certifications)
			
			case 1: // GitHub Node
				ui.DrawCyberBanner(conn)
				ui.FetchAndDrawGitHub(conn, "daniel-cmdline")
			
			case 2: // Meus Projetos
				ui.DrawProjects(conn, profile.Projects)
			
			case 3: // Comms Window (Orquestração Segura)
				ui.DrawCommsEnvelopeTop(conn)

				name := input.ReadLine(conn, reader, "         \033[1;33m│\033[0m  \033[1;32mFROM (NOME):\033[0m ")
				conn.Write([]byte("         \033[1;33m│\033[0m                                                             \033[1;33m│\033[0m\r\n"))

				contact := input.ReadLine(conn, reader, "         \033[1;33m│\033[0m  \033[1;32mUP-LINK (EMAIL/LINK):\033[0m ")
				conn.Write([]byte("         \033[1;33m│\033[0m                                                             \033[1;33m│\033[0m\r\n"))
				
				conn.Write([]byte("         \033[1;33m├─────────────────────────────────────────────────────────────┤\033[0m\r\n"))
				conn.Write([]byte("         \033[1;33m│\033[0m  \033[1;36mPAYLOAD (MENSAGEM):\033[0m                                         \033[1;33m│\033[0m\r\n"))

				msg := input.ReadLine(conn, reader, "         \033[1;33m│\033[0m  ➔  ")

				ui.DrawCommsEnvelopeBottom(conn)

				conn.Write([]byte("\r\n\r\n  \033[1;31m[!] SEALING ENVELOPE... \033[0m\r\n"))
				time.Sleep(400 * time.Millisecond)
				conn.Write([]byte("  \033[1;36m[*] DISPATCHING PACKETS VIA SECURE HTTP POST... \033[0m"))

				err := api.SendContact(name, contact, msg)
				if err != nil {
					conn.Write([]byte(fmt.Sprintf("\r\n\r\n  \033[1;31m[!] TRANSMISSION FAILURE: %v\033[0m\r\n", err)))
				} else {
					conn.Write([]byte("\r\n\r\n  \033[1;32m[+] SUCCESS: MAIL SENT DIRECTLY TO DANIEL'S CORE HANDSET.\033[0m\r\n"))
				}
				conn.Write([]byte("\033[1;30m  ─────────────────────────────────────────────────────────────────────────────────\033[0m\r\n"))
			
			case 4: // Sair
				conn.Write([]byte("\033[1;31m\r\n[!] DISCONNECTING FROM HOST... BYE.\r\n\033[0m"))
				return
			}

			conn.Write([]byte("\r\n\033[5;1;91m➔ PRESS ANY KEY TO RETURN TO CORE OS...\033[0m\r\n"))
			reader.ReadByte() 
		}
	}
}

func main() {
	fmt.Println("Iniciando o servidor TCP na porta 2222...")
	listener, err := net.Listen("tcp", ":2222")
	if err != nil {
		fmt.Println("Erro ao abrir o socket:", err)
		os.Exit(1)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleConnection(conn)
	}
}
