.data
$0000 0x02
$0001 0x0a
$0002 0xff
$0003 0xf0

.code $0200
	movi 0xff R0
	tra R0
	addi 1
	int 1
	hlt

:overflow_int
	din
	movi 0xee R1
	str R1 $ffff
	ein
	crn

:stack_overflow_int
	din
	movi 0xe1 R1
	str R1 $ffff
	ein
	crn