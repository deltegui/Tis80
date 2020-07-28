.data
$5000 "DONE!"

.code $4100
    movi 0x00 R1
    tra R1
:loop
    addi 1
    jfg 0x00 end
    jmp loop

:end
    cfg 0x00
    movm $5000 $0100
	movm $3000 $0102
	int 4
	crn