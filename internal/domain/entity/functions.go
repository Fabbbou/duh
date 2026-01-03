package entity

type Script struct {
	Name         string
	PathToFile   string
	Functions    []Function
	DataToInject string
	Warnings     []Warning
}

// Warnings detected while loading the script
type Warning struct {
	Line    int
	Details string
}

type Function struct {
	Name string
	// Extracted from comments, each line as a separate string
	Documentation []string
}
