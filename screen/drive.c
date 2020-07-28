#include "drive.h"

FloppyDrive drive;

void init_drive() {
	drive.file = NULL;
	drive.is_end = false;
	drive.readed = 0;
	drive.length = 0;
}

bool open_drive(const char* rom_name) {
	drive.file = fopen(rom_name, "r");
	if(drive.file == NULL) {
		printf("Error reading file %s\n", rom_name);
		return false;
	}
	fseek(drive.file, 0, SEEK_END);
	size_t file_size = ftell(drive.file);
	drive.length = file_size;
	rewind(drive.file);
	return true;
}

bool is_drive_end() {
	if(drive.readed >= drive.length) {
		return true;
	}
	return false;
}

uint8_t read_drive() {
	if(is_drive_end()) {
		return 0x00;
	}
	uint8_t b;
	fread(&b, sizeof(uint8_t), 1, drive.file);
	drive.readed++;
	return b;
}

void close_drive() {
	fclose(drive.file);
	init_drive();
}