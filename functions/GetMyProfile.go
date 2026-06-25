package functions

import (
    "portfolio-ssh/types"
)


// type Profile struct {
// 	Name           string
// 	Age            int
// 	Role           string
// 	Education      string
// 	Stack          []string //nolint:revive
// 	Certifications []string
// 	Projects       []ui.UIProject // FIX: Usando o tipo correto do pacote ui
// }

func GetMyProfile() types.Profile {
	p1 := types.Project{ // FIX: Instanciando como ui.UIProject
		Name:        "Seguramos",
		Description: "Plataforma full-stack de corretagem de seguros digital corporativa.",
		TechStack:   []string{"React", "Typescript", "Node.js", "PostgreSQL"},
	}

	p2 := types.Project{ // FIX: Instanciando como ui.UIProject
		Name:        "Go TUI Portfolio",
		Description: "Servidor SSH concorrente multiplataforma assíncrono e criptografado escrito do zero.",
		TechStack:   []string{"Go", "SSH Protocol", "RFC 4251", "Cryptography", "Linux Kernel"},
	}

	return types.Profile{
		Name:      "Daniel Caesar Mantilha",
		Age:       34,
		Role:      "Systems & Full-Stack Software Engineer // Network Engineer",
		Education: "Sistemas de Informação (Foco em Engenharia de Software)",
		Stack: []string{
			"GNU/Linux", "Golang", "Node.js", "TypeScript",
			"PostgreSQL", "Docker", "HTTP/Websockets", "REST APIs", "React/Next.js", "Python",
		},
		Certifications: []string{
			"CCNA (Cisco Certified Network Associate) - ID: Enterprise & Security Core",
			"CAE (Certificate in Advanced English) - University of Cambridge",
		},
		Projects: []types.Project{p1, p2},
	}
}