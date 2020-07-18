; rin es Read INdirection
; win es Write INdirection

.data
    $2FFE 30
    $2FFF 00
    $3000 "Hello world"

.code $4000
:print
    rin $2FFE R1
    tra R1
    jeq final
    tar R1
    sti R1 $2000
    ldi $2FFF R1
    tra R1
    addi 1
    jfg sumBigDir 0     ; jfg Jump if FLag [mem] [flag number or instruction with flag name]
    tar R1
    sti R1 $2FFF
    jmp print

:sumBigDir
    sti $2FFF 0
    ldi $2FFE R1
    tra R1
    addi 1
    tar R1
    sti R1 $2FFE
    jmp print

:final
    htl
