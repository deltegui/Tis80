#include "tis.h"
#include "cpu.h"
#include "loader.h"

void init_tis(RomReader reader) {
	init_cpu();
	init_loader(reader);
	load_rom("kernal.rom");
}

void free_tis() {
	free_cpu();
}

void print_status() {
	print_cpu_status();
}