package main

import (
	"math/rand"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

type Hardware interface {
	Init() error
	Close()
	Int7() uint8
	Update(keys *uint16, quit *bool)
	Draw(pixels []uint8) error
}

type SDLHardware struct {
	window   *sdl.Window
	renderer *sdl.Renderer
	rand     *rand.Rand
	keyMap   map[sdl.Keycode]uint8
}

const size = 10

func (hardware *SDLHardware) Init() error {
	// Random
	source := rand.NewSource(time.Now().Unix())
	rand := rand.New(source)
	hardware.rand = rand
	// Display
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		return err
	}
	window, err := sdl.CreateWindow("chip8", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		displayWidth*10, displayHeigh*10, sdl.WINDOW_SHOWN)
	if err != nil {
		return err
	}
	renderer, err := sdl.CreateRenderer(window, -1, 0)
	if err != nil {
		return err
	}
	hardware.window = window
	hardware.renderer = renderer
	return nil
}

func (hardware *SDLHardware) Close() {
	_ = hardware.renderer.Destroy()
	_ = hardware.window.Destroy()
}

func (hardware *SDLHardware) Int7() uint8 {
	value := hardware.rand.Intn(256)
	return uint8(value)
}

func (hardware *SDLHardware) Update(keys *uint16, quit *bool) {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch et := event.(type) {
		case *sdl.QuitEvent:
			*quit = true
		case *sdl.KeyboardEvent:
			if et.Type == sdl.KEYUP {
				if key, ok := hardware.keyMap[et.Keysym.Sym]; ok {
					*keys &= 0 << key
				}
			} else if et.Type == sdl.KEYDOWN {
				if et.Keysym.Sym == sdl.K_ESCAPE {
					*quit = true
				} else if key, ok := hardware.keyMap[et.Keysym.Sym]; ok {
					*keys |= 1 << key
				}
			}
		}
	}
}

func (hardware *SDLHardware) Draw(pixels []uint8) error {
	err := hardware.renderer.SetDrawColor(0, 0, 0, 255)
	if err != nil {
		return err
	}
	err = hardware.renderer.Clear()
	if err != nil {
		return err
	}

	var x, y int32
	for x = 0; x < displayWidth; x++ {
		for y = 0; y < displayHeigh; y++ {
			i := x + y*displayWidth
			pixel := pixels[i]
			if pixel == 0 {
				err = hardware.renderer.SetDrawColor(0, 0, 0, 255)
				if err != nil {
					return err
				}
			} else {
				err = hardware.renderer.SetDrawColor(255, 255, 255, 255)
				if err != nil {
					return err
				}
			}
			err = hardware.renderer.FillRect(&sdl.Rect{
				X: x * size, Y: y * size, W: size, H: size,
			})
			if err != nil {
				return err
			}
		}
	}

	hardware.renderer.Present()
	return nil
}

func NewSDLHardware(keyMap map[sdl.Keycode]uint8) Hardware {
	hardware := SDLHardware{
		keyMap: keyMap,
	}
	return &hardware
}

func NewDefaultSDLHardware() Hardware {
	keyMap := map[sdl.Keycode]uint8{
		sdl.K_1: 0x1, sdl.K_2: 0x2, sdl.K_3: 0x3, sdl.K_4: 0xC,
		sdl.K_q: 0x4, sdl.K_w: 0x5, sdl.K_e: 0x6, sdl.K_r: 0xD,
		sdl.K_a: 0x7, sdl.K_s: 0x8, sdl.K_d: 0x9, sdl.K_f: 0xE,
		sdl.K_z: 0xA, sdl.K_x: 0x0, sdl.K_c: 0xB, sdl.K_v: 0xF,
	}
	return NewSDLHardware(keyMap)
}
