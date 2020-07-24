#ifndef tiscpu_h
#define tiscpu_h

#include <stdbool.h>
#include <stdlib.h>
#include "rom.h"
#include "error.h"

typedef struct {
	uint8_t* bytes;
	int size;
} ScreenData;

TisErr init_tis(RomReader reader);
void free_tis();

void print_status();

TisErr execute_instruction();

/*
void dispatch_key_down(int key);

ScreenData* get_screen_data();

uint8_t* get_memory_snapshot();
*/

#endif