package ev

type PutOp struct {
	Ev  *Ev
	Err chan error
}

func NewPutOp(ev *Ev) PutOp {
	return PutOp{
		Ev:  ev,
		Err: make(chan error, 1),
	}
}

func MustPut(sink Sink) func(ev *Ev) {
	return func(ev *Ev) {
		ce(sink.Put(ev))
	}
}
