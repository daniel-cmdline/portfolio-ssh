package ui

import (
	"fmt"
	"net"
	"strings"
	"time"
)

// DrawAboutMe renderiza a tela "Sobre Mim" sem duplicações de stack
func DrawAboutMe(conn net.Conn, name string, age int, role string, education string, stack []string, certifications []string) {
	DrawCyberBanner(conn)
	conn.Write([]byte("\033[1;36m─── [ OVERRIDING SECURITY PROTOCOL // DECRYPTING OPERATIVE PROFILE ] ─────────────\033[0m\r\n\r\n"))

	manifesto := " \033[1;33m⚡ CORE SUMMARY:\033[0m\r\n" +
		"    Engenheiro de Software Full-Stack com foco cirúrgico em Infraestrutura, Redes\r\n" +
		"    e Segurança de Sistemas. Especialista em arquitetar ecossistemas robustos e\r\n" +
		"    escaláveis utilizando Node.js, TypeScript e React/Next.js, integrados a sockets\r\n" +
		"    de baixo nível, gerenciamento de processos e persistência avançada de dados.\r\n" +
		"    Diferenciado pela capacidade de transitar entre a escovação de bits na camada de\r\n" +
		"    rede e o tuning de performance no ecossistema web moderno.\r\n\r\n"
	Typewriter(conn, manifesto, 8*time.Millisecond)

	lineIdentity := fmt.Sprintf(" \033[1;32m➔ IDENTITY:\033[0m       %s\r\n", name)
	Typewriter(conn, lineIdentity, 10*time.Millisecond)

	lineRole := fmt.Sprintf(" \033[1;32m➔ FOCUS SPECIALTY:\033[0m %s (Security & Infrastructure Driven)\r\n", role)
	Typewriter(conn, lineRole, 10*time.Millisecond)

	lineGlobal := " \033[1;32m➔ GLOBAL MOBILITY:\033[0m  Vivência internacional (Canadá Node) │ Fluência profissional nativa.\r\n"
	Typewriter(conn, lineGlobal, 10*time.Millisecond)

	lineAcademics := fmt.Sprintf(" \033[1;32m➔ ACADEMICS:\033[0m       🎓 %s\r\n", education)
	Typewriter(conn, lineAcademics, 10*time.Millisecond)

	// Agora apenas unimos a fatia limpa do main.go sem inventar moda no fim da string
	stackStr := strings.Join(stack, " │ ")
	lineStack := fmt.Sprintf(" \033[1;32m➔ TECH STACK:\033[0m      🛠️  [ %s ]\r\n", stackStr)
	conn.Write([]byte(lineStack))

	conn.Write([]byte("\r\n \033[1;32m➔ VALIDATED CREDENTIALS:\033[0m\r\n"))
	
	cert1 := "🌐 CCNA (Cisco Certified Network Associate) — ID: Enterprise Routing, Switching & Security Master"
	textoCert1 := fmt.Sprintf("    ├─ %s\r\n", cert1)
	Typewriter(conn, textoCert1, 10*time.Millisecond)

	cert2 := "📜 Cambridge CAE (C1 Advanced) — University of Cambridge (International Native Fluency)"
	textoCert2 := fmt.Sprintf("    └─ %s\r\n", cert2)
	Typewriter(conn, textoCert2, 10*time.Millisecond)

	conn.Write([]byte("\r\n\033[1;36m──────────────────────────────────────────────────────────────────────────────────\033[0m\r\n"))
}
