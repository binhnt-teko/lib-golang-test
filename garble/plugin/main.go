package main

import (
	_ "runtime/cgo"

	"github.com/blcvn/lib-golang-test/garble/plugin/lib"
)

//export PublicVar
var PublicVar int = lib.ImportedFunc()

//export privateFunc
func privateFunc(n int) { println("Hello, number", n) }

//export PublicFunc
func PublicFunc() { privateFunc(PublicVar) }
