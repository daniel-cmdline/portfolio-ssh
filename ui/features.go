package ui

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

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

type UIProject struct {
	Name        string
	Description string
	TechStack   []string
}

// DrawMatrixIntro aceita qualquer coisa que leia e escreva (io.ReadWriter)
func DrawMatrixIntro(w io.ReadWriter) {
	w.Write([]byte("\033[2J\033[H"))
	time.Sleep(500 * time.Millisecond)

	greenNeon := "\033[1;92m"
	reset := "\033[0m"

	Typewriter(w, greenNeon+" Wake up, Neo...\r\n\r\n"+reset, 80*time.Millisecond)
	time.Sleep(1200 * time.Millisecond)

	Typewriter(w, greenNeon+" The Matrix has you...\r\n\r\n"+reset, 80*time.Millisecond)
	time.Sleep(1200 * time.Millisecond)

	neoASCII := greenNeon + 
		"  ⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿\r\n" +
		"  ⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿\r\n" +
		"  ⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿\r\n" +
		"  ⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠿⠿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿\r\n" +
		"  ⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡿⠟⠉⠉⠉⠀⠀⠉⠙⠿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿\r\n" +
		"  ⣿⣿⣿⣿⣿⣿⣿⣿⣿⠿⠛⠛⠉⠉⠉⠉⠙⠛⠿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠟⠉⠀⠰⣏⣤⣤⣆⣠⣀⣀⣀⣀⡉⠛⣿⣿⣿⣿⣿⣿⣿\r\n" +
		"  ⣿⣿⣿⣿⣿⣿⣿⠋⠀⢲⢢⢤⣀⡀⠀⠀⠀⠀⠀⠀⠈⢻⣿⣿⣿⣿⣿⡟⠀⠀⠀⠀⠘⣿⣾⡝⠻⢷⠿⠟⡿⠋⠀⠘⣿⣿⣿⣿⣿⣿   [ OPERATIVE DETECTED ]\r\n" +
		"  ⣿⣿⣿⣿⠟⠋⠀⠠⣔⣮⣭⠉⠙⠉⢁⡥⠞⠀⠀⠀⠀⠈⣿⣿⣿⣿⣿⡇⠂⠀⠀⠒⢄⠈⠉⠿⣧⡄⠲⠟⣿⣶⣤⡄⠈⠻⣿⣿⣿⣿\r\n" +
		"  ⣿⣿⣋⣥⣤⣲⣼⣽⣿⣿⠗⣦⣄⡴⠋⠔⠁⠀⠀⠀⠀⠀⢸⣿⣿⣿⣿⣿⡤⠀⠀⠀⠀⠉⠈⢀⠹⠿⡟⠒⠛⢿⣿⣾⣤⠄⠙⢿⣿⣿   TRACING LINK...\r\n" +
		"  ⣿⣿⣿⣿⣿⣿⣿⣿⣿⠏⠉⠀⠀⠀⡰⠀⢩⠂⠀⣆⠀⠀⣸⣿⣿⣿⣿⣿⣷⡤⠀⠒⠢⡄⢄⠀⠡⣄⠑⣤⡄⠀⠙⢻⣿⣿⣿⣾⣿⣿\r\n" +
		"  ⣿⣿⣿⣿⣿⣿⣿⣿⣃⡀⠈⣷⡆⢠⣅⣶⠀⢤⣿⠁⠀⣸⣿⣿⣿⣿⣿⣿⣿⣿⣷⡀⠀⠹⣀⠈⠻⣿⣽⡛⢿⣧⠀⠀⠻⣿⣿⣿⣿⣿\r\n" +
		"  ⣿⣿⣿⣿⣿⣿⣿⠇⠀⠀⣽⡟⢠⣿⡿⠀⠀⡼⠤⢄⣰⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣷⡄⠀⠙⣆⠁⠨⠿⣋⠀⠙⢦⠀⠀⢹⣿⣿⣿⣿\r\n" +
		"  ⣿⣿⣿⣿⣿⣿⣿⠠⡀⠀⣿⠃⠾⡿⠃⠀⢠⠃⠀⢨⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣷⣆⢀⣹⡄⠀⠀⠹⡄⠀⠈⢆⠀⢨⣿⣿⣿⣿\r\n" +
		"  ⣿⣿⣿⣿⣿⣿⣿⣧⣄⣤⠃⠀⢸⠁⠀⠀⣾⣀⣴⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣧⣤⣴⣧⣤⣤⣾⣷⣾⣿⣿⣿⣿\r\n" +
		"  ⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣧⣦⣼⣷⣦⣼⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣷⣿⣿⣿⣿⣿\r\n" +
		"  ⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿\r\n\r\n" + reset

	Typewriter(w, neoASCII, 5*time.Millisecond)
	time.Sleep(1000 * time.Millisecond)

	w.Write([]byte("\033[2J\033[H"))

	alienLines := []string{
		"                     ⢀⣤⣶⣶⣶⣶⣦⣤⣀⠀⠀⠀⠀⠀",
		"              ⠀⠀⢀⣤⣶⣶⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣦⡀⠀⠀",
		"              ⠀⢠⣿⣿⣿⣿⣿⠿⣿⣿⣿⣿⠿⠿⢿⣿⣿⣷⡀⠀",
		"              ⠀⢸⣿⡿⠋⠁⠀⠀⠀⠉⠉⠀⠀⠀⠀⠈⢹⣿⡇⠀",
		"              ⠀⢸⣿⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣿⡇⠀",
		"              ⠀⢸⣿⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣿⡇⠀",
		"              ⠀⢸⣿⣠⣤⣶⣶⣶⣦⣀⣀⣴⣶⣶⣶⣤⣄⣿⡇⡀       [ SYSTEM DECRYPT OVERRIDE ]",
		"              ⣿⣿⣿⠻⣿⣿⣿⣿⣿⠟⠻⣿⣿⣿⣿⣿⠟⣿⣿⣿       MEMORY_DUMP: ",
		"              ⣿⣿⣿⠀⠈⠉⠛⠋⠉⠀⠀⠉⠙⠛⠉⠁⠀⣿⣿⣿       SECTOR: ",
		"              ⠙⢿⣿⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣿⡿⠃        STATUS: ",
		"              ⠀⠸⣿⣧⠀⠀⠀⢀⣀⣀⣀⣀⡀⠀⠀⠀⣼⣿⠇⠀       INTRUSION: SUCCESS",
		"              ⠀⠀⠙⢿⣷⣄⠀⠈⠉⠉⠉⠉⠁⠀⣠⣾⡿⠃⠀⠀",
		"              ⠀⠀⠀⢸⣿⣿⣷⣤⣀⣀⣀⣀⣤⣾⣿⣿⡅⠀⠀⠀",
		"              ⠀⠀⢰⣿⣿⣿⣿⣿⣿⡿⠿⢿⣿⣿⣿⣿⣿⡄⠀⠀",
		"              ⠀⠀⠻⠿⠿⠿⠿⠿⠿⠷⠴⠿⠿⠿⠿⠿⠿⠇⠀⠀",
	}

	for i, line := range alienLines {
		var firulaExtra string
		switch i {
		case 7:
			firulaExtra = fmt.Sprintf("0x%X%X%X", time.Now().UnixNano()%999, i*7, time.Now().Second())
		case 8:
			firulaExtra = fmt.Sprintf("0x%X-NET_CORE", 0x3FFF-i*112)
		case 9:
			firulaExtra = "\033[5;1;91mCRITICAL_BYPASS\033[0m\033[1;92m"
		}
		
		w.Write([]byte(greenNeon + line + firulaExtra + "\r\n" + reset))
		time.Sleep(35 * time.Millisecond)
	}
	time.Sleep(1000 * time.Millisecond)

	Typewriter(w, "\r\n"+greenNeon+" Knock, knock.\r\n"+reset, 80*time.Millisecond)
	time.Sleep(1500 * time.Millisecond)

	w.Write([]byte("\033[2J\033[H"))
	w.Write([]byte("\033[1;31m [💥] PROTOCOL OVERRIDDEN // INITIALIZING DANIEL_CMD_LINE OS...\033[0m\r\n"))
	time.Sleep(600 * time.Millisecond)
}

func DrawCommsEnvelopeTop(w io.Writer) {
	w.Write([]byte("\033[1;35m  📡 COM-LINK DETECTED // SECURE INBOUND CHANNEL // ROUTING: DANIEL_CORE\033[0m\r\n"))
	w.Write([]byte("\033[1;30m  ─────────────────────────────────────────────────────────────────────────────────\033[0m\r\n\r\n"))
	w.Write([]byte("         \033[1;33m┌─────────────────────────────────────────────────────────────┐\033[0m\r\n"))
	w.Write([]byte("         \033[1;33m│\033[0m  \033[1;34m✉️  DIGITAL ENVELOPE (ESTABLISH COMMUNICATIONS)\033[0m            \033[1;33m│\033[0m\r\n"))
	w.Write([]byte("         \033[1;33m├─────────────────────────────────────────────────────────────┤\033[0m\r\n"))
	w.Write([]byte("         \033[1;33m│\033[0m                                                             \033[1;33m│\033[0m\r\n"))
}

func DrawCommsEnvelopeBottom(w io.Writer) {
	w.Write([]byte("         \033[1;33m│\033[0m                                                             \033[1;33m│\033[0m\r\n"))
	w.Write([]byte("         \033[1;33m└─────────────────────────────────────────────────────────────┘\033[0m\r\n"))
}

// DrawAboutMe restaurado com efeito Typewriter letra por letra!
func DrawAboutMe(w io.Writer, name string, age int, role string, edu string, stack []string, certs []string) {
	Typewriter(w, "\033[1;36m┌── [ DECRYPTED OPERATIVE DOSSIER ] ──────────────────────────────────────────────┐\033[0m\r\n", 3*time.Millisecond)
	
	Typewriter(w, fmt.Sprintf("  \033[1;32m⚡ OPERATIVE:\033[0m %s\r\n", name), 8*time.Millisecond)
	Typewriter(w, fmt.Sprintf("  \033[1;32m⚡ AGE:\033[0m       %d Years Old\r\n", age), 8*time.Millisecond)
	Typewriter(w, fmt.Sprintf("  \033[1;32m⚡ POSITION:\033[0m  %s\r\n", role), 8*time.Millisecond)
	Typewriter(w, fmt.Sprintf("  \033[1;32m⚡ ACAD:\033[0m      %s\r\n", edu), 8*time.Millisecond)
	
	Typewriter(w, "\033[1;36m├── [ CORE OPERATIONS & PROFILE OVERVIEW ] ───────────────────────────────────────┤\033[0m\r\n", 3*time.Millisecond)
	
	// Logs de sistema escovando as suas reais habilidades
	Typewriter(w, "  \033[1;30m[sys/init]\033[0m \033[1;32mNATIVE LINUX USER\033[0m -> Dominio de ambiente Unix, Shell Scripting e automacao.\r\n", 6*time.Millisecond)
	Typewriter(w, "  \033[1;30m[sys/net] \033[0m \033[1;34mNETWORK INFRASTRUCTURE\033[0m -> Arquitetura de Redes, Protocolos TCP/IP e roteamento.\r\n", 6*time.Millisecond)
	Typewriter(w, "  \033[1;30m[sys/core]\033[0m \033[1;33mBACKEND SPECIALIST\033[0m -> Construcao de APIs HTTP robustas, servicos em Go e Node.\r\n", 6*time.Millisecond)
	Typewriter(w, "  \033[1;30m[sys/sync]\033[0m \033[1;35mREAL-TIME SYSTEMS\033[0m -> Implementacao de comunicacao bidirecional via WebSockets.\r\n", 6*time.Millisecond)
	Typewriter(w, "  \033[1;30m[sys/web] \033[0m \033[1;36mWEB DEVELOPMENT\033[0m -> Engenharia Full-Stack ponta a ponta com TypeScript e React.\r\n", 6*time.Millisecond)
	
	Typewriter(w, "\033[1;36m├── [ HARDWARE & SOFTWARE STACK ] ────────────────────────────────────────────────┤\033[0m\r\n", 3*time.Millisecond)
	Typewriter(w, "  \033[1;33m🛠️  TECH RUNTIMES:\033[0m\r\n   ", 8*time.Millisecond)
	
	for _, tech := range stack {
		Typewriter(w, "\033[1;30;107m " + tech + " \033[0m ", 4*time.Millisecond)
	}
	
	Typewriter(w, "\r\n\r\n  \033[1;35m🎖️  VALIDATED NETWORK & LANG CERTIFICATIONS:\033[0m\r\n", 8*time.Millisecond)
	for _, cert := range certs {
		Typewriter(w, "   ➔ \033[1;94m" + cert + "\033[0m\r\n", 8*time.Millisecond)
	}
	
	Typewriter(w, "\033[1;36m└─────────────────────────────────────────────────────────────────────────────────┘\r\n", 3*time.Millisecond)
}

// DrawProjects restaurado com efeito Typewriter bloco a bloco!
func DrawProjects(w io.Writer, projects []UIProject) {
	Typewriter(w, "\033[1;36m┌── [ LOCAL DATABASE SECTORS - PRODUCTION DEPLOYS ] ─────────────────────────────┐\033[0m\r\n", 4*time.Millisecond)
	
	for _, p := range projects {
		Typewriter(w, "  \033[1;33m📁 PROJECT: " + p.Name + "\033[0m\r\n", 10*time.Millisecond)
		Typewriter(w, "  │  \033[90mDescription:\033[0m " + p.Description + "\r\n", 8*time.Millisecond)
		Typewriter(w, "  │  \033[90mInfrastructure Stack:\033[0m ", 10*time.Millisecond)
		
		for _, t := range p.TechStack {
			Typewriter(w, "\033[1;34m[" + t + "]\033[0m ", 5*time.Millisecond)
		}
		Typewriter(w, "\r\n  │\r\n", 10*time.Millisecond)
	}
	
	Typewriter(w, "\033[1;36m└─────────────────────────────────────────────────────────────────────────────────┘\r\n", 4*time.Millisecond)
}

func FetchAndDrawGitHub(w io.Writer, username string) {
	w.Write([]byte("\033[1;33m[!] CONNECTING TO CORE GITHUB API NODE...\033[0m\r\n"))
	client := &http.Client{Timeout: 6 * time.Second}
	respUser, err := client.Get(fmt.Sprintf("https://api.github.com/users/%s", username))
	if err != nil || respUser.StatusCode != 200 {
		w.Write([]byte("\033[1;31m[💥] ERROR: UNABLE TO FETCH USER NODE\033[0m\r\n"))
		return
	}
	defer respUser.Body.Close()
	var user GitHubUser
	json.NewDecoder(respUser.Body).Decode(&user)

	respRepos, err := client.Get(fmt.Sprintf("https://api.github.com/users/%s/repos?sort=updated&per_page=3", username))
	if err != nil || respRepos.StatusCode != 200 {
		w.Write([]byte("\033[1;31m[💥] ERROR: UNABLE TO FETCH REPOSITORIES NODE\033[0m\r\n"))
		return
	}
	defer respRepos.Body.Close()
	var repos []GitHubRepo
	// FIX: Alimentando a struct com os dados reais trafegados no cabo
	json.NewDecoder(respRepos.Body).Decode(&repos)

	w.Write([]byte("\033[1;32m[✓] SECTOR SYNCHRONIZED! RENDERING DATAFEED...\033[0m\r\n\r\n"))
	w.Write([]byte("\033[1;36m┌── [ GITHUB REMOTE TELEMETRY ] ──────────────────────────────────────────────────┐\033[0m\r\n"))
	w.Write([]byte("  \033[1;35m  _(\\ _/)_ \033[0m    \033[1;32mOPERATIVE:\033[0m  " + user.Login + "\r\n"))
	w.Write([]byte("  \033[1;35m ((  \"  )) \033[0m    \033[1;32mNET REPOS:\033[0m  " + fmt.Sprintf("%d Public Units", user.PublicRepos) + "\r\n"))
	w.Write([]byte("  \033[1;35m  /\\-V-/\\  \033[0m    \033[1;32mFOLLOWERS:\033[0m  " + fmt.Sprintf("%d Active Nodes", user.Followers) + "\r\n"))
	w.Write([]byte("  \033[1;35m (___|___) \033[0m    \033[1;32mBIOGRAPHY:\033[0m  " + user.Bio + "\r\n"))
	w.Write([]byte("\033[1;36m├─────────────────────────────────────────────────────────────────────────────────┤\033[0m\r\n"))
	w.Write([]byte("  \033[1;33m📡 LIVE DEPLOYMENTS (RECENTLY UPDATED ON GITHUB):\033[0m\r\n"))

	for _, repo := range repos {
		desc := repo.Description
		if desc == "" {
			desc = "No description provided."
		}
		textoRepo := fmt.Sprintf("\r\n    \033[1;36m📁 %s\033[0m\r\n    ├─ \033[90mSYS_DESC:\033[0m %s\r\n    └─ \033[90mCORE_LNG:\033[0m \033[1;34m%s\033[0m  │  \033[1;33m⭐ %d\033[0m\r\n", repo.Name, desc, repo.Language, repo.Stargazers)
		Typewriter(w, textoRepo, 10*time.Millisecond)
	}
	w.Write([]byte("\033[1;36m└─────────────────────────────────────────────────────────────────────────────────┘\r\n"))
}

func Typewriter(w io.Writer, text string, delay time.Duration) {
	for _, char := range text {
		w.Write([]byte(string(char)))
		time.Sleep(delay)
	}
}

func DrawCyberBanner(w io.Writer) {
	green := "\033[1;92m"
	cyan  := "\033[1;96m"
	gray  := "\033[90m"
	reset := "\033[0m"
	w.Write([]byte(green + "  ____   _   _   _ ___ _____ _       ____ __  __ ____   _     ___ _   _ _____\r\n" + reset))
	w.Write([]byte(green + " |  _ \\ /_\\ | \\ | |_ _| ____| |     / ___|  \\/  |  _ \\  | |   |_ _| \\ | | ____|\r\n" + reset))
	w.Write([]byte(cyan  + " | | | / _ \\|  \\| | | ||  _| | |    | |   | |\\/| | | | | | |    | ||  \\| |  _|  \r\n" + reset))
	w.Write([]byte(cyan  + " | |_| / ___ \\ |\\  | | || |___| |___ | |___| |  | | |_| | | |___ | || |\\  | |___ \r\n" + reset))
	w.Write([]byte(cyan  + " |____/_/   \\_\\_| \\_|___|_____|_____| \\____|_|  |_|____/  |_____|___|_| \\_|_____|\r\n" + reset))
	w.Write([]byte(gray  + "  ─── [ HOST OVERRIDE: DANIEL_CMD_LINE ] ─── v2.0-RAW ─────────────────────────────\r\n\r\n" + reset))
}

func DrawStatusBar(w io.Writer, currentOption string) {
	w.Write([]byte("\r\n"))
	w.Write([]byte("\033[1;30;106m 🖥️  PORT: 2222 \033[1;37;45m ⚡ ENGINE: GO \033[1;37;40m ➔ ACTIVE: " + currentOption + " \033[0m\r\n"))
}
