package main

type Config struct {
	Random   Random
	Keyboard Keyboard
	Display  Display
}

func NewDefaultTermboxConfig() (config Config, err error) {
	display, err := NewDefaultTermboxDisplay()
	if err != nil {
		return config, err
	}
	config = Config{
		Random:   NewGoRandom(),
		Keyboard: NewDefaultTermboxKeyboard(),
		Display:  display,
	}
	return config, nil
}
