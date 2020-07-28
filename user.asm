.data
$5000 "Bienvenido a esta demo del Tis80. Puede hacer por el momento pocas cosas."

.code $4100
	movm $5000 $0100
	movm $3000 $0102
	int 3
	crn