//go:build js && wasm

package platform

import "syscall/js"

func Browser() bool { return true }
func Args() []string {
	params := js.Global().Get("URLSearchParams").New(js.Global().Get("location").Get("search"))
	args := []string{}
	for _, name := range []string{"width", "height", "cell", "speed", "variance", "debug", "lines"} {
		if params.Call("has", name).Bool() {
			args = append(args, "-"+name+"="+params.Call("get", name).String())
		}
	}
	return args
}

var ready bool

// Ready reports the first actual game update, rather than just a completed download.
func Ready() {
	if !ready {
		ready = true
		callback := js.Global().Get("appReady")
		if callback.Type() == js.TypeFunction {
			callback.Invoke()
		}
	}
}
func Input() (bool, bool, bool) {
	input := js.Global().Get("appInput")
	if input.Type() != js.TypeObject {
		return false, false, false
	}
	left, right, jump := input.Get("left").Bool(), input.Get("right").Bool(), input.Get("jump").Bool()
	input.Set("jump", false)
	return left, right, jump
}
