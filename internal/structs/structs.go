// Package structs are some helpers functions extractef from
// https://github.com/fatih/structs
package structs

import (
	"reflect"
	"runtime"
	"strings"
)

func strctVal(s interface{}) reflect.Value {
	v := reflect.ValueOf(s)

	// if pointer get the underlying element
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	return v
}

// Name returns the name of a struct
func Name(s interface{}) string {
	return strctVal(s).Type().Name()
}

// Methods returns all of a struct methods
func Methods(s interface{}) (methods []reflect.Method) {
	t := reflect.ValueOf(s).Type()

	for i := 0; i < t.NumMethod(); i++ {
		methods = append(methods, t.Method(i))
	}

	return
}

// FuncName returns a function's name
func FuncName(f interface{}) string {
	v := reflect.ValueOf(f)

	if v.Kind() != reflect.Func {
		panic("Not a function")
	}

	p := runtime.FuncForPC(v.Pointer())
	if p == nil {
		return ""
	}

	trimedName := strings.TrimSuffix(p.Name(), "-fm")
	lastDotIndex := strings.LastIndex(trimedName, ".")
	return trimedName[lastDotIndex+1:]
}
