package function

type FunctionInfo struct {
	Name          string
	StartLine     uint
	EndLine       uint
	HasDocs       bool
	Documentation []string
}

type CodeOutsideFunction struct {
	Line    uint
	Content string
}
