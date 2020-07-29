.data
$5000 "Bienvenido a esta demo del Tis80. Esta version grafica obtiene la memoria de video y la representa en esta pantalla."

.code $4100
	movm $5000 $0100
	movm $3000 $0102
	int 4
	crn