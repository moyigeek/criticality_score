package writer

type Writer interface {
	Open() error
	Close() error
	Write(url string) error
}
