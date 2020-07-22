#include <stdio.h>
#include "../cpu/tis.h"
#include "../cpu/error.h"

typedef struct {
	FILE* file;
	bool is_end;
	size_t readed;
	size_t length;
} Rom;

Rom rom;

void init_rom() {
	rom.file = NULL;
	rom.is_end = false;
	rom.readed = 0;
	rom.length = 0;
}

bool open_rom(const char* rom_name) {
	rom.file = fopen(rom_name, "r");
	if(rom.file == NULL) {
		printf("Error reading file %s\n", rom_name);
		return false;
	}
	fseek(rom.file, 0, SEEK_END);
	size_t file_size = ftell(rom.file);
	rom.length = file_size;
	rewind(rom.file);
	return true;
}

bool is_rom_end() {
	if(rom.readed >= rom.length) {
		return true;
	}
	return false;
}

uint8_t read_rom() {
	if(is_rom_end()) {
		return 0x00;
	}
	uint8_t b;
	fread(&b, sizeof(uint8_t), 1, rom.file);
	rom.readed++;
	return b;
}

void close_rom() {
	fclose(rom.file);
}

int main() {
	RomReader reader = {
		.open = &open_rom,
		.is_at_end = &is_rom_end,
		.read = &read_rom,
		.close = &close_rom,
	};
	TisErr err = init_tis(reader);
	if(err != ErrNone) {
		printf("Error while initializing Tis80: %s. Emitted when reading byte %zu\n", tis_error_string(err), rom.readed);
		exit(1);
	}
	print_status();
	err = ErrNone;
	while(err == ErrNone) {
		err = execute_instruction();
	}
	if(err != ErrNone && err != ErrExecEnd) {
		printf("Error while executing assembler: %s. \n", tis_error_string(err));
		free_tis();
		return 1;
	}
	print_status();
	free_tis();
	return 0;
}