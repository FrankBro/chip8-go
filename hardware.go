package main

import (
	"math/rand"
	"time"

	"github.com/nsf/termbox-go"
)

type Hardware interface {
	// Setup
	Init() error
	Close()
	// Random
	Int7() uint8
	// Keypad
	GetKeys() uint16
	// Display
	GetPixels() []uint8
	Draw() error
	WriteSprite(sprite []uint8, x, y uint8) bool
	Clear()
}

type TermboxHardware struct {
	rand   *rand.Rand
	fg, bg termbox.Attribute
	keyMap map[rune]uint8
	pixels [displaySize]uint8
	keys   uint16
}

func (hardware *TermboxHardware) Init() error {
	// Random
	source := rand.NewSource(time.Now().Unix())
	rand := rand.New(source)
	hardware.rand = rand
	// Display
	if err := termbox.Init(); err != nil {
		return err
	}

	termbox.HideCursor()

	err := termbox.Clear(hardware.bg, hardware.bg)
	if err != nil {
		return err
	}

	err = termbox.Flush()
	if err != nil {
		return err
	}

	go func(hardware *TermboxHardware) {
		for {
			event := termbox.PollEvent()
			if event.Type != termbox.EventKey {
				continue
			}
			if event.Ch == 0 {
				if event.Key == termbox.KeyEsc {
					panic("esc")
				}
				continue
			}
			if key, ok := hardware.keyMap[event.Ch]; ok {
				hardware.keys |= 1 << key
			}
		}
	}(hardware)

	return nil
}

func (hardware *TermboxHardware) Close() {
	termbox.Close()
}

func (hardware *TermboxHardware) Int7() uint8 {
	value := hardware.rand.Intn(256)
	return uint8(value)
}

func (hardware *TermboxHardware) GetKeys() uint16 {
	return hardware.keys
}

func (hardware *TermboxHardware) GetPixels() []uint8 {
	return hardware.pixels[:]
}

func (hardware *TermboxHardware) Draw() error {
	for y := 0; y < displayHeigh-1; y++ {
		for x := 0; x < displayWidth-1; x++ {
			index := y*displayWidth + x
			v := ' '

			if hardware.pixels[index] == 0x01 {
				v = '█'
			}

			termbox.SetCell(x, y, v, hardware.fg, hardware.bg)
		}
	}
	return termbox.Flush()
}

func (hardware *TermboxHardware) WriteSprite(sprite []uint8, x, y uint8) (collision bool) {
	n := len(sprite)
	for yl := 0; yl < n; yl++ {
		r := sprite[yl]

		for xl := 0; xl < 8; xl++ {
			i := 0x80 >> uint8(xl)
			on := (r & byte(i)) == byte(i)
			xp := uint16(x) + uint16(xl)
			if xp >= displayWidth {
				xp -= displayWidth
			}
			yp := uint16(y) + uint16(yl)
			if yp >= displayHeigh {
				yp -= displayHeigh
			}

			index := xp + yp*displayWidth
			if hardware.pixels[index] == 0x01 {
				collision = true
			}
			var v uint8
			if on {
				v = 0x01
			}
			hardware.pixels[index] ^= v
		}
	}
	return collision
}

func (hardware *TermboxHardware) Clear() {
	for y := 0; y < displayHeigh-1; y++ {
		for x := 0; x < displayWidth-1; x++ {
			index := y*displayWidth + x
			hardware.pixels[index] = 0
		}
	}
}

func NewTermboxHardware(fg, bg termbox.Attribute, keyMap map[rune]uint8) Hardware {
	hardware := TermboxHardware{
		fg:     fg,
		bg:     bg,
		keyMap: keyMap,
	}
	return &hardware
}

func NewDefaultTermboxHardware() Hardware {
	keyMap := map[rune]uint8{
		'1': 0x1, '2': 0x2, '3': 0x3, '4': 0xC,
		'q': 0x4, 'w': 0x5, 'e': 0x6, 'r': 0xD,
		'a': 0x7, 's': 0x8, 'd': 0x9, 'f': 0xE,
		'z': 0xA, 'x': 0x0, 'c': 0xB, 'v': 0xF,
	}
	return NewTermboxHardware(termbox.ColorDefault, termbox.ColorDefault, keyMap)
}
