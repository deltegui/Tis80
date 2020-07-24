.data
$0000 0x01
$0001 0x02
$0102 22
$0110 "example.rom"

.code $0200
	inr $0000 R0
	movi 0xff R15
	inw R15 $0000
	dsk $0110
	movi 0xff R0
	tra R0
	movi 0x00 R0
	movi 0x01 R1
	movi 0x02 R2
	movi 0x03 R3
	movi 0x04 R4
	movi 0x05 R5
	movi 0x06 R6
	movi 0x07 R7
	movi 0x08 R8
	movi 0x09 R9
	movi 0x0a R10
	movi 0x0b R11
	movi 0x0c R12
	movi 0x0d R13
	movi 0x0e R14
	movi 0x0f R16
	cll $4000
	movi 3 R3
	hlt
