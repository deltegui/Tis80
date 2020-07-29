.data
$4100 "user.rom"
$4120 "Fatal error: Stack overflow"

.code $0200
	; Register all interruptions
	movm overflow_int $0000
	movm stack_overflow_int $0002
	movm io_error $0004
	movm strcpy $0008

	; Call user code stored in $4100
	dsk $4100
	pmd
	cll $4100
	hlt

;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
;	ACC Overflow Interruption	;
;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
:overflow_int
	crn

;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
;	 Stack Overflow Interruption	;
;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
:stack_overflow_int
	din
	cfg 0x01
	movm $4120 $0100
	movm $3000 $0102
	cll strcpy
	hlt

;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
;    I/O Error Interruption    ;
;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
:io_error
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

;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
;      8 BIT INT TO STR      ;
;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
:itoa
	din
	; Copies the int to R0
	ldr $0100 R0

	; Copies the destiny direction where to write the str into $1000
	ldr $0102 R1
	str R1 $1000
	ldr $0103 R1
	str R1 $1001

	;Initializes registres where digits in decimal will be stored: (R10, R11, R12) = (HUNDREDS, TENS, UNITS)
	movi 0x00 R10
	movi 0x00 R11
	movi 0x00 R12

	; Start the algorithm. It will be substrating a number util the result is 0. Meanwhile R10, R11 and R12 will be updated.

	; R1 will be the counter
	movi 0x00 R1

	; R4 will have the limit of R10, R11, R12 (10)
	movi 10 R4

	; Here we substract R0 (the int we want to pass to ASCII) and the counter.
	; If it is 0 or we overflow flag is setted we are done
:itoa_loop
	tra R0
	sub R1
	jfg 0x00 itoa_end
	jeq itoa_end
	tra R1
	addi 1
	jfg 0x00 itoa_end
	tar R1

	; If we reach this part, counter is incremented, so will be R10, R11 and R12
	tra R12					; UNITS++
	addi 1
	tar R12
	sub R4					; UNITS - 10
	jeq itoa_inc_tens 		; IF UNITS == 10, we should set UNITS = 0 and TENS++
:itoa_return
	jmp itoa_loop

:itoa_inc_tens
	movi 0x00 R12			; set UNITS = 0
	tra R11					; TENS++
	addi 1
	tar R11
	sub R4 					; TENS - 10
	jeq itoa_inc_hundreds	; IF TENS == 10, we should set TENS = 0 and HUNDREDS++
	jmp itoa_return

:itoa_inc_hundreds
	movi 0x00 R11			; set TENS = 0
	tra R10					; HUNDREDS++
	addi 1
	tar R10
	jmp itoa_return

:itoa_end
	cfg 0x00
	; Write R10, R11 and R12 to the direction stored in $1000
	ldr $1000 R1
	ldr $1001 R2

	; Write HUNDREDS
	tra R10
	addi 0x30
	tar R10
	inw R10 $1000

	; Write TENS
	tra R2
	addi 1
	jfg 0x00 itoa_destiny_overflow_one
	str R2 $1001
:itoa_destiny_overflow_return_one
	tra R11
	addi 0x30
	tar R11
	inw R11 $1000

	; Write UNITS
	tra R2
	addi 1
	jfg 0x00 itoa_destiny_overflow_two
	str R2 $1001
:itoa_destiny_overflow_return_two
	tra R12
	addi 0x30
	tar R12
	inw R12 $1000

	; FINISH
	ein
	crn

:itoa_destiny_overflow_one
	cfg 0x00
	tra R1
	addi 1
	tar R1
	str R1 $1000
	movi 0x00 R2
	str R2 $1001
	jmp itoa_destiny_overflow_return_one

:itoa_destiny_overflow_two
	cfg 0x00
	tra R1
	addi 1
	tar R1
	str R1 $1000
	movi 0x00 R2
	str R2 $1001
	jmp itoa_destiny_overflow_return_two