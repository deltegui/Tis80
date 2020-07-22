#include <stdlib.h>
#include <stdio.h>
#include "cpu.h"

#define MEMORY_LENGTH 65536
#define INIT_KERNAL_ROM 0x0200

#define OP_ADD	0x00

#define OP_MOVI	0X33
#define OP_TAR 	0X34
#define OP_TRA 	0X35

#define OP_HLT 	0X41

typedef struct {
	uint8_t* memory;
	uint8_t registers[16];
	uint8_t acc;
	uint8_t* pc;
	bool halt;
} Cpu;

Cpu cpu;

void init_cpu() {
	cpu.memory = malloc(sizeof(uint8_t)*MEMORY_LENGTH);
	cpu.pc = cpu.memory + INIT_KERNAL_ROM - 1;
	cpu.halt = false;
}

void free_cpu() {
	free(cpu.memory);
}

void dispatch_interruption(Interruption interruption) {

}

void write_byte(uint16_t direction, uint8_t data) {
	cpu.memory[direction - 1] = data;
}

uint8_t read_byte(uint16_t direction) {
	return cpu.memory[direction - 1];
}

static void set_register(uint8_t r, uint8_t value) {
	if(r < 0 || r > 15) {
		return;
	}
	cpu.registers[r] = value;
}

static uint8_t get_register(uint8_t r) {
	if(r < 0 || r > 15) {
		return 0x00;
	}
	return cpu.registers[r];
}

static bool is_pc_out_bounds() {
	return cpu.pc < cpu.memory || cpu.pc >= (cpu.memory + MEMORY_LENGTH);
}

static uint8_t read_pc() {
	if(is_pc_out_bounds()) {
		return 0x00;
	}
	printf("Reading direction %04x, value %02x\n", cpu.pc - cpu.memory + 1, *cpu.pc);
	return *cpu.pc++;
}

void print_cpu_status() {
	printf("------TIS 80 CPU STATUS-----\n");
	printf("\n");
	printf("ACC register: %02x\n", cpu.acc);
	printf("\n");
	for(int i = 0; i < 16; i++) {
		printf("R%d: %02x\n", i, cpu.registers[i]);
	}
	printf("\n");
	printf("\n");
	for(int i = 0; i < MEMORY_LENGTH; i++) {
		if(i % 16 == 0) {
			printf("\n $%04x:", i);
		}
		printf(" %02x", cpu.memory[i]);
	}
}

TisErr cpu_execute_instruction() {
	if(cpu.halt) {
		return ErrExecEnd;
	}
	if(is_pc_out_bounds()) {
		return ErrMemOutBounds;
	}
	switch(read_pc()) {
	case OP_ADD: {
		uint8_t r = read_pc();
		cpu.acc = cpu.acc + get_register(r);
		printf("[ADD] ACC + R%d\n", r);
		break;
	}
	case OP_MOVI: {
		uint8_t number = read_pc();
		uint8_t r = read_pc();
		cpu.registers[r] = number;
		printf("[MOVI] %d -> R%d\n", number, r);
		break;
	}
	case OP_TRA: {
		uint8_t r = read_pc();
		cpu.acc = get_register(r);
		printf("[TRA] R%d -> ACC\n", r);
		break;
	}
	case OP_TAR: {
		uint8_t r = read_pc();
		set_register(r, cpu.acc);
		printf("[ADD] ACC -> R%d\n", r);
		break;
	}
	case OP_HLT:
		printf("[HLT]\n");
		return ErrExecEnd;
	default:
		return ErrExecInstruction;
	}
	return ErrNone;
}
