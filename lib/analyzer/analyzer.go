package analyzer

type Analyzer interface {
	Walk(path string, recursive bool, callback func(codeBlock, filePath string, line int) error) error
}
