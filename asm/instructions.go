package tisasm

import "fmt"

type ParseParams func(parser Parser)
type Diassemble func(dasm Diassembler)

type Instruction struct {
	Literal     string
	OpCode      byte
	TokenSize   int
	MemorySize  int
	ParseParams ParseParams
	Diassemble  Diassemble
}

func paramsNone(prs Parser) {}

func paramsRegister(prs Parser) {
	prs.emitRegister(prs.scanner.Scan())
}

func paramsNumber(prs Parser) {
	prs.emitNumber(prs.scanner.Scan())
}

func paramsJump(prs Parser) {
	prs.emitJumpDest(prs.scanner.Scan())
}

func paramsJumpJump(prs Parser) {
	prs.emitJumpDest(prs.scanner.Scan())
	prs.emitJumpDest(prs.scanner.Scan())
}

func paramsJumpRegister(prs Parser) {
	prs.emitJumpDest(prs.scanner.Scan())
	prs.emitRegister(prs.scanner.Scan())
}

func paramsRegisterJump(prs Parser) {
	prs.emitRegister(prs.scanner.Scan())
	prs.emitJumpDest(prs.scanner.Scan())
}

func paramsNumberJump(prs Parser) {
	prs.emitNumber(prs.scanner.Scan())
	prs.emitJumpDest(prs.scanner.Scan())
}

func paramsMemoryRegister(prs Parser) {
	prs.emitMemory(prs.scanner.Scan())
	prs.emitRegister(prs.scanner.Scan())
}

func paramsRegisterMemory(prs Parser) {
	prs.emitRegister(prs.scanner.Scan())
	prs.emitMemory(prs.scanner.Scan())
}

func paramsRegisterRegister(prs Parser) {
	prs.emitRegister(prs.scanner.Scan())
	prs.emitRegister(prs.scanner.Scan())
}

func paramsNumberRegister(prs Parser) {
	prs.emitNumber(prs.scanner.Scan())
	prs.emitRegister(prs.scanner.Scan())
}

func diassembleNone(dasm Diassembler) {
}

func diassembleRegister(dasm Diassembler) {
	dasm.readRegister()
}

func diassembleNumber(dasm Diassembler) {
	dasm.readNumber()
}

func diassembleJump(dasm Diassembler) {
	dasm.readMemory()
}

func diassembleJumpJump(dasm Diassembler) {
	dasm.readMemory()
	fmt.Print(" ")
	dasm.readMemory()
}

func diassembleJumpRegister(dasm Diassembler) {
	dasm.readMemory()
	fmt.Print(" ")
	dasm.readRegister()
}

func diassembleRegisterJump(dasm Diassembler) {
	dasm.readRegister()
	fmt.Print(" ")
	dasm.readMemory()
}

func diassembleNumberJump(dasm Diassembler) {
	dasm.readNumber()
	fmt.Print(" ")
	dasm.readMemory()
}

func diassembleMemoryRegister(dasm Diassembler) {
	dasm.readMemory()
	fmt.Print(" ")
	dasm.readRegister()
}

func diassembleRegisterMemory(dasm Diassembler) {
	dasm.readRegister()
	fmt.Print(" ")
	dasm.readMemory()
}

func diassembleRegisterRegister(dasm Diassembler) {
	dasm.readRegister()
	fmt.Print(" ")
	dasm.readRegister()
}

func diassembleNumberRegister(dasm Diassembler) {
	dasm.readNumber()
	fmt.Print(" ")
	dasm.readRegister()
}

var instructions []Instruction = []Instruction{
	// Aritmetico-Logicos 0x0 y 0x1
	{
		Literal:     "add", // Add register. acc + Rx -> acc
		OpCode:      0x01,
		TokenSize:   2,
		MemorySize:  2,
		ParseParams: paramsRegister,
		Diassemble:  diassembleRegister,
	},
	{
		Literal:     "addi", // Add integer. acc + INT -> acc
		OpCode:      0x02,
		TokenSize:   2,
		MemorySize:  2,
		ParseParams: paramsNumber,
		Diassemble:  diassembleNumber,
	},
	{
		Literal:     "sub", // Substract register. acc - Rx -> acc
		OpCode:      0x03,
		TokenSize:   2,
		MemorySize:  2,
		ParseParams: paramsRegister,
		Diassemble:  diassembleRegister,
	},
	{
		Literal:     "subi", // Substract integer. acc - INT -> acc
		OpCode:      0x04,
		TokenSize:   2,
		MemorySize:  2,
		ParseParams: paramsNumber,
		Diassemble:  diassembleNumber,
	},
	{
		Literal:     "sil", // Shift left. acc << 1 -> acc
		OpCode:      0x05,
		TokenSize:   1,
		MemorySize:  1,
		ParseParams: paramsNone,
		Diassemble:  diassembleNone,
	},
	{
		Literal:     "sir", // Sift right. acc >> 1 -> acc
		OpCode:      0x06,
		TokenSize:   1,
		MemorySize:  1,
		ParseParams: paramsNone,
		Diassemble:  diassembleNone,
	},
	{
		Literal:     "and",
		OpCode:      0x07,
		TokenSize:   2,
		MemorySize:  2,
		ParseParams: paramsRegister,
		Diassemble:  diassembleRegister,
	},
	{
		Literal:     "or",
		OpCode:      0x08,
		TokenSize:   2,
		MemorySize:  2,
		ParseParams: paramsRegister,
		Diassemble:  diassembleRegister,
	},
	{
		Literal:     "not",
		OpCode:      0x09,
		TokenSize:   1,
		MemorySize:  1,
		ParseParams: paramsNone,
		Diassemble:  diassembleNone,
	},
	{
		Literal:     "xor", // eXclusive OR
		OpCode:      0x0a,
		TokenSize:   2,
		MemorySize:  2,
		ParseParams: paramsRegister,
		Diassemble:  diassembleRegister,
	},

	// Salto 0x2
	{
		Literal:     "jmp", // Inconditional jump
		OpCode:      0x20,
		TokenSize:   2,
		MemorySize:  3,
		ParseParams: paramsJump,
		Diassemble:  diassembleJump,
	},
	{
		Literal:     "jeq", // Jump equals. If acc == 0, jump to mem
		OpCode:      0x21,
		TokenSize:   2,
		MemorySize:  3,
		ParseParams: paramsJump,
		Diassemble:  diassembleJump,
	},
	{
		Literal:     "jne", // Jump not equal. If acc != 0, jump to mem
		OpCode:      0x22,
		TokenSize:   2,
		MemorySize:  3,
		ParseParams: paramsJump,
		Diassemble:  diassembleJump,
	},
	{
		Literal:     "jgt", // Jump Greater Than. If acc > 0, jump to mem
		OpCode:      0x23,
		TokenSize:   2,
		MemorySize:  3,
		ParseParams: paramsJump,
		Diassemble:  diassembleJump,
	},
	{
		Literal:     "jlt", // Jump Lower Than. If acc < 0, jump to mem
		OpCode:      0x24,
		TokenSize:   2,
		MemorySize:  3,
		ParseParams: paramsJump,
		Diassemble:  diassembleJump,
	},
	{
		Literal:     "jfg", // Jump If flag is setted. If Flags[INT], jump to mem
		OpCode:      0x25,
		TokenSize:   3,
		MemorySize:  4,
		ParseParams: paramsNumberJump,
		Diassemble:  diassembleNumberJump,
	},

	// Movimiento 0x3
	{
		Literal:     "ldr", // Load register. $mem -> Rx
		OpCode:      0x30,
		TokenSize:   3,
		MemorySize:  4,
		ParseParams: paramsMemoryRegister,
		Diassemble:  diassembleMemoryRegister,
	},
	{
		Literal:     "str", // Store register. Rx -> $mem
		OpCode:      0x31,
		TokenSize:   3,
		MemorySize:  4,
		ParseParams: paramsRegisterMemory,
		Diassemble:  diassembleRegisterMemory,
	},
	{
		Literal:     "mov", // Move. Rx -> Ry
		OpCode:      0x32,
		TokenSize:   3,
		MemorySize:  3,
		ParseParams: paramsRegisterRegister,
		Diassemble:  diassembleRegisterRegister,
	},
	{
		Literal:     "movi", // Move Integer. INT -> Rx
		OpCode:      0x33,
		TokenSize:   3,
		MemorySize:  3,
		ParseParams: paramsNumberRegister,
		Diassemble:  diassembleNumberRegister,
	},
	{
		Literal:     "tar", // Translate ACC to Rx. acc -> Rx
		OpCode:      0x34,
		TokenSize:   2,
		MemorySize:  2,
		ParseParams: paramsRegister,
		Diassemble:  diassembleRegister,
	},
	{
		Literal:     "tra", // Translate Rx to ACC. Rx -> ACC
		OpCode:      0x35,
		TokenSize:   2,
		MemorySize:  2,
		ParseParams: paramsRegister,
		Diassemble:  diassembleRegister,
	},
	{
		Literal:     "inr", // Read indirection. Reads the byte that points the memory stored in MEM
		OpCode:      0x36,
		TokenSize:   3,
		MemorySize:  4,
		ParseParams: paramsJumpRegister,
		Diassemble:  diassembleJumpRegister,
	},
	{
		Literal:     "inw", // Write indirection. Writes the byte that points the memory stored in MEM
		OpCode:      0x37,
		TokenSize:   3,
		MemorySize:  4,
		ParseParams: paramsRegisterJump,
		Diassemble:  diassembleRegisterJump,
	},
	{
		Literal:     "dsk", // Writes disk content into memory direction
		OpCode:      0x38,
		TokenSize:   2,
		MemorySize:  3,
		ParseParams: paramsJump,
		Diassemble:  diassembleJump,
	},
	{
		Literal:     "movm", // Write two bytes into destiny direction and the next one
		OpCode:      0x39,
		TokenSize:   3,
		MemorySize:  5,
		ParseParams: paramsJumpJump,
		Diassemble:  diassembleJumpJump,
	},

	// Llamadas 0x4
	{
		Literal:     "int", // Call interruption
		OpCode:      0x40,
		TokenSize:   2,
		MemorySize:  2,
		ParseParams: paramsNumber,
		Diassemble:  diassembleNumber,
	},
	{
		Literal:     "hlt", // Halt execution
		OpCode:      0x41,
		TokenSize:   1,
		MemorySize:  1,
		ParseParams: paramsNone,
		Diassemble:  diassembleNone,
	},
	{
		Literal:     "cll", // Call subrutine that starts from $mem
		OpCode:      0x42,
		TokenSize:   2,
		MemorySize:  3,
		ParseParams: paramsJump,
		Diassemble:  diassembleJump,
	},
	{
		Literal:     "crn", // Returns control to calling subrutine
		OpCode:      0x43,
		TokenSize:   1,
		MemorySize:  1,
		ParseParams: paramsNone,
		Diassemble:  diassembleNone,
	},
	{
		Literal:     "pmd", // Enable protected mode.
		OpCode:      0x44,
		TokenSize:   1,
		MemorySize:  1,
		ParseParams: paramsNone,
		Diassemble:  diassembleNone,
	},
	{
		Literal:     "ein", // Enable interrputions.
		OpCode:      0x45,
		TokenSize:   1,
		MemorySize:  1,
		ParseParams: paramsNone,
		Diassemble:  diassembleNone,
	},
	{
		Literal:     "din", // Disable protected mode.
		OpCode:      0x46,
		TokenSize:   1,
		MemorySize:  1,
		ParseParams: paramsNone,
		Diassemble:  diassembleNone,
	},
	{
		Literal:     "cfg", // Clear flag with number
		OpCode:      0x47,
		TokenSize:   2,
		MemorySize:  2,
		ParseParams: paramsNumber,
		Diassemble:  diassembleNumber,
	},

	// Stack manipulation 0x5
	{
		Literal:     "psa", // Push acc.
		OpCode:      0x50,
		TokenSize:   1,
		MemorySize:  1,
		ParseParams: paramsNone,
		Diassemble:  diassembleNone,
	},
	{
		Literal:     "poa", // Pop to acc.
		OpCode:      0x51,
		TokenSize:   1,
		MemorySize:  1,
		ParseParams: paramsNone,
		Diassemble:  diassembleNone,
	},
	{
		Literal:     "psr", // Push register Rx.
		OpCode:      0x52,
		TokenSize:   2,
		MemorySize:  2,
		ParseParams: paramsRegister,
		Diassemble:  diassembleRegister,
	},
	{
		Literal:     "por", // Pop to register Rx.
		OpCode:      0x53,
		TokenSize:   2,
		MemorySize:  2,
		ParseParams: paramsRegister,
		Diassemble:  diassembleRegister,
	},
}

func GetInstruction(str string) (Instruction, error) {
	return findInstruction(func(ins Instruction) bool {
		return ins.Literal == str
	}, "Undefined instruction %s", str)
}

func GetInstructionUsingOpcode(opcode byte) (Instruction, error) {
	return findInstruction(func(ins Instruction) bool {
		return ins.OpCode == opcode
	}, "Undefined instruction opcode %x", opcode)
}

func findInstruction(match func(Instruction) bool, msg string, params ...interface{}) (Instruction, error) {
	for _, ins := range instructions {
		if match(ins) {
			return ins, nil
		}
	}
	return Instruction{"", 0x00, 0, 0, paramsNone, diassembleNone}, fmt.Errorf(msg, params...)
}
