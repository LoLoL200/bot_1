package events

type Type int

const (
	Unkdown Type = iota
	Message
)

type Event struct {
	Type Type
	Text string
	Meta interface{}
}

type Fetcher interface {
	Fetch(limit int) ([]Event, error)
}
type Processor interface {
	Process(e Event) error
}
