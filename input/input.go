package input

import (
	"bufio"
	"io"
	"strings"
)

// ReadLine captura a digitação byte a byte e faz o ECHO manual exigido pelo protocolo SSH
func ReadLine(w io.Writer, reader *bufio.Reader, prompt string) string {
	w.Write([]byte(prompt))
	
	// Ativa o cursor visível e piscante
	w.Write([]byte("\033[?25h"))
	
	var inputBuffer []byte

	for {
		b, err := reader.ReadByte()
		if err != nil {
			break
		}

		// Se for Enter (Carriage Return \r ou New Line \n), quebra o loop de digitação
		if b == 13 || b == 10 {
			w.Write([]byte("\r\n")) // Cospe uma quebra de linha pro cursor descer visualmente
			break
		}

		// Trata o Backspace (ASCII 127 ou 8) para o usuário conseguir apagar se errar
		if b == 127 || b == 8 {
			if len(inputBuffer) > 0 {
				inputBuffer = inputBuffer[:len(inputBuffer)-1]
				// Sequência ANSI: move o cursor 1 para trás, joga um espaço em branco (apaga), move 1 para trás de novo
				w.Write([]byte("\b \b"))
			}
			continue
		}

		// Adiciona o byte digitado ao nosso buffer interno
		inputBuffer = append(inputBuffer, b)

		// ECO: Escreve o caractere de volta pro canal para ele APARECER no terminal do usuário!
		w.Write([]byte{b})
	}
	
	// Desativa o cursor de novo para manter a estética TUI limpa
	w.Write([]byte("\033[?25l"))
	
	return strings.TrimSpace(string(inputBuffer))
}
