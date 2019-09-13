package main

import (
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	if len(os.Args) != 2 {
		log.Println("Usage: main rom")
	}
	name := os.Args[1]
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

	var quit bool
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, syscall.SIGINT)
	go func() {
		<-sig
		quit = true
	}()
	clock := time.NewTicker(time.Second / 60)
	for range clock.C {
		if quit {
			panic("quit")
		}
		err := cpu.cycle()
		if err != nil {
			panic(err)
		}
	}
}
