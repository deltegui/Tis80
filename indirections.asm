.data
$0000 0x01
$0001 0x02
$0102 22

.code $0200
	inr $0000 R0
	movi 0xff R15
	inw R15 $0000
	hlt
