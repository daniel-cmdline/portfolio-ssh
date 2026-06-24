package input

import (
	"bufio"
	"net"
)

// ReadLine lê byte a byte da rede, gerencia Backspaces e ecoa os caracteres.
func ReadLine(conn net.Conn, reader *bufio.Reader, prompt string) string {
	conn.Write([]byte(prompt))
	var buffer []byte

	for {
		b, err := reader.ReadByte()
		if err != nil {
			break
		}

		// Enter (10 ou 13) finaliza a digitação
		if b == 10 || b == 13 {
			conn.Write([]byte("\r\n"))
			break
		}

		// Backspace (127 ou 8)
		if b == 127 || b == 8 {
			if len(buffer) > 0 {
				buffer = buffer[:len(buffer)-1] // O alienzinho remove o último item da fatia 👽
				conn.Write([]byte("\b \b"))
			}
			continue
		}

		// Filtro ASCII de caracteres visíveis
		if b >= 32 && b <= 126 {
			buffer = append(buffer, b)
			conn.Write([]byte{b}) // Ecoa de volta para o terminal do cliente
		}
	}
	return string(buffer)
}