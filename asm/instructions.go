package tisasm

import "fmt"

type ParseParams func(parser Parser)

type Instruction struct {
	Literal     string
	OpCode      byte
	TokenSize   int
	MemorySize  int
	ParseParams ParseParams
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

var instructions []Instruction = []Instruction{
	// Aritmetico-Logicos 0x0 y 0x1
	{
		Literal:     "add", // Add register. acc + Rx -> acc
		OpCode:      0x00,
		TokenSize:   2,
		MemorySize:  2,
		ParseParams: paramsRegister,
	},
	{
		Literal:     "addi", // Add integer. acc + INT -> acc
		OpCode:      0x01,
		TokenSize:   2,
		MemorySize:  2,
		ParseParams: paramsNumber,
	},
	{
		Literal:     "sub", // Substract register. acc - Rx -> acc
		OpCode:      0x02,
		TokenSize:   2,
		MemorySize:  2,
		ParseParams: paramsRegister,
	},
	{
		Literal:     "subi", // Substract integer. acc - INT -> acc
		OpCode:      0x03,
		TokenSize:   2,
		MemorySize:  2,
		ParseParams: paramsNumber,
	},
	{
		Literal:     "sil", // Shift left. acc << 1 -> acc
		OpCode:      0x04,
		TokenSize:   1,
		MemorySize:  1,
		ParseParams: paramsNone,
	},
	{
		Literal:     "sir", // Sift right. acc >> 1 -> acc
		OpCode:      0x05,
		TokenSize:   1,
		MemorySize:  1,
		ParseParams: paramsNone,
	},
	{
		Literal:     "and",
		OpCode:      0x06,
		TokenSize:   2,
		MemorySize:  2,
		ParseParams: paramsRegister,
	},
	{
		Literal:     "or",
		OpCode:      0x07,
		TokenSize:   2,
		MemorySize:  2,
		ParseParams: paramsRegister,
	},
	{
		Literal:     "not",
		OpCode:      0x08,
		TokenSize:   1,
		MemorySize:  1,
		ParseParams: paramsNone,
	},
	{
		Literal:     "xor", // eXclusive OR
		OpCode:      0x09,
		TokenSize:   2,
		MemorySize:  2,
		ParseParams: paramsRegister,
	},

	// Salto 0x2
	{
		Literal:     "jmp", // Inconditional jump
		OpCode:      0x20,
		TokenSize:   2,
		MemorySize:  3,
		ParseParams: paramsJump,
	},
	{
		Literal:     "jeq", // Jump equals. If acc == 0, jump to mem
		OpCode:      0x21,
		TokenSize:   2,
		MemorySize:  3,
		ParseParams: paramsJump,
	},
	{
		Literal:     "jne", // Jump not equal. If acc != 0, jump to mem
		OpCode:      0x22,
		TokenSize:   2,
		MemorySize:  3,
		ParseParams: paramsJump,
	},
	{
		Literal:     "jgt", // Jump Greater Than. If acc > 0, jump to mem
		OpCode:      0x23,
		TokenSize:   2,
		MemorySize:  3,
		ParseParams: paramsJump,
	},
	{
		Literal:     "jlt", // Jump Lower Than. If acc < 0, jump to mem
		OpCode:      0x24,
		TokenSize:   2,
		MemorySize:  3,
		ParseParams: paramsJump,
	},
	{
		Literal:     "jfg", // Jump If flag is setted. If Flags[INT], jump to mem
		OpCode:      0x25,
		TokenSize:   3,
		MemorySize:  4,
		ParseParams: paramsNumberJump,
	},

	// Movimiento 0x3
	{
		Literal:     "ldr", // Load register. $mem -> Rx
		OpCode:      0x30,
		TokenSize:   3,
		MemorySize:  4,
		ParseParams: paramsMemoryRegister,
	},
	{
		Literal:     "str", // Store register. Rx -> $mem
		OpCode:      0x31,
		TokenSize:   3,
		MemorySize:  4,
		ParseParams: paramsRegisterMemory,
	},
	{
		Literal:     "mov", // Move. Rx -> Ry
		OpCode:      0x32,
		TokenSize:   3,
		MemorySize:  3,
		ParseParams: paramsRegisterRegister,
	},
	{
		Literal:     "movi", // Move Integer. INT -> Rx
		OpCode:      0x33,
		TokenSize:   3,
		MemorySize:  3,
		ParseParams: paramsNumberRegister,
	},
	{
		Literal:     "tar", // Translate ACC to Rx. acc -> Rx
		OpCode:      0x34,
		TokenSize:   2,
		MemorySize:  2,
		ParseParams: paramsRegister,
	},
	{
		Literal:     "tra", // Translate Rx to ACC. Rx -> ACC
		OpCode:      0x35,
		TokenSize:   2,
		MemorySize:  2,
		ParseParams: paramsRegister,
	},
	{
		Literal:     "inr", // Read indirection. Reads the byte that points the memory stored in MEM
		OpCode:      0x36,
		TokenSize:   2,
		MemorySize:  3,
		ParseParams: paramsJumpRegister,
	},
	{
		Literal:     "inw", // Write indirection. Writes the byte that points the memory stored in MEM
		OpCode:      0x35,
		TokenSize:   2,
		MemorySize:  3,
		ParseParams: paramsRegisterJump,
	},
	{
		Literal:     "dsk", // Load disk from memory direction
		OpCode:      0x36,
		TokenSize:   2,
		MemorySize:  3,
		ParseParams: paramsJump,
	},

	// Llamadas 0x4
	{
		Literal:     "int", // Call interruption
		OpCode:      0x40,
		TokenSize:   2,
		MemorySize:  2,
		ParseParams: paramsNumber,
	},
	{
		Literal:     "hlt", // Halt execution
		OpCode:      0x41,
		TokenSize:   1,
		MemorySize:  1,
		ParseParams: paramsNone,
	},
	{
		Literal:     "cll", // Call subrutine that starts from $mem
		OpCode:      0x42,
		TokenSize:   2,
		MemorySize:  3,
		ParseParams: paramsJump,
	},
	{
		Literal:     "crn", // Returns control to calling subrutine
		OpCode:      0x43,
		TokenSize:   1,
		MemorySize:  1,
		ParseParams: paramsNone,
	},
	{
		Literal:     "pmd", // Enable protected mode.
		OpCode:      0x44,
		TokenSize:   1,
		MemorySize:  1,
		ParseParams: paramsNone,
	},
	{
		Literal:     "ein", // Enable interrputions.
		OpCode:      0x45,
		TokenSize:   1,
		MemorySize:  1,
		ParseParams: paramsNone,
	},
	{
		Literal:     "din", // Disable protected mode.
		OpCode:      0x46,
		TokenSize:   1,
		MemorySize:  1,
		ParseParams: paramsNone,
	},

	// Stack manipulation 0x5
	{
		Literal:     "psa", // Push acc.
		OpCode:      0x50,
		TokenSize:   1,
		MemorySize:  1,
		ParseParams: paramsNone,
	},
	{
		Literal:     "poa", // Pop to acc.
		OpCode:      0x51,
		TokenSize:   1,
		MemorySize:  1,
		ParseParams: paramsNone,
	},
	{
		Literal:     "psr", // Push register Rx.
		OpCode:      0x52,
		TokenSize:   1,
		MemorySize:  1,
		ParseParams: paramsNone,
	},
	{
		Literal:     "por", // Pop to register Rx.
		OpCode:      0x53,
		TokenSize:   1,
		MemorySize:  1,
		ParseParams: paramsNone,
	},
}

func GetInstruction(str string) (Instruction, error) {
	for _, ins := range instructions {
		if ins.Literal == str {
			return ins, nil
		}
	}
	return Instruction{"", 0x00, 0, 0, paramsNone}, fmt.Errorf("Undefined instruction %s", str)
}
