package ev

type Sink interface {
	Put(ev *Ev) error
}

type PutFunc func(ev *Ev) error

var _ Sink = PutFunc(nil)

func (p PutFunc) Put(ev *Ev) error {
	return p(ev)
}
