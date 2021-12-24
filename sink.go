package ev

type Sink interface {
	Put(ev *Ev) error
}
