package main

import (
	"bytes"
	"fmt"
	"io"
)

type CPU struct {
	hardware   Hardware
	memory     [memorySize]uint8
	display    [displaySize]uint8
	v          [registerCount]uint8
	stack      [stackSize]uint16
	log        io.Writer
	keys       uint16
	i          uint16
	dt         uint8
	st         uint8
	pc         uint16
	sp         uint8
	quit       bool
	shouldDraw bool
	opcode     opcode
}

func NewCPU(hardware Hardware, log io.Writer) (*CPU, error) {
	cpu := CPU{
		hardware: hardware,
		pc:       pcStart,
		log:      log,
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

func (cpu *CPU) Quit() bool {
	return cpu.quit
}

func (cpu *CPU) fetch() {
	cpu.opcode = opcode(cpu.memory[cpu.pc])<<8 | opcode(cpu.memory[cpu.pc+1])
}

func (cpu *CPU) UpdateTimers() {
	if cpu.dt > 0 {
		cpu.dt--
	}
	if cpu.st > 0 {
		cpu.st--
	}
}

func (cpu *CPU) execute() error {
	cmd := fmt.Sprintf("%X\n", cpu.opcode)
	_, err := cpu.log.Write([]byte(cmd))
	if err != nil {
		return err
	}
	cpu.pc += 2
	x := cpu.opcode.x()
	vx := cpu.v[x]
	y := cpu.opcode.y()
	vy := cpu.v[y]
	switch cpu.opcode & 0xF000 {
	case 0x0000:
		switch cpu.opcode {
		case 0x00E0:
			cpu.clear()
		case 0x00EE:
			cpu.jump(cpu.stack[cpu.sp])
			cpu.sp--
		default:
			return fmt.Errorf("CPU.execute: Unknown opcode: %X", cpu.opcode)
		}
	case 0x1000:
		cpu.jump(cpu.opcode.nnn())
	case 0x2000:
		cpu.call(cpu.opcode.nnn())
	case 0x3000:
		cpu.skip(vx == cpu.opcode.kk())
	case 0x4000:
		cpu.skip(vx != cpu.opcode.kk())
	case 0x5000:
		cpu.skip(vx == vy)
	case 0x6000:
		cpu.assignRegister(x, cpu.opcode.kk())
	case 0x7000:
		cpu.assignRegister(x, vx+cpu.opcode.kk())
	case 0x8000:
		switch cpu.opcode & 0x000F {
		case 0x0000:
			cpu.assignRegister(x, vy)
		case 0x0001:
			cpu.assignRegister(x, vx|vy)
		case 0x0002:
			cpu.assignRegister(x, vx&vy)
		case 0x0003:
			cpu.assignRegister(x, vx^vy)
		case 0x0004:
			cpu.add(x, vx, vy)
		case 0x0005:
			cpu.sub(x, vx, vy)
		case 0x0006:
			cpu.shr(x, vx)
		case 0x0007:
			cpu.sub(x, vy, vx)
		case 0x000E:
			cpu.shl(x, vx)
		default:
			return fmt.Errorf("CPU.execute: Unknown opcode: %X", cpu.opcode)
		}
	case 0x9000:
		switch cpu.opcode & 0x000F {
		case 0x0000:
			cpu.skip(vx != vy)
		default:
			return fmt.Errorf("CPU.execute: Unknown opcode: %X", cpu.opcode)
		}
	case 0xA000:
		cpu.assignI(cpu.opcode.nnn())
	case 0xB000:
		cpu.jump(cpu.opcode.nnn() + uint16(cpu.v[0]))
	case 0xC000:
		cpu.assignRegister(x, cpu.hardware.Int7()+cpu.opcode.kk())
	case 0xD000:
		cpu.draw(vx, vy, cpu.opcode.n())
	case 0xE000:
		switch cpu.opcode & 0x00FF {
		case 0x009E:
			cpu.skip(cpu.isKeySet(vx))
		case 0x00A1:
			cpu.skip(!cpu.isKeySet(vx))
		default:
			return fmt.Errorf("CPU.execute: Unknown opcode: %X", cpu.opcode)
		}
	case 0xF000:
		switch cpu.opcode & 0x00FF {
		case 0x0007:
			cpu.assignRegister(x, cpu.dt)
		case 0x000A:
			cpu.assignRegister(x, cpu.waitForKey())
		case 0x0015:
			cpu.assignDT(vx)
		case 0x0018:
			cpu.assignST(vx)
		case 0x001E:
			cpu.assignI(cpu.i + uint16(vx))
		case 0x0029:
			cpu.assignI(uint16(vx) * bytesPerSprite)
		case 0x0033:
			cpu.storeBcd(vx)
		case 0x0055:
			cpu.toMemory(x)
		case 0x0065:
			cpu.fromMemory(x)
		default:
			return fmt.Errorf("CPU.execute: Unknown opcode: %X", cpu.opcode)
		}
	}
	return nil
}

func (cpu *CPU) cycle() error {
	cpu.fetch()
	err := cpu.execute()
	if err != nil {
		return err
	}
	cpu.hardware.Update(&cpu.keys, &cpu.quit)
	if cpu.shouldDraw {
		cpu.shouldDraw = false
		err = cpu.hardware.Draw(cpu.display[:])
		if err != nil {
			return err
		}
	}
	// cpu.updateTimers()
	return nil
}
