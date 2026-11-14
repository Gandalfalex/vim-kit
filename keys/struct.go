package keys

import (
	"reflect"
	"strings"
)

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

// NormalizeStruct rewrites every key in the []string fields of the struct v
// points to (nested structs included) to its canonical spelling. Call it on
// a config struct after loading, when the app also reads the fields directly.
func NormalizeStruct(v any) {
	normalize(reflect.Indirect(reflect.ValueOf(v)))
}

func normalize(v reflect.Value) {
	if v.Kind() != reflect.Struct {
		return
	}
	for i := range v.NumField() {
		fv := v.Field(i)
		switch {
		case !fv.CanSet():
		case fv.Kind() == reflect.Struct:
			normalize(fv)
		case fv.Type() == reflect.TypeOf([]string(nil)):
			for j := range fv.Len() {
				fv.Index(j).SetString(Normalize(fv.Index(j).String()))
			}
		}
	}
}
