package main

import (
	"io/ioutil"
	"os"
	"time"
)

func main() {
	// scanner := bufio.NewScanner(os.Stdin)
	// scanner.Scan()
	name := "roms/test_opcode.ch8"
	if len(os.Args) == 2 {
		name = os.Args[1]
	}
	f, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	// out, err := os.Create("out")
	// if err != nil {
	// 	panic(err)
	// }
	// defer out.Close()
	out := ioutil.Discard
	hardware := NewDefaultTermboxHardware()
	err = hardware.Init()
	if err != nil {
		panic(err)
	}
	defer hardware.Close()
	cpu, err := NewCPU(hardware, out)
	if err != nil {
		panic(err)
	}
	program, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	err = cpu.LoadProgram(program)
	if err != nil {
		panic(err)
	}

	clock := time.NewTicker(time.Second / 60)
	for range clock.C {
		if hardware.Quit() {
			break
		}
		err := cpu.cycle()
		if err != nil {
			panic(err)
		}
	}
}
