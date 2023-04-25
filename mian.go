package main

import (
	"github.com/IBAX-io/ibax-cli/cmd"
	"runtime"
)

func main() {
	runtime.LockOSThread()
	cmd.Execute()
}
