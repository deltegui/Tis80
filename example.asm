.code $0200
	movi 0 R0
	tra R0

:start
	addi 1
	jfg 0 end
	jmp start

:end
	movi 1 R1
	cfg 0
	hlt