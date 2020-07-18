.data
$5000 "Hello world"

.code $4000
add R1
movi R1 10
# tra R1
:loop
# subi 1
jne loop
# int 0
