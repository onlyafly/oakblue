package main

import (
	"fmt"
	"github.com/onlyafly/oakblue/internal/vm"
)

func main() {
	m := vm.NewMachine()
	m.LoadBytecode([]byte{
		0x30, 0x00, // .ORIG 0x3000
		16, 33, // ADD R0 R0 1
		0b11110000, 0x25, // TRAP 0x25
	})
	err := m.Execute()
	if err != nil {
		fmt.Println("Error: " + err.Error())
	}
	fmt.Println(m.RegisterDump())
}
