package keys

import (
	"reflect"
	"strings"
)

// FromStruct builds Bindings from a config struct (or pointer to one) whose
// []string fields hold keys, naming each action after the field's tag
// (e.g. `yaml:"next_result"`), or the field name without one. Nested
// structs are flattened. This lets apps keep their own config layout.
func FromStruct(v any, tag string) Bindings {
	b := Bindings{}
	collect(reflect.Indirect(reflect.ValueOf(v)), tag, b)
	return b
}

func collect(v reflect.Value, tag string, b Bindings) {
	if v.Kind() != reflect.Struct {
		return
	}
	t := v.Type()
	for i := range t.NumField() {
		f := t.Field(i)
		if !f.IsExported() {
			continue
		}
		fv := v.Field(i)
		switch {
		case fv.Kind() == reflect.Struct:
			collect(fv, tag, b)
		case fv.Type() == reflect.TypeOf([]string(nil)):
			name, _, _ := strings.Cut(f.Tag.Get(tag), ",")
			if name == "-" {
				continue
			}
			if name == "" {
				name = f.Name
			}
			b[Action(name)] = append([]string(nil), fv.Interface().([]string)...)
		}
	}
}

// Only returns the bindings for the given actions, e.g. the subset of
// motions a view supports.
func (b Bindings) Only(actions ...Action) Bindings {
	o := Bindings{}
	for _, a := range actions {
		if ks, ok := b[a]; ok {
			o[a] = append([]string(nil), ks...)
		}
	}
	return o
}
