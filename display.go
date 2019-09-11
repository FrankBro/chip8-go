package main

import (
	"github.com/nsf/termbox-go"
)

type Display interface {
	Clear(pixels []uint8)
	Draw(pixels []uint8) error
	WriteSprite(pixels []uint8, sprite []uint8, x, y uint8) bool
	Close()
}

type TermboxDisplay struct {
	fg, bg termbox.Attribute
}

func (display TermboxDisplay) Clear(pixels []uint8) {
	for y := 0; y < displayHeigh-1; y++ {
		for x := 0; x < displayWidth-1; x++ {
			index := y*displayWidth + x
			pixels[index] = 0
		}
	}
}

func (display TermboxDisplay) Draw(pixels []uint8) error {
	for y := 0; y < displayHeigh-1; y++ {
		for x := 0; x < displayWidth-1; x++ {
			index := y*displayWidth + x
			v := ' '

			if pixels[index] == 0x01 {
				v = '█'
			}

			termbox.SetCell(x, y, v, display.fg, display.bg)
		}
	}
	return termbox.Flush()
}

func (display TermboxDisplay) WriteSprite(pixels []uint8, sprite []uint8, x, y uint8) (collision bool) {
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
			if pixels[index] == 0x01 {
				collision = true
			}
			var v uint8
			if on {
				v = 0x01
			}
			pixels[index] ^= v
		}
	}
	return collision
}

func (display TermboxDisplay) Close() {
	termbox.Close()
}

func NewTermboxDisplay(fg, bg termbox.Attribute) (Display, error) {
	display := TermboxDisplay{
		fg: fg,
		bg: bg,
	}
	if err := termbox.Init(); err != nil {
		return display, err
	}

	termbox.HideCursor()

	err := termbox.Clear(bg, bg)
	if err != nil {
		return display, err
	}

	err = termbox.Flush()
	if err != nil {
		return display, err
	}

	return display, nil
}

func NewDefaultTermboxDisplay() (Display, error) {
	return NewTermboxDisplay(termbox.ColorDefault, termbox.ColorDefault)
}
