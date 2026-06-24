package functions

import (
	"bufio"
	"fmt"
	"time"

	"golang.org/x/crypto/ssh"
	"portfolio-ssh/ui"
	"portfolio-ssh/utils"
)

// Essa função só roda internamente dentro do pacote
func handleGlobalRequests(requests <-chan *ssh.Request) {
	for req := range requests {
		if req.Type == "shell" || req.Type == "pty-req" {
			req.Reply(true, nil)
		} else {
			req.Reply(false, nil)
		}
	}
}

// HandleSSHChannel orquestra a TUI cyberpunk dentro da sessão criptografada do SSH
func HandleSSHChannel(ch ssh.Channel, requests <-chan *ssh.Request) {
	defer ch.Close()

	// Despacha os metadados do SSH em background (manter vivo/pty) de forma isolada
	go handleGlobalRequests(requests)

	// Roda a introdução cinematográfica da Matrix assim que o canal SSH abre
	ui.DrawMatrixIntro(ch)
	time.Sleep(1500 * time.Millisecond)

	profile := GetMyProfile()
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
				name := ReadLine(ch, reader, "         \033[1;33m│\033[0m  \033[1;32mFROM (NOME):\033[0m ", 0)
				ch.Write([]byte("         \033[1;33m│\033[0m                                                             \033[1;33m│\033[0m\r\n"))

				contact := ReadLine(ch, reader, "         \033[1;33m│\033[0m  \033[1;32mUP-LINK (EMAIL/LINK):\033[0m ", 0)
				ch.Write([]byte("         \033[1;33m│\033[0m                                                             \033[1;33m│\033[0m\r\n"))
				ch.Write([]byte("         \033[1;33m├─────────────────────────────────────────────────────────────┤\033[0m\r\n"))
				ch.Write([]byte("         \033[1;33m│\033[0m  \033[1;36mPAYLOAD (MENSAGEM):\033[0m                                         \033[1;33m│\033[0m\r\n"))

				msg := ReadLine(ch, reader, "         \033[1;33m│\033[0m  ➔  ", 50)
				ui.DrawCommsEnvelopeBottom(ch)

				ch.Write([]byte("\r\n\r\n  \033[1;31m[!] SEALING ENVELOPE... \033[0m\r\n"))
				time.Sleep(400 * time.Millisecond)
				ch.Write([]byte("  \033[1;36m[*] DISPATCHING PACKETS VIA SECURE HTTP POST... \033[0m"))

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
