#include "tis.h"
#include "cpu.h"
#include "loader.h"

TisErr init_tis(RomReader reader) {
	init_cpu();
	init_loader(reader);
	TisErr err = load_rom("kernal.rom");
	if(err != ErrNone) {
		free_cpu();
	}
	return err;
}

void free_tis() {
	free_cpu();
}

void print_status() {
	print_cpu_status();
}

TisErr execute_instruction() {
	return cpu_execute_instruction();
}