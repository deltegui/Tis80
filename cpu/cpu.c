#include <stdlib.h>
#include <stdio.h>
#include "cpu.h"

#define MEMORY_LENGTH 65536

typedef struct {
	uint8_t* memory;
	uint8_t registers[16];
	uint8_t acc;
	uint8_t* pc;
} Cpu;

Cpu cpu;

void init_cpu() {
	cpu.memory = malloc(sizeof(uint8_t)*MEMORY_LENGTH);
}

void free_cpu() {
	free(cpu.memory);
}

bool execute_instruction() {
	return false;
}

void dispatch_interruption(Interruption interruption) {

}

void write_byte(uint16_t direction, uint8_t data) {
	cpu.memory[direction] = data;
}

uint8_t read_byte(uint16_t direction) {
	return cpu.memory[direction];
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