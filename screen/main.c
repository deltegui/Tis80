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
	init_rom();
}

void print_status(CpuStatus* status) {
	printf("------TIS 80 CPU STATUS-----\n");
	printf("\n");
	printf("ACC register: %02x\n", status->acc);
	printf("\n");
	for(int i = 0; i < 16; i++) {
		printf("R%d: %02x\n", i, status->registers[i]);
	}
	printf("\n");
	printf("Protected Mode: %d\n", status->protected_mode);
	printf("Enabled Interruptions: %d\n", status->enabled_interruptions);
	printf("Overflow: %d\n", status->flags[FLAG_ACC_OVERFLOW]);
	printf("Stack Overflow: %d\n", status->flags[FLAG_STACK_OVERFLOW]);
	printf("IO error: %d\n", status->flags[FLAG_IO_ERROR]);
	printf("\n");
	printf("\n");
	for(int i = 0; i < MEMORY_LENGTH; i++) {
		if(i % 16 == 0) {
			printf("\n $%04x:", i);
		}
		printf(" %02x", status->memory[i]);
	}
	printf("\n");
}

void print_screen(CpuStatus* status) {
	uint8_t* init_screen = &status->memory[INIT_VID_MEM];
	status->memory[INIT_KEYBOARD_BUFFER - 1] = '\0';
	printf("------------OUTPUT-----------\n");
	printf("%s\n", init_screen);
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
	err = ErrNone;
	while(err == ErrNone) {
		err = execute_instruction();
	}
	if(err != ErrNone && err != ErrExecEnd) {
		printf("Error while executing assembler: %s. \n", tis_error_string(err));
		free_tis();
		return 1;
	}
	CpuStatus* status = get_status();
	print_status(status);
	print_screen(status);
	free_status(status);
	free_tis();
	return 0;
}