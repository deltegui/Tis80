;Ejemplo de un stack overflow y como el kernel lo maneja
.code $4100
    movi 0x00 R1
    tra R1
:loop
    psa
    jfg 0x01 end
    addi 1
    jfg 0x00 set_acc
:continue
    jmp loop

:set_acc
    cfg 0x00
    tra R1
    jmp continue

:end
    cfg 0x01
	crn