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

TisErr execute_instruction() {
	return cpu_execute_instruction();
}

CpuStatus* get_status() {
	return get_cpu_status();
}

void free_status(CpuStatus* status) {
	free_cpu_status(status);
}