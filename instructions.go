package main

func (cpu *CPU) clear() {
	cpu.hardware.Clear()
}

func (cpu *CPU) jump(addr uint16) {
	cpu.pc = addr
}

func (cpu *CPU) call(addr uint16) {
	cpu.sp++
	cpu.stack[cpu.sp] = cpu.pc
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
	result := uint16(lhs) + uint16(rhs)
	var carry uint8
	if result > 255 {
		carry = 1
	}
	cpu.v[0xf] = carry
	cpu.v[register] = uint8(result)
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
	cpu.v[0xf] = value & 0x80
	cpu.v[register] = value << 1
}

func (cpu *CPU) assignI(value uint16) {
	cpu.i = value
}

func (cpu *CPU) draw(size, x, y uint8) error {
	var collision byte
	sprite := cpu.memory[cpu.i : cpu.i+uint16(size)]
	if cpu.hardware.WriteSprite(sprite, x, y) {
		collision = 1
	}
	cpu.v[0xF] = collision
	return cpu.hardware.Draw()
}

func (cpu *CPU) isKeySet(pos uint8) bool {
	return (cpu.hardware.GetKeys() & (1 << pos)) != 0
}

func log2n(n uint16) uint8 {
	if n > 1 {
		return 1 + log2n(n/2)
	}
	return 0
}

func isPowerOfTwo(n uint16) bool {
	return n&(^(n & (n - 1))) != 0
}

func findOnlySetBit(n uint16) (uint8, bool) {
	if !isPowerOfTwo(n) {
		return 0, false
	}
	return log2n(n) + 1, true
}

func (cpu *CPU) waitForKey() (key uint8) {
	// Maybe this could be done better
	// Make sure no key is pressed first
	var keys uint16
	for {
		keys = cpu.hardware.GetKeys()
		if keys == 0 {
			break
		}
	}
	for {
		keys = cpu.hardware.GetKeys()
		if bit, valid := findOnlySetBit(keys); valid {
			key = bit
			break
		}
	}
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
		cpu.v[i] = cpu.memory[cpu.i+i]
	}
}

func (cpu *CPU) fromMemory(length uint8) {
	var i uint16
	for ; i < uint16(length)+1; i++ {
		cpu.v[i] = cpu.memory[cpu.i+i]
	}
}
