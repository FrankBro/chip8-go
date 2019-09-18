package main

func (cpu *CPU) clear() {
	for i := range cpu.display {
		cpu.display[i] = 0
	}
	cpu.shouldDraw = true
}

func (cpu *CPU) jump(addr uint16) {
	cpu.pc = addr
}

func (cpu *CPU) call(addr uint16) {
	cpu.stack[cpu.sp] = cpu.pc
	cpu.sp++
	cpu.pc = addr
}

func (cpu *CPU) skip(cond bool) {
	if cond {
		cpu.pc += 2
	}
}

func (cpu *CPU) assignRegister(register, value uint8) {
	cpu.v[register] = value
}

func (cpu *CPU) add(register, lhs, rhs uint8) {
	var carry uint8
	if lhs > 0xFF-rhs {
		carry = 1
	}
	cpu.v[0xf] = carry
	cpu.v[register] = lhs + rhs
}

func (cpu *CPU) sub(register, lhs, rhs uint8) {
	var borrow uint8
	if lhs > rhs {
		borrow = 1
	}
	cpu.v[0xf] = borrow
	cpu.v[register] = lhs - rhs
}

func (cpu *CPU) shr(register, value uint8) {
	cpu.v[0xf] = value & 1
	cpu.v[register] = value >> 1
}

func (cpu *CPU) shl(register, value uint8) {
	cpu.v[0xf] = value >> 7
	cpu.v[register] = value << 1
}

func (cpu *CPU) assignI(value uint16) {
	cpu.i = value
}

func (cpu *CPU) draw(x, y, height uint8) {
	var xl, yl uint16
	cpu.v[0xF] = 0
	for yl = 0; yl < uint16(height); yl++ {
		// yp := (uint16(y) + yl) % displayHeigh
		pixel := cpu.memory[cpu.i+yl]
		for xl = 0; xl < 8; xl++ {
			// xp := (uint16(x) + xl) % displayWidth
			if pixel&(0x80>>xl) != 0 {
				// index := xp + yp*displayWidth
				index := ((uint16(x) + xl) + (uint16(y)+yl)*displayWidth) % displaySize
				if cpu.display[index] == 1 {
					cpu.v[0xF] = 1
				}
				cpu.display[index] ^= 1
			}
		}
	}
	cpu.shouldDraw = true
}

func (cpu *CPU) isKeySet(pos uint8) bool {
	return (cpu.keys & (1 << pos)) != 0
}

func (cpu *CPU) waitForKey() (key uint8) {
	for key = 0; key < 16; key++ {
		if cpu.keys&(1<<key) != 0 {
			return key
		}
	}
	// Stay at the same instruction
	cpu.pc -= 2
	return key
}

func (cpu *CPU) assignDT(value uint8) {
	cpu.dt = value
}

func (cpu *CPU) assignST(value uint8) {
	cpu.st = value
}

func (cpu *CPU) storeBcd(value uint8) {
	cpu.memory[cpu.i] = value / 100
	cpu.memory[cpu.i+1] = (value / 10) % 10
	cpu.memory[cpu.i+2] = value % 10
}

func (cpu *CPU) toMemory(length uint8) {
	var i uint16
	for ; i < uint16(length)+1; i++ {
		// cpu.memory[cpu.i+i] = cpu.v[i]
		cpu.memory[cpu.i] = cpu.v[i]
		cpu.i++
	}
}

func (cpu *CPU) fromMemory(length uint8) {
	var i uint16
	for ; i < uint16(length)+1; i++ {
		// cpu.v[i] = cpu.memory[cpu.i+i]
		cpu.v[i] = cpu.memory[cpu.i]
		cpu.i++
	}
}
