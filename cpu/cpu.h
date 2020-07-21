#ifndef tiscpu_cpu_h
#define tiscpu_cpu_h

#include <stdbool.h>

typedef enum {
	KEYBOARD_INT,
} Interruption;

void init_cpu();

void free_cpu();

bool execute_instruction();

void dispatch_interruption(Interruption interruption);

void write_byte(uint16_t direction, uint8_t data);

uint8_t read_byte(uint16_t direction);

void print_cpu_status();

#endif