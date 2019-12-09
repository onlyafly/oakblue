package main

import (
	"fmt"
	"github.com/onlyafly/oakblue/internal/vm"
)

func main() {
	m := vm.NewMachine()
	m.LoadMemory([]byte{16, 33}, 0x3000)
	m.Execute()
	fmt.Println(m.RegisterDump())
}
