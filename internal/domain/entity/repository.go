package entity

type Repository struct {
	Name    string
	Aliases map[string]string
	Exports map[string]string
}
