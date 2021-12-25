package ev

type EmitFunc = func(any) error

func NewEmitter(sink Sink) EmitFunc {
	return func(obj any) error {
		ev := ToEv(obj)
		return sink.Put(ev)
	}
}

func MustEmit(fn EmitFunc) func(any) {
	return func(obj any) {
		ce(fn(obj))
	}
}
