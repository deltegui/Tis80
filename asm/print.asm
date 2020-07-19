; rin es Read INdirection
; win es Write INdirection
; aisha death in vegas
; texas

.data
    $2FFE 30
    $2FFF 00
    $3000 "Hello world"

.code $4000
:print
    inr $2FFE R1
    tra R1
    jeq final
    tar R1
    str R1 $2000
    ldr $2FFF R1
    tra R1
    addi 1
    jfg 0 sumBigDir ; jfg Jump if FLag [mem] [flag number or instruction with flag name]
    tar R1
    str R1 $2FFF
    jmp print

:sumBigDir
    str R1 $2FFF
    ldr $2FFE R1
    tra R1
    addi 1
    tar R1
    str R1 $2FFE
    jmp print

:final
    hlt
