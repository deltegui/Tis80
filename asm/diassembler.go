package tisasm

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

const (
	DataSectionByte byte = 0x00
	CodeSectoinByte      = 0x01
)

const (
	NumberType  byte = 0x02
	StringType       = 0x01
	SectionType      = 0x00
)

type Diassembler struct {
	binaryFile  *os.File
	currentLine int
	eof         bool
}

func NewDiassembler(file *os.File) Diassembler {
	return Diassembler{file, 0, false}
}

func (dasm Diassembler) Diasemble() {
	dasm.readSectionFlag()
	switch dasm.readByte() {
	case DataSectionByte:
		dasm.readDataSection()
	case CodeSectoinByte:
		dasm.readCodeSection()
	default:
		ShowError("Expected valid section after section flag")
	}
}

func (dasm Diassembler) readDataSection() {
	fmt.Println(".data")
	for !dasm.eof {
		dasm.readMemoryIgnoring()
		fmt.Print(" ")
		dataType := dasm.readByte()
		if dataType == SectionType {
			break
		} else if dataType == NumberType {
			dasm.readNumber()
		} else if dataType == StringType {
			dasm.readString()
		} else {
			ShowErrorf("Unkown data type: %x", dataType)
		}
		fmt.Println()
	}
	if dasm.eof {
		ShowError("Unexpected end of file in data section")
	}
	fmt.Println()
	dasm.readSectionFlag()
	dasm.expectBytes(0x01)
	dasm.readCodeSection()
}

func (dasm Diassembler) readCodeSection() {
	fmt.Print(".code ")
	dasm.startLineAnnotation()
	fmt.Println()
	for !dasm.eof {
		b := dasm.readByte()
		if dasm.eof {
			break
		}
		ins, err := GetInstructionUsingOpcode(b)
		if err != nil {
			ShowErrorf("%e", err)
		}
		fmt.Printf("%s ", ins.Literal)
		ins.Diassemble(dasm)
		fmt.Printf("   \t\t;")
		dasm.emitLineAnnotation()
		dasm.currentLine += ins.MemorySize
	}
}

func (dasm *Diassembler) startLineAnnotation() {
	high := dasm.readByte()
	low := dasm.readByte()
	str := fmt.Sprintf("%02x%02x", high, low)
	bytes, err := hex.DecodeString(str)
	if err != nil {
		ShowError("Malformed start code memory direction")
	}
	dasm.currentLine = int(binary.BigEndian.Uint16(bytes))
	fmt.Printf("$%s", str)
}

func (dasm Diassembler) emitLineAnnotation() {
	fmt.Print("$")
	fmt.Printf("%04x\n", dasm.currentLine)
}

func (dasm Diassembler) readSectionFlag() {
	dasm.expectBytes(0xff, 0xfe, 0xfe, 0xff)
}

func (dasm Diassembler) readMemoryIgnoring() {
	high := dasm.readByte()
	low := dasm.readByte()
	if high == 0x00 && low == 0x00 {
		return
	}
	fmt.Printf("$%02x%02x", high, low)
}

func (dasm Diassembler) readMemory() {
	high := dasm.readByte()
	low := dasm.readByte()
	fmt.Printf("$%02x%02x", high, low)
}

func (dasm Diassembler) readNumber() {
	fmt.Printf("%d", dasm.readByte())
}

func (dasm Diassembler) readRegister() {
	fmt.Print("R")
	dasm.readNumber()
}

func (dasm Diassembler) readString() {
	fmt.Print("\"")
	for {
		current := dasm.readByte()
		if current == 0x00 {
			break
		}
		fmt.Printf("%c", current)
	}
	fmt.Print("\"")
}

func (dasm Diassembler) expectBytes(bytes ...byte) {
	for i := 0; i < len(bytes); i++ {
		b := dasm.readByte()
		if b != bytes[i] {
			ShowErrorf("Expected %x, but %d byte is %x (not %x)\n", bytes, i+1, b, bytes[i])
		}
	}
}

func (dasm *Diassembler) readByte() byte {
	buffer := make([]byte, 1)
	length, err := dasm.binaryFile.Read(buffer)
	if err == io.EOF {
		dasm.eof = true
		return 0x00
	}
	if length != 1 || err != nil {
		ShowErrorf("Error while reading file %e\n", err)
	}
	return buffer[0]
}
