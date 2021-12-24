package ev

type Ev struct {
	Name       string
	Attrs      []Attr
	ExtraAttrs []Attr
	Subs       []*Ev
}

type Attr struct {
	Name  string
	Value any
}

type Evs = []*Ev
