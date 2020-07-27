.data
$5000 "La vida es dura pero mas dura es mi verdura!"

.code $4100
	movm $5000 $0100
	movm $3000 $0102
	int 3
	crn