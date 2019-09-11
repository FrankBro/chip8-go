package main

import (
	"github.com/nsf/termbox-go"
)

type Keyboard interface {
	GetKey() (uint8, error)
}

type TermboxKeyboard struct {
	KeyMap map[rune]uint8
}

func (keyboard TermboxKeyboard) GetKey() (uint8, error) {
	for {

		event := termbox.PollEvent()

		if event.Type != termbox.EventKey {
			continue
		}

		if event.Ch == 0 {
			if event.Key == termbox.KeyEsc {
				return 0x0, ErrQuit
			}
			continue
		}

		if key, ok := keyboard.KeyMap[event.Ch]; ok {
			return key, nil
		}
	}
}

func NewDefaultTermboxKeyboard() Keyboard {
	keyboard := TermboxKeyboard{
		KeyMap: map[rune]uint8{
			'1': 0x1, '2': 0x2, '3': 0x3, '4': 0xC,
			'q': 0x4, 'w': 0x5, 'e': 0x6, 'r': 0xD,
			'a': 0x7, 's': 0x8, 'd': 0x9, 'f': 0xE,
			'z': 0xA, 'x': 0x0, 'c': 0xB, 'v': 0xF,
		},
	}
	return keyboard
}
