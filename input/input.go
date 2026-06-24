package input

import (
	"bufio"
	"io"
	"strings"
)

// ReadLine lê strings de forma segura de qualquer io.ReadWriter (TCP ou SSH Channel)
func ReadLine(w io.Writer, reader *bufio.Reader, prompt string) string {
	w.Write([]byte(prompt))
	
	// Ativa o cursor visível e piscante pro usuário digitar
	w.Write([]byte("\033[?25h"))
	
	line, _ := reader.ReadString('\r')
	
	// Desativa o cursor de novo para manter a estética TUI limpa
	w.Write([]byte("\033[?25l"))
	
	return strings.TrimSpace(strings.ReplaceAll(line, "\r", ""))
}
