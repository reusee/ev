package ev

import "github.com/reusee/e4"

var (
	we = e4.Wrap.With(e4.WrapStacktrace)
	ce = e4.Check.With(e4.WrapStacktrace)
	he = e4.Handle
)
