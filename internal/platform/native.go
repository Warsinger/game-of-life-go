//go:build !js

package platform

import "os"

func Browser() bool             { return false }
func Args() []string            { return os.Args[1:] }
func Ready()                    {}
func Input() (bool, bool, bool) { return false, false, false }
