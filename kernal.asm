.data
$4100 "user.rom"

.code $0200
	; Register all interruptions
	movm overflow_int $0000
	movm stack_overflow_int $0002
	movm io_error $0004
	movm strcpy $0006

	; Call user code stored in $4100
	dsk $4100
	pmd
	cll $4100
	hlt

;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
;	ACC Overflow Interruption	;
;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
:overflow_int
	din
	movi 0xee R1
	str R1 $ffff
	ein
	crn

;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
;	 Stack Overflow Interruption	;
;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
:stack_overflow_int
	din
	movi 0xe1 R1
	str R1 $ffff
	ein
	crn

;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
;    I/O Error Interruption    ;
;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
:io_error
	din
	movi 0xe2 R2
	str R2 $fffd
	ein
	crn

;;;;;;;;;;;;;;;;;;;;;
;    String Copy    ;
;;;;;;;;;;;;;;;;;;;;;
:strcpy
	din
	; Copies a string stored in indirection $0100 to the direction stored in $0102

	; str source is stored in R0 R1 and $1000 $1001
	ldr $0100 R0
	ldr $0101 R1
	str R0 $1000
	str R1 $1001

	; str destiny is stored in R2 R3 and $1002 $1003
	ldr $0102 R2
	ldr $0103 R3
	str R2 $1002
	str R3 $1003

	; Iterate over str source until a 0x00 is readed
:loop_strcpy
	inr $1000 R5 						; if str[i] == 0 goto end_strcpy
	tra R5
	jeq end_strcpy
	inw R5 $1002						; else write str[i] to destiny[j]

	tra R1								; i++
	addi 1
	jfg 0x00 origin_overflow_strcpy 	; if lower part of direction have an overflow, fixit
	tar R1								; else store R1
	str R1 $1001

:origin_contiune_strcpy
	tra R3								; j++
	addi 1
	jfg 0x00 destiny_overflow_strcpy 	; if lower part of direction have an overflow, fixit
	tar R3
	str R3 $1003

:destiny_continue_strcpy
	jmp loop_strcpy

:end_strcpy
	ein
	crn

:origin_overflow_strcpy					; fix overflow of the lower part of the memory (origin str)
	cfg 0x00
	tra R0
	addi 1
	tar R0
	str R0 $1000
	movi 0x00 R1
	str R1 $1001
	jmp origin_contiune_strcpy

:destiny_overflow_strcpy				; fix overflow of the lower part of the memory (destiny str)
	cfg 0x00
	tra R2
	addi 1
	tar R2
	str R2 $1002
	movi 0x00 R3
	str R3 $1003
	jmp destiny_continue_strcpy

