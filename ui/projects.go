package ui

import (
	"fmt"
	"net"
	"strings"
	"time"
)

// UIProject é a estrutura que o pacote UI usa para renderizar os projetos
type UIProject struct {
	Name        string
	Description string
	TechStack   []string
}

// DrawProjects renderiza a lista de repositórios sem molduras travadas
func DrawProjects(conn net.Conn, projects []UIProject) {
	DrawCyberBanner(conn)
	conn.Write([]byte("\033[1;35m─── [ FETCHING REPOSITORIES // LOCAL DATABASE ] ──────────────────────────────────\033[0m\r\n"))
	
	for _, proj := range projects {
		textoProj := fmt.Sprintf("\r\n \033[1;33m⚡ TARGET:\033[0m %s\r\n    ├─ \033[90mDESC:\033[0m %s\r\n", proj.Name, proj.Description)
		Typewriter(conn, textoProj, 15*time.Millisecond)

		stack := strings.Join(proj.TechStack, " │ ")
		conn.Write([]byte(fmt.Sprintf("    └─ \033[90mSTACK:\033[0m \033[1;34m[ %s ]\033[0m\r\n", stack)))
	}
	conn.Write([]byte("\r\n\033[1;35m──────────────────────────────────────────────────────────────────────────────────\033[0m\r\n"))
}
