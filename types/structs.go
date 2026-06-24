package types

type Project struct {
	Name        string
	Description string
	TechStack   []string
}
type Profile struct {
	Name           string
	Age            int
	Role           string
	Education      string
	Stack          []string //nolint:revive
	Certifications []string
	Projects       []Project // FIX: Usando o tipo correto do pacote ui
}
