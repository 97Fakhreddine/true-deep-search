package index

type Document struct {
	ID       string
	Path     string
	Title    string
	Content  string
	Metadata map[string]string
}
