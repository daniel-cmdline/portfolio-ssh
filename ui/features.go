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

// DrawCommsEnvelopeTop renderiza a parte superior e cabeçalho do envelope cyberpunk
func DrawCommsEnvelopeTop(conn net.Conn) {
	conn.Write([]byte("\033[1;35m  📡 COM-LINK DETECTED // SECURE INBOUND CHANNEL // ROUTING: DANIEL_CORE\033[0m\r\n"))
	conn.Write([]byte("\033[1;30m  ─────────────────────────────────────────────────────────────────────────────────\033[0m\r\n\r\n"))

	// O Desenho do Envelope ASCII (O topo da carta)
	conn.Write([]byte("         \033[1;33m┌─────────────────────────────────────────────────────────────┐\033[0m\r\n"))
	conn.Write([]byte("         \033[1;33m│\033[0m  \033[1;34m✉️  DIGITAL ENVELOPE (ESTABLISH COMMUNICATIONS)\033[0m            \033[1;33m│\033[0m\r\n"))
	conn.Write([]byte("         \033[1;33m├─────────────────────────────────────────────────────────────┤\033[0m\r\n"))
	conn.Write([]byte("         \033[1;33m│\033[0m                                                             \033[1;33m│\033[0m\r\n"))
}

// DrawCommsEnvelopeBottom renderiza o fundo do envelope selado
func DrawCommsEnvelopeBottom(conn net.Conn) {
	conn.Write([]byte("         \033[1;33m│\033[0m                                                             \033[1;33m│\033[0m\r\n"))
	conn.Write([]byte("         \033[1;33m└─────────────────────────────────────────────────────────────┘\033[0m\r\n"))
}

// FetchAndDrawGitHub puxa usuário + repositórios recentes e desenha com ASCII Art do Octocat
func FetchAndDrawGitHub(conn net.Conn, username string) {
	conn.Write([]byte("\033[1;33m[!] CONNECTING TO CORE GITHUB API NODE...\033[0m\r\n"))

	client := &http.Client{Timeout: 6 * time.Second}

	respUser, err := client.Get(fmt.Sprintf("https://api.github.com/users/%s", username))
	if err != nil || respUser.StatusCode != 200 {
		conn.Write([]byte("\033[1;31m[💥] ERROR: UNABLE TO FETCH USER NODE\033[0m\r\n"))
		return
	}
	defer respUser.Body.Close()
	var user GitHubUser
	json.NewDecoder(respUser.Body).Decode(&user)

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

	conn.Write([]byte("\033[1;36m┌── [ GITHUB REMOTE TELEMETRY ] ──────────────────────────────────────────────────┐\033[0m\r\n"))
	conn.Write([]byte("  \033[1;35m  _(\\ _/)_ \033[0m    \033[1;32mOPERATIVE:\033[0m  " + user.Login + "\r\n"))
	conn.Write([]byte("  \033[1;35m ((  \"  )) \033[0m    \033[1;32mNET REPOS:\033[0m  " + fmt.Sprintf("%d Public Units", user.PublicRepos) + "\r\n"))
	conn.Write([]byte("  \033[1;35m  /\\-V-/\\  \033[0m    \033[1;32mFOLLOWERS:\033[0m  " + fmt.Sprintf("%d Active Nodes", user.Followers) + "\r\n"))
	conn.Write([]byte("  \033[1;35m (___|___) \033[0m    \033[1;32mBIOGRAPHY:\033[0m  " + user.Bio + "\r\n"))
	
	conn.Write([]byte("\033[1;36m├─────────────────────────────────────────────────────────────────────────────────┤\033[0m\r\n"))
	conn.Write([]byte("  \033[1;33m📡 LIVE DEPLOYMENTS (RECENTLY UPDATED ON GITHUB):\033[0m\r\n"))

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

// DrawCyberBanner cospe um logo gigante em ASCII Art
func DrawCyberBanner(conn net.Conn) {
	green := "\033[1;92m"
	cyan  := "\033[1;96m"
	gray  := "\033[90m"
	reset := "\033[0m"

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