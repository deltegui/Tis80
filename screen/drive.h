#ifndef tisscreen_drive_h
#define tisscreen_drive_h

#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>

typedef struct {
	FILE* file;
	bool is_end;
	size_t readed;
	size_t length;
} FloppyDrive;

void init_drive();
bool open_drive(const char* rom_name);
uint8_t read_drive();
bool is_drive_end();
void close_drive();

#endif