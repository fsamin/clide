package clide

//File represents a file on a cloud storage provider
type File struct {
	Filename string
	URL      string
}

//ProgressPrinter is a type used to print progression
type ProgressPrinter func(string, ...interface{}) (int, error)
