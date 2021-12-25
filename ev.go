package ev

import "reflect"

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

func ToEv(obj any) *Ev {
	value := reflect.ValueOf(obj)
	t := value.Type()
	ev := new(Ev)
	name := t.Name()
	ev.Name = name
	for i := 0; i < t.NumField(); i++ {
		name := t.Field(i).Name
		value := value.Field(i).Interface()
		ev.Attrs = append(ev.Attrs, Attr{
			Name:  name,
			Value: value,
		})
	}
	return ev
}
