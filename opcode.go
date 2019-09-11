package main

type opcode uint16

func (op opcode) nnn() uint16 {
	return uint16(op & opcode(0x0FFF))
}

func (op opcode) n() uint8 {
	return uint8(op & opcode(0x000F))
}

func (op opcode) x() uint8 {
	return uint8((op & opcode(0x0F00)) >> opcode(8))
}

func (op opcode) y() uint8 {
	return uint8((op & opcode(0x00F0)) >> opcode(4))
}

func (op opcode) kk() uint8 {
	return uint8(op & opcode(0x00FF))
}
