package main

import (
	"plugin"
	_ "runtime/cgo"
)

func main() {
	p, err := plugin.Open("plugin.so")
	if err != nil {
		panic(err)
	}
	v, err := p.Lookup("PublicVar")
	if err != nil {
		panic(err)
	}
	f, err := p.Lookup("PublicFunc")
	if err != nil {
		panic(err)
	}
	*v.(*int) = 7
	f.(func())()
}
