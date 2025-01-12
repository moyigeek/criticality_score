package writer

import "fmt"

type StdOutWriter struct{}

func NewStdOutWriter() *StdOutWriter {
	return &StdOutWriter{}
}

func (w StdOutWriter) Open() error {
	return nil
}

func (w StdOutWriter) Close() error {
	return nil
}

func (w StdOutWriter) Write(url string) error {
	fmt.Println(url)
	return nil
}
