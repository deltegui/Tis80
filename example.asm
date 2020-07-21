.data
$4FFF 255
$5000 "Hello world"

.code $5050
:loop
	jne stop
	movi 1 R1
	add R1
	jne loop
	cll nothing

:stop
	hlt

:nothing
	pmd
	din
	ein
	crn