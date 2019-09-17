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
	hardware := NewDefaultSDLHardware()
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

	go func(cpu *CPU) {
		ticker := time.NewTicker(time.Second / 60)
		for range ticker.C {
			cpu.UpdateTimers()
		}
	}(cpu)
	ticker := time.NewTicker(time.Second / 600)
	for range ticker.C {
		if cpu.Quit() {
			break
		}
		err := cpu.cycle()
		if err != nil {
			panic(err)
		}
	}
}

// import "github.com/veandco/go-sdl2/sdl"

// func main() {
// 	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
// 		panic(err)
// 	}
// 	defer sdl.Quit()

// 	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
// 		800, 600, sdl.WINDOW_SHOWN)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer window.Destroy()

// 	surface, err := window.GetSurface()
// 	if err != nil {
// 		panic(err)
// 	}
// 	surface.FillRect(nil, 0)

// 	rect := sdl.Rect{0, 0, 200, 200}
// 	surface.FillRect(&rect, 0xffff0000)
// 	window.UpdateSurface()

// 	running := true
// 	for running {
// 		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
// 			switch event.(type) {
// 			case *sdl.QuitEvent:
// 				println("Quit")
// 				running = false
// 				break
// 			}
// 		}
// 	}
// }
