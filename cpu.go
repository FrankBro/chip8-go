package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
)

type CPU struct {
	hardware Hardware
	memory   [memorySize]uint8
	v        [registerCount]uint8
	i        uint16
	dt       uint8
	st       uint8
	pc       uint16
	sp       uint8
	stack    [stackSize]uint16
}

func NewCPU(hardware Hardware) (*CPU, error) {
	cpu := CPU{
		hardware: hardware,
		pc:       pcStart,
	}
	err := cpu.loadFont()
	return &cpu, err
}

func (cpu *CPU) load(offset int, r io.Reader) (int, error) {
	return r.Read(cpu.memory[offset:])
}

func (cpu *CPU) loadFont() error {
	font := font()
	_, err := cpu.load(0, bytes.NewReader(font))
	if err != nil {
		return fmt.Errorf("CPU.loadFont: CPU.load err: %s", err.Error())
	}
	return nil
}

func (cpu *CPU) LoadProgram(program []byte) error {
	_, err := cpu.load(pcStart, bytes.NewReader(program))
	if err != nil {
		return fmt.Errorf("CPU.LoadProgram: CPU.load err: %s", err.Error())
	}
	return nil
}

func (cpu *CPU) fetch() opcode {
	return opcode(cpu.memory[cpu.pc])<<8 | opcode(cpu.memory[cpu.pc+1])
}

func (cpu *CPU) updateTimers() {
	if cpu.dt > 0 {
		cpu.dt--
	}
	if cpu.st > 0 {
		cpu.st--
	}
}

func isKeySet(keys uint16, pos uint8) bool {
	return (keys & (1 << pos)) != 0
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

func (cpu *CPU) execute(opcode opcode) error {
	switch opcode & 0xF000 {
	case 0x0000:
		switch opcode {
		case 0x00E0:
			// 00E0
			cpu.hardware.Clear()
			cpu.pc += 2
		case 0x00EE:
			// 00EE
			cpu.pc = cpu.stack[cpu.sp]
			cpu.sp--
			cpu.pc += 2
		default:
			return fmt.Errorf("CPU.execute: Unknown opcode: %d", opcode)
		}
	case 0x1000:
		// 1nnn
		nnn := opcode.nnn()
		cpu.pc = nnn
	case 0x2000:
		// 2nnn
		nnn := opcode.nnn()
		cpu.sp++
		cpu.stack[cpu.sp] = cpu.pc
		cpu.pc = nnn
	case 0x3000:
		// 3xkk
		x := opcode.x()
		kk := opcode.kk()
		vx := cpu.v[x]
		cpu.pc += 2
		if vx == kk {
			cpu.pc += 2
		}
	case 0x4000:
		// 4xkk
		x := opcode.x()
		kk := opcode.kk()
		vx := cpu.v[x]
		cpu.pc += 2
		if vx != kk {
			cpu.pc += 2
		}
	case 0x5000:
		// 5xy0
		x := opcode.x()
		y := opcode.y()
		vx := cpu.v[x]
		vy := cpu.v[y]
		cpu.pc += 2
		if vx == vy {
			cpu.pc += 2
		}
	case 0x6000:
		// 6xkk
		x := opcode.x()
		kk := opcode.kk()
		cpu.v[x] = kk
		cpu.pc += 2
	case 0x7000:
		// 7xkk
		x := opcode.x()
		kk := opcode.kk()
		cpu.v[x] += kk
		cpu.pc += 2
	case 0x8000:
		switch opcode & 0x000F {
		case 0x0000:
			// 8xy0
			x := opcode.x()
			y := opcode.y()
			cpu.v[x] = cpu.v[y]
			cpu.pc += 2
		case 0x0001:
			// 8xy1
			x := opcode.x()
			y := opcode.y()
			cpu.v[x] |= cpu.v[y]
			cpu.pc += 2
		case 0x0002:
			// 8xy2
			x := opcode.x()
			y := opcode.y()
			cpu.v[x] &= cpu.v[y]
			cpu.pc += 2
		case 0x0003:
			// 8xy3
			x := opcode.x()
			y := opcode.y()
			cpu.v[x] ^= cpu.v[y]
			cpu.pc += 2
		case 0x0004:
			// 8xy4
			x := opcode.x()
			y := opcode.y()
			vx := cpu.v[x]
			vy := cpu.v[y]
			r := uint16(vx) + uint16(vy)
			var carry uint8
			if r > 255 {
				carry = 1
			}
			cpu.v[0xf] = carry
			cpu.v[x] = uint8(r)
			cpu.pc += 2
		case 0x0005:
			// 8xy5
			x := opcode.x()
			y := opcode.y()
			var carry uint8
			if cpu.v[x] > cpu.v[y] {
				carry = 1
			}
			cpu.v[0xf] = carry
			cpu.v[x] -= cpu.v[y]
			cpu.pc += 2
		case 0x0006:
			// 8xy6
			x := opcode.x()
			var carry uint8
			if cpu.v[x]&0x01 == 0x01 {
				carry = 1
			}
			cpu.v[0xf] = carry
			cpu.v[x] >>= 2
			cpu.pc += 2
		case 0x0007:
			// 8xy7
			x := opcode.x()
			y := opcode.y()
			var carry uint8
			if cpu.v[y] > cpu.v[x] {
				carry = 1
			}
			cpu.v[0xf] = carry
			cpu.v[x] = cpu.v[y] - cpu.v[x]
			cpu.pc += 2
		case 0x000E:
			// 8xyE
			x := opcode.x()
			var carry uint8
			if cpu.v[x]&0x80 == 0x80 {
				carry = 1
			}
			cpu.v[0xf] = carry
			cpu.v[x] <<= 2
			cpu.pc += 2
		default:
			return fmt.Errorf("CPU.execute: Unknown opcode: %d", opcode)
		}
	case 0x9000:
		switch opcode & 0x000F {
		case 0x0000:
			// 9xy0
			x := opcode.x()
			y := opcode.y()

			cpu.pc += 2
			if cpu.v[x] != cpu.v[y] {
				cpu.pc += 2
			}
		default:
			return fmt.Errorf("CPU.execute: Unknown opcode: %d", opcode)
		}
	case 0xA000:
		// Annn
		nnn := opcode.nnn()
		cpu.i = nnn
		cpu.pc += 2
	case 0xB000:
		// Bnnn
		nnn := opcode.nnn()
		v0 := cpu.v[0]
		cpu.pc = nnn + uint16(v0)
	case 0xC000:
		// Cxkk
		x := opcode.x()
		kk := opcode.kk()
		value := cpu.hardware.Int7()
		cpu.v[x] = value & kk
		cpu.pc += 2
	case 0xD000:
		// Dxyn
		x := opcode.x()
		y := opcode.y()
		n := opcode.n()
		vx := cpu.v[x]
		vy := cpu.v[y]

		var cf byte

		sprite := cpu.memory[cpu.i : cpu.i+uint16(n)]
		if cpu.hardware.WriteSprite(sprite, vx, vy) {
			cf = 0x01
		}

		cpu.v[0xF] = cf
		cpu.pc += 2

		err := cpu.hardware.Draw()
		if err != nil {
			return err
		}
	case 0xE000:
		switch opcode & 0x00FF {
		case 0x009E:
			// Ex9E
			x := opcode.x()
			keys := cpu.hardware.GetKeys()
			cpu.pc += 2
			if isKeySet(keys, cpu.v[x]) {
				cpu.pc += 2
			}
		case 0x00A1:
			// ExA1
			x := opcode.x()
			keys := cpu.hardware.GetKeys()
			cpu.pc += 2
			if !isKeySet(keys, cpu.v[x]) {
				cpu.pc += 2
			}
		default:
			return fmt.Errorf("CPU.execute: Unknown opcode: %d", opcode)
		}
	case 0xF000:
		switch opcode & 0x00FF {
		case 0x0007:
			// Fx07
			x := opcode.x()
			cpu.v[x] = cpu.dt
			cpu.pc += 2
		case 0x000A:
			// Fx0A
			x := opcode.x()
			// Maybe this could be done better
			// Make sure no key is pressed first
			var keys uint16
			for {
				keys = cpu.hardware.GetKeys()
				if keys == 0 {
					break
				}
			}
			var key uint8
			for {
				keys = cpu.hardware.GetKeys()
				if bit, valid := findOnlySetBit(keys); valid {
					key = bit
					break
				}
			}
			cpu.v[x] = key
			cpu.pc += 2
		case 0x0015:
			// Fx15
			x := opcode.x()
			vx := cpu.v[x]
			cpu.dt = vx
			cpu.pc += 2
		case 0x0018:
			// Fx18
			x := opcode.x()
			vx := cpu.v[x]
			cpu.st = vx
			cpu.pc += 2
		case 0x001E:
			// Fx1E
			x := opcode.x()
			vx := cpu.v[x]
			cpu.i += uint16(vx)
			cpu.pc += 2
		case 0x0029:
			// Fx29
			x := opcode.x()
			vx := cpu.v[x]
			cpu.i = uint16(vx) * bytesPerSprite
			cpu.pc += 2
		case 0x0033:
			// Fx33
			x := opcode.x()
			vx := cpu.v[x]
			cpu.memory[cpu.i] = vx / 100
			cpu.memory[cpu.i+1] = (vx / 10) % 10
			cpu.memory[cpu.i+2] = (vx % 100) % 10
			cpu.pc += 2
		case 0x0055:
			// Fx55
			x := uint16(opcode.x())
			var i uint16
			for ; i < x; i++ {
				cpu.memory[cpu.i+i] = cpu.v[i]
			}
			cpu.pc += 2
		case 0x0065:
			// Fx65
			x := uint16(opcode.x())
			var i uint16
			for ; i < x; i++ {
				cpu.v[i] = cpu.memory[cpu.i+i]
			}
			cpu.pc += 2
		default:
			return fmt.Errorf("CPU.execute: Unknown opcode: %d", opcode)
		}
	}
	return nil
}

func (cpu *CPU) cycle() error {
	opcode := cpu.fetch()
	err := cpu.execute(opcode)
	if err != nil {
		return err
	}
	cpu.updateTimers()
	return nil
}
