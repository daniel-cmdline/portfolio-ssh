package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"portfolio-ssh/ui" // Pacote modularizado ativo
)

type Project struct {
	Name        string
	Description string
	TechStack   []string
}

type Profile struct {
	Name           string
	Age            int
	Role           string
	Education      string
	Stack          []string
	Certifications []string
	Projects       []Project
}

func getMyProfile() Profile {
	p1 := Project{
		Name:        "Seguramos",
		Description: "Plataforma full-stack de corretagem de seguros digital.",
		TechStack:   []string{"React", "Typescript", "Node.js", "PostgreSQL"},
	}

	p2 := Project{
		Name:        "Go TUI Portfolio",
		Description: "Servidor TCP concorrente multiplataforma escrito do zero.",
		TechStack:   []string{"Go", "TCP Sockets", "Linux Kernel(FD)"},
	}

	myProfile := Profile{
		Name:      "Daniel Caesar Mantilha",
		Age:       34,
		Role:      "Full-Stack Software Engineer & Network Engineer",
		Education: "Sistemas de Informação - FMU (3º Ano)",
		Stack:     []string{"Node.js", "Express", "TypeScript", "React", "Python", "Go", "PostgresSQL", "DrizzleORM"},
		Certifications: []string{
			"CCNA (Cisco Certified Network Associate)",
			"CAE (Certificate in Advanced English) - University of Cambridge",
		},
		Projects: []Project{p1, p2},
	}

	return myProfile
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	profile := getMyProfile()
	menuItems := []string{
		"[1] DECRYPT PROFILE (SOBRE MIM)",
		"[2] SYNC GITHUB NODE (LIVE DATA)",
		"[3] ACCESS DATABASE (PROJETOS)",
		"[4] TERMINATE SESSION (SAIR)",
	}
	cursor := 0

	//Otimizador de leitura de rede que impede seu servidor de engasgar fazendo requisições em excesso pro Kernel
	reader := bufio.NewReader(conn)

	for {
		// Limpa a tela e joga o cursor pro topo
		conn.Write([]byte("\033[2J\033[H"))

		// FIRULA 1: ASCII Art Gigante do Pacote UI
		ui.DrawCyberBanner(conn)

		// Moldura de instrução estendida (Alinhada com o tamanho da ASCII Art nova)
		conn.Write([]byte("\033[1;32m┌── SELECT DESTINATION PROTOCOL ──────────────────────────────────────────────────┐\r\n"))

		// Desenha os itens do menu
		for idx, item := range menuItems {
			if idx == cursor {
				// Cursor selecionado brilhando em Ciano Neon com largura correta
				conn.Write([]byte(fmt.Sprintf("│  \033[1;30;106m ➔ %-76s \033[0m│\r\n", item)))
			} else {
				// Opções normais em tom esverdeado hacker apagado
				conn.Write([]byte(fmt.Sprintf("│     \033[0;32m%-76s\033[0m │\r\n", item)))
			}
		}
		conn.Write([]byte("\033[1;32m└─────────────────────────────────────────────────────────────────────────────────┘\r\n\033[0m"))

		// FIRULA 2: Barra de Status Estilo Tmux/LazyVim no rodapé
		ui.DrawStatusBar(conn, menuItems[cursor][4:])


		//A mecânica de cima e baixo e rederização interceptando da sequencia de bytes ASCII digitados no terminal do usuario 
		b, err := reader.ReadByte()
		if err != nil {
			break
		}

		// Captura de setas
		if b == 27 {
			b2, _ := reader.ReadByte()
			b3, _ := reader.ReadByte()

			if b2 == 91 {
    			if b3 == 65 { // Seta para CIMA
        			if cursor > 0 { 
            		cursor-- 
        			} else { 
           			 cursor = len(menuItems) - 1 // Se passou do topo, pula pro final!
        			} 
    			} 

    			if b3 == 66 { // Seta para BAIXO
        			if cursor < len(menuItems)-1 { 
            		cursor++ 
        			} else { 
            		cursor = 0 // Se passou do final, reseta pro topo!
        			} 
    			} 
			}
			continue
		}

		// Abre o menu se a sequencia for 10 ou 13
		if b == 10 || b == 13 {
			conn.Write([]byte("\033[2J\033[H")) // Limpa a tela para a sub-view
			ui.DrawCyberBanner(conn)

			if cursor == 0 { // Sobre Mim (Ultra Modificado!)
				conn.Write([]byte("\033[1;36m┌── [ DECRYPTING CORE DATA... ] ──────────────────────────────────────────────────┐\033[0m\r\n"))
				
				// Montando blocos de texto estilizados com Typewriter
				infoBasica := fmt.Sprintf(" \033[1;32m> IDENTITY:\033[0m    %s\r\n \033[1;32m> LIFESPAN:\033[0m    %d YEARS IN EARTH NODE\r\n \033[1;32m> CORE ROLE:\033[0m   %s\r\n \033[1;32m> ACADEMICS:\033[0m   🎓 %s\r\n", 
					profile.Name, profile.Age, profile.Role, profile.Education)
				ui.Typewriter(conn, infoBasica, 15*time.Millisecond)

				// Injetando as Stacks
				stackStr := strings.Join(profile.Stack, " │ ")
				conn.Write([]byte(fmt.Sprintf(" \033[1;32m> TECH STACK:\033[0m  🛠️  [ %s ]\r\n", stackStr)))

				// Injetando Certificações com loop dedicado
				conn.Write([]byte(" \033[1;32m> CERTIFICATIONS:\033[0m  \r\n"))
				for _, cert := range profile.Certifications {
					icone := "📜"
					if strings.Contains(cert, "CCNA") {
						icone = "🌐"
					}
					textoCert := fmt.Sprintf("    ├─ %s %s\r\n", icone, cert)
					ui.Typewriter(conn, textoCert, 15*time.Millisecond)
				}
				
				conn.Write([]byte("\033[1;36m└─────────────────────────────────────────────────────────────────────────────────┘\033[0m\r\n"))

			} else if cursor == 1 { // GitHub Live Data Node!
				ui.FetchAndDrawGitHub(conn, "daniel-cmdline")

			} else if cursor == 2 { // Meus Projetos
				conn.Write([]byte("\033[1;35m┌── [ FETCHING REPOSITORIES... ] ─────────────────────────────────────────────────┐\033[0m\r\n"))
				for _, proj := range profile.Projects {
					textoProj := fmt.Sprintf("\r\n \033[1;33m⚡ TARGET:\033[0m %s\r\n    ├─ \033[90mDESC:\033[0m %s\r\n", proj.Name, proj.Description)
					ui.Typewriter(conn, textoProj, 15*time.Millisecond)

					stack := strings.Join(proj.TechStack, ", ")
					conn.Write([]byte(fmt.Sprintf("    └─ \033[90mSTACK:\033[0m \033[1;34m[%s]\033[0m\r\n", stack)))
				}
				conn.Write([]byte("\r\n\033[1;35m└─────────────────────────────────────────────────────────────────────────────────┘\033[0m\r\n"))

			} else if cursor == 3 { // Sair
				conn.Write([]byte("\033[1;31m\r\n[!] DISCONNECTING FROM HOST... BYE.\r\n\033[0m"))
				break
			}

			conn.Write([]byte("\r\n\033[5;1;91m➔ PRESS ANY KEY TO RETURN TO CORE OS...\033[0m\r\n"))
			reader.ReadByte() // Trava a tela
		}
	}
}

//Inicia o socket, realiza a conexão e passa ela na função handleConnection como um goroutine
func main() {
	fmt.Println("Iniciando o servidor TCP na porta 2222...")

	bytes := []byte("Sequencia randominca de bytes")

	fmt.Println("Sequência characters transformada em bytes", bytes)

	listener, err := net.Listen("tcp", ":2222")
	if err != nil {
		fmt.Println("Erro ao abrir o socket:", err)
		os.Exit(1)
	}
	defer listener.Close()

	for {
		fmt.Println("Aguardando nova conexão...")
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Erro ao aceitar conexão:", err)
			continue
		}
		go handleConnection(conn)
	}
}