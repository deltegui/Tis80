#ifndef tiscpu_cpu_h
#define tiscpu_cpu_h

#include <stdbool.h>
#include "error.h"

#define REGISTER_MIN 0
#define REGISTER_MAX 15

#define MEMORY_LENGTH 65536

#define INIT_INT				0x0000
#define INIT_PARAMS				0x0100
#define INIT_STACK 				0x0104
#define INIT_KERNAL_ROM 		0x0200
#define INIT_VID_MEM			0x3000
#define INIT_KEYBOARD_BUFFER	0x4000
#define INIT_RAM 				0x4100

#define FLAG_COUNT 3

#define FLAG_ACC_OVERFLOW 	0
#define FLAG_STACK_OVERFLOW 1
#define FLAG_IO_ERROR 		2

typedef enum {
	ACC_OVERFLOW_INT,
	STACK_OVERFLOW_INT,
	IO_ERROR_INT,
	KEYBOARD_INT,
} Interruption;

typedef struct {
	uint8_t* memory;
	uint8_t registers[REGISTER_MAX + 1];
	uint8_t acc;
	uint16_t pc;
	uint16_t stack_top;
	bool flags[FLAG_COUNT];
	bool halt;
	bool protected_mode;
	bool enabled_interruptions;
} CpuStatus;

void init_cpu();

void free_cpu();

TisErr cpu_execute_instruction();

void dispatch_interruption(Interruption interruption);

void write_byte(uint16_t direction, uint8_t data);

uint8_t read_byte(uint16_t direction);

CpuStatus* get_cpu_status();

void free_cpu_status(CpuStatus* status);

#endif