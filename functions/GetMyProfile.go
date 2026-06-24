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
			"Linux/GNU", "Go (Golang)", "C Lang", "Node.js", "TypeScript",
			"PostgreSQL", "HTTP/Websockets", "REST APIs", "React/Next.js",
		},
		Certifications: []string{
			"CCNA (Cisco Certified Network Associate) - ID: Enterprise & Security Core",
			"CAE (Certificate in Advanced English) - University of Cambridge",
		},
		Projects: []types.Project{p1, p2},
	}
}