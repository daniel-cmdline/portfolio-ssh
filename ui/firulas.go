package ui

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"
)

// GitHubUser mapeia exatamente o que a API do GitHub devolve e que nos interessa
type GitHubUser struct {
	Login       string `json:"login"`
	PublicRepos int    `json:"public_repos"`
	Followers   int    `json:"followers"`
	Bio         string `json:"bio"`
}

type GitHubRepo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Language    string `json:"language"`
	Stargazers  int    `json:"stargazers_count"`
}

// FetchAndDrawGitHub puxa usuário + repositórios recentes e desenha com ASCII Art do Octocat
func FetchAndDrawGitHub(conn net.Conn, username string) {
	conn.Write([]byte("\033[1;33m[!] CONNECTING TO CORE GITHUB API NODE...\033[0m\r\n"))

	client := &http.Client{Timeout: 6 * time.Second}

	// 1. Pega dados do Usuário
	respUser, err := client.Get(fmt.Sprintf("https://api.github.com/users/%s", username))
	if err != nil || respUser.StatusCode != 200 {
		conn.Write([]byte("\033[1;31m[💥] ERROR: UNABLE TO FETCH USER NODE\033[0m\r\n"))
		return
	}
	defer respUser.Body.Close() // <-- LINHA CORRIGIDA AQUI!
	var user GitHubUser
	json.NewDecoder(respUser.Body).Decode(&user)

	// 2. Pega os 3 repositórios mais recentes
	respRepos, err := client.Get(fmt.Sprintf("https://api.github.com/users/%s/repos?sort=updated&per_page=3", username))
	if err != nil || respRepos.StatusCode != 200 {
		conn.Write([]byte("\033[1;31m[💥] ERROR: UNABLE TO FETCH REPOSITORIES NODE\033[0m\r\n"))
		return
	}
	defer respRepos.Body.Close()
	var repos []GitHubRepo
	json.NewDecoder(respRepos.Body).Decode(&repos)

	conn.Write([]byte("\033[1;32m[✓] SECTOR SYNCHRONIZED! RENDERING DATAFEED...\033[0m\r\n\r\n"))
	time.Sleep(200 * time.Millisecond)

	// ASCII Art do Octocat + Bloco de Informações do Usuário
	conn.Write([]byte("\033[1;36m┌── [ GITHUB REMOTE TELEMETRY ] ──────────────────────────────────────────────────┐\033[0m\r\n"))
	
	// Printando dados básicos junto com a silhueta do Octocat
	conn.Write([]byte("  \033[1;35m  _(\\ _/)_ \033[0m    \033[1;32mOPERATIVE:\033[0m  " + user.Login + "\r\n"))
	conn.Write([]byte("  \033[1;35m ((  \"  )) \033[0m    \033[1;32mNET REPOS:\033[0m  " + fmt.Sprintf("%d Public Units", user.PublicRepos) + "\r\n"))
	conn.Write([]byte("  \033[1;35m  /\\-V-/\\  \033[0m    \033[1;32mFOLLOWERS:\033[0m  " + fmt.Sprintf("%d Active Nodes", user.Followers) + "\r\n"))
	conn.Write([]byte("  \033[1;35m (___|___) \033[0m    \033[1;32mBIOGRAPHY:\033[0m  " + user.Bio + "\r\n"))
	
	conn.Write([]byte("\033[1;36m├─────────────────────────────────────────────────────────────────────────────────┤\033[0m\r\n"))
	conn.Write([]byte("  \033[1;33m📡 LIVE DEPLOYMENTS (RECENTLY UPDATED ON GITHUB):\033[0m\r\n"))

	// Loop para varrer e printar os repositórios reais da nuvem!
	for _, repo := range repos {
		desc := repo.Description
		if desc == "" {
			desc = "No description provided."
		}
		
		textoRepo := fmt.Sprintf("\r\n    \033[1;36m📁 %s\033[0m\r\n    ├─ \033[90mSYS_DESC:\033[0m %s\r\n    └─ \033[90mCORE_LNG:\033[0m \033[1;34m%s\033[0m  │  \033[1;33m⭐ %d\033[0m\r\n", 
			repo.Name, desc, repo.Language, repo.Stargazers)
		
		Typewriter(conn, textoRepo, 10*time.Millisecond)
	}

	conn.Write([]byte("\033[1;36m└─────────────────────────────────────────────────────────────────────────────────┘\033[0m\r\n"))
}

// Typewriter digita os caracteres de forma gradual na rede
func Typewriter(conn net.Conn, text string, delay time.Duration) {
	for _, char := range text {
		conn.Write([]byte(string(char)))
		time.Sleep(delay)
	}
}

// DrawCyberBanner cospe um logo gigante em ASCII Art com o nome de guerra DEFINITIVO
func DrawCyberBanner(conn net.Conn) {
	// Cores Neon de alta intensidade
	green := "\033[1;92m"
	cyan  := "\033[1;96m"
	gray  := "\033[90m"
	reset := "\033[0m"

	// ASCII ART: DANIEL CMD LINE (Letras totalmente isoladas e limpas)
	conn.Write([]byte(green + "  ____   _   _   _ ___ _____ _       ____ __  __ ____   _     ___ _   _ _____\r\n" + reset))
	conn.Write([]byte(green + " |  _ \\ /_\\ | \\ | |_ _| ____| |     / ___|  \\/  |  _ \\  | |   |_ _| \\ | | ____|\r\n" + reset))
	conn.Write([]byte(cyan  + " | | | / _ \\|  \\| | | ||  _| | |    | |   | |\\/| | | | | | |    | ||  \\| |  _|  \r\n" + reset))
	conn.Write([]byte(cyan  + " | |_| / ___ \\ |\\  | | || |___| |___ | |___| |  | | |_| | | |___ | || |\\  | |___ \r\n" + reset))
	conn.Write([]byte(cyan  + " |____/_/   \\_\\_| \\_|___|_____|_____| \\____|_|  |_|____/  |_____|___|_| \\_|_____|\r\n" + reset))
	conn.Write([]byte(gray  + "  ─── [ HOST OVERRIDE: DANIEL_CMD_LINE ] ─── v2.0-RAW ─────────────────────────────\r\n\r\n" + reset))
}

// DrawStatusBar cria uma barra sólida estilo Vim/Tmux no rodapé da tela
func DrawStatusBar(conn net.Conn, currentOption string) {
	conn.Write([]byte("\r\n"))
	conn.Write([]byte("\033[1;30;106m 🖥️  PORT: 2222 \033[1;37;45m ⚡ ENGINE: GO \033[1;37;40m ➔ ACTIVE: " + currentOption + " \033[0m\r\n"))
}