package input

import (
	"bufio"
	"io"
	"strings"
)

// ReadLine captura o input e aceita um maxLen para quebrar a linha mantendo a borda do envelope
func ReadLine(w io.Writer, reader *bufio.Reader, prompt string, maxLen int) string {
	w.Write([]byte(prompt))
	
	// Ativa o cursor visível e piscante
	w.Write([]byte("\033[?25h"))
	
	var inputBuffer []byte
	currentLineLen := 0

	for {
		b, err := reader.ReadByte()
		if err != nil {
			break
		}

		// Se for Enter
		if b == 13 || b == 10 {
			w.Write([]byte("\r\n"))
			break
		}

		// Trata o Backspace
		if b == 127 || b == 8 {
			if len(inputBuffer) > 0 {
				inputBuffer = inputBuffer[:len(inputBuffer)-1]
				
				// Se apagou o primeiro caractere de uma linha quebrada, volta para a linha de cima
				if currentLineLen == 0 && maxLen > 0 {
					// Sequência ANSI: Sobe 1 linha (\033[A), vai pro fim da linha da caixa
					w.Write([]byte("\033[A\r         \033[1;33m│\033[0m  ➔  "))
					// Reposiciona o tamanho atual baseado no limite
					currentLineLen = maxLen
					// Simula o apagamento do caractere que estava no fim da linha anterior
					w.Write([]byte(strings.Repeat(" ", maxLen) + "\r         \033[1;33m│\033[0m  ➔  " + string(inputBuffer[len(inputBuffer)-maxLen:])))
				} else {
					currentLineLen--
					w.Write([]byte("\b \b"))
				}
			}
			continue
		}

		// Se tiver limite ativo (maxLen > 0) e estourar a largura do envelope
		if maxLen > 0 && currentLineLen >= maxLen {
			// 1. Fecha o visual da linha atual jogando espaços até a borda direita (se necessário) e fechando com │
			// Mas pro nosso layout simples, vamos só pular a linha mantendo a estética interna:
			w.Write([]byte("\r\n         \033[1;33m│\033[0m  "))
			currentLineLen = 0
		}

		// Adiciona o byte ao buffer
		inputBuffer = append(inputBuffer, b)
		currentLineLen++

		// ECO do byte na tela
		w.Write([]byte{b})
	}
	
	// Desativa o cursor
	w.Write([]byte("\033[?25l"))
	
	return strings.TrimSpace(string(inputBuffer))
}
