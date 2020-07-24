#include <stdlib.h>
#include <stdio.h>
#include "cpu.h"

#define UINT8_COUNT (UINT8_MAX + 1) 

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

#define OP_ADD	0x00
#define OP_ADDI 0x01
#define OP_SUB	0x02
#define OP_SUBI 0x03
#define OP_SIL	0x04
#define OP_SIR	0x05
#define OP_AND	0x06
#define OP_OR	0x07
#define OP_XOR	0x08
#define OP_NOT	0x09
#define OP_JMP	0x20
#define OP_JEQ	0x21
#define OP_JNE	0x22
#define OP_JGT	0x23
#define OP_JLT	0x24
#define OP_JFG	0x25
#define OP_LDR	0x30
#define OP_STR	0x31
#define OP_MOV	0x32
#define OP_MOVI	0x33
#define OP_TAR	0x34
#define OP_TRA	0x35
#define OP_INR	0x36
#define OP_INW	0x37
#define OP_DSK	0x38
#define OP_INT	0x40
#define OP_HLT	0x41
#define OP_CLL	0x42
#define OP_CRN 	0x43
#define OP_PMD 	0x44
#define OP_EIN 	0x45
#define OP_DIN	0x46
#define OP_CFG	0x47
#define OP_PSA	0x50
#define OP_POA	0x51
#define OP_PSR	0x52
#define OP_POR	0x53

typedef struct {
	uint8_t* memory;
	uint8_t registers[REGISTER_MAX + 1];
	uint8_t acc;
	uint8_t* pc;
	uint8_t* stack_top;
	bool halt;
	bool overflow_flag;
	bool protected_mode;
	bool enabled_interruptions;
	bool stack_overflow;
} Cpu;

Cpu cpu;

void init_cpu() {
	cpu.memory = malloc(sizeof(uint8_t)*MEMORY_LENGTH);
	cpu.pc = cpu.memory + INIT_KERNAL_ROM;
	cpu.stack_top = cpu.memory + INIT_STACK;
	cpu.halt = false;
	cpu.overflow_flag = false;
	cpu.stack_overflow = false;
	cpu.protected_mode = false;
	cpu.enabled_interruptions = true;
}

void free_cpu() {
	free(cpu.memory);
}

void dispatch_interruption(Interruption interruption) {
	if(!cpu.enabled_interruptions) {
		return;
	}
}

void write_byte(uint16_t direction, uint8_t data) {
	cpu.memory[direction] = data;
}

uint8_t read_byte(uint16_t direction) {
	return cpu.memory[direction];
}

static void stack_push(uint8_t value) {
	uint16_t next_pos = (uint16_t)cpu.stack_top + 1;
	if(next_pos >= INIT_KERNAL_ROM) {
		cpu.stack_overflow = true;
		return;
	}
	*cpu.stack_top = value;
	cpu.stack_top++;
}

static uint8_t stack_pop() {
	uint8_t value = *cpu.stack_top;
	uint16_t current_pos = (uint16_t)cpu.stack_top;
	if(current_pos > INIT_STACK) {
		cpu.stack_top--;
	}
	return value;

}

static void set_register(uint8_t r, uint8_t value) {
	printf("Set register %d, value %x\n", r, value);
	if(r >= REGISTER_MIN && r <= REGISTER_MAX) {
		printf("SETTED!\n");
		cpu.registers[r] = value;
	}
}

static uint8_t get_register(uint8_t r) {
	if(r < REGISTER_MIN || r > REGISTER_MAX) {
		return 0x00;
	}
	return cpu.registers[r];
}

static bool is_pc_out_bounds() {
	return cpu.pc < cpu.memory || cpu.pc > (cpu.memory + MEMORY_LENGTH);
}

static uint8_t read_pc() {
	if(is_pc_out_bounds()) {
		return 0x00;
	}
	printf("Reading direction %04x, value %02x\n", cpu.pc - cpu.memory, *cpu.pc);
	return *cpu.pc++;
}

static uint16_t to_direction(uint8_t byte_high, uint8_t byte_low) {
	uint16_t high = (uint16_t)byte_high;
	uint16_t low = (uint16_t)byte_low;
	high = (high << 8) & 0xff00;
	low = low & 0x00ff;
	return high + low;
}

static uint16_t read_memory() {
	return to_direction(read_pc(), read_pc());
}

static uint8_t read_indirection(uint16_t direction) {
	uint8_t high = cpu.memory[direction];
	uint8_t low = cpu.memory[direction + 1];
	uint16_t indirection = to_direction(high, low);
	return cpu.memory[indirection];
}

static void write_indirection(uint16_t direction, uint8_t data) {
	uint8_t high = cpu.memory[direction];
	uint8_t low = cpu.memory[direction + 1];
	uint16_t indirection = to_direction(high, low);
	cpu.memory[indirection] = data;
}

static void jump_to(uint16_t direction) {
	cpu.pc = cpu.memory + direction;
}

static void alu_add(uint8_t number) {
	int result = cpu.acc + number;
	printf("SUMA %x\n", result);
	if(result >= UINT8_COUNT) {
		cpu.overflow_flag = true;
	}
	cpu.acc = (uint8_t)result;
}

static void alu_sub(uint8_t number) {
	int result = cpu.acc - number;
	if(result < 0) {
		cpu.overflow_flag = true;
	}
	cpu.acc = (uint8_t)result;
}

static void alu_sift_left() {
	cpu.acc = cpu.acc << 1;
}

static void alu_sift_right() {
	cpu.acc = cpu.acc >> 1;
}

static void alu_and(uint8_t number) {
	cpu.acc = cpu.acc & number;
}

static void alu_or(uint8_t number) {
	cpu.acc = cpu.acc | number;
}

static void alu_not() {
	cpu.acc = ~cpu.acc;
}

static void alu_xor(uint8_t number) {
	cpu.acc = cpu.acc ^ number;
}

static bool is_flag_setted_from_code(uint8_t code) {
	switch(code) {
	case 0:
		return cpu.overflow_flag;
	case 1:
		return cpu.stack_overflow;
	default:
		return false;
	}
}

static void clear_flag_from_code(uint8_t code) {
#define DISABLE_FLAG(flag) cpu.flag = false; break
	switch(code) {
	case 0: {
		DISABLE_FLAG(overflow_flag);
	}	
	case 1: {
		DISABLE_FLAG(stack_overflow);
	}
	}
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
	printf("Protected Mode: %d\n", cpu.protected_mode);
	printf("Enabled Interruptions: %d\n", cpu.enabled_interruptions);
	printf("Overflow: %d\n", cpu.overflow_flag);
	printf("Stack Overflow: %d\n", cpu.stack_overflow);
	printf("\n");
	printf("\n");
	for(int i = 0; i < MEMORY_LENGTH; i++) {
		if(i % 16 == 0) {
			printf("\n $%04x:", i);
		}
		printf(" %02x", cpu.memory[i]);
	}
}

inline TisErr cpu_execute_instruction() {
#define READ_REG() get_register(read_pc())
#define JUMP_IF(OP) do {\
	uint16_t direction = read_memory();\
	if(cpu.acc OP 0) {\
		jump_to(direction);\
	}\
} while(false)

	if(cpu.halt) {
		return ErrExecEnd;
	}
	if(is_pc_out_bounds()) {
		return ErrMemOutBounds;
	}
	switch(read_pc()) {
	case OP_ADD: {
		alu_add(READ_REG());
		printf("[ADD]\n");
		break;
	}
	case OP_ADDI: {
		alu_add(read_pc());
		printf("[ADDI]\n");
		break;
	}
	case OP_SUB: {
		alu_sub(READ_REG());
		printf("[SUB] ACC\n");
		break;
	}
	case OP_SUBI: {
		alu_sub(read_pc());
		printf("[SUBI]\n");
		break;
	}
	case OP_SIL: {
		alu_sift_left();
		printf("[SIL]\n");
		break;
	}
	case OP_SIR: {
		alu_sift_right();
		printf("[SIR]\n");
		break;
	}
	case OP_AND: {
		alu_and(READ_REG());
		printf("[AND]\n");
		break;
	}
	case OP_OR: {
		alu_or(READ_REG());
		printf("[OR]\n");
		break;
	}
	case OP_NOT: {
		alu_not();
		printf("[NOT]\n");
		break;
	}
	case OP_XOR: {
		alu_xor(READ_REG());
		printf("[XOR]\n");
		break;
	}
	case OP_JMP: {
		jump_to(read_memory());
		printf("[JMP]\n");
		break;
	}
	case OP_JEQ: {
		JUMP_IF(==);
		printf("[JEQ]\n");
		break;
	}
	case OP_JNE: {
		JUMP_IF(!=);
		printf("[JNE]\n");
		break;
	}
	case OP_JGT: {
		JUMP_IF(>);
		printf("[JGT]\n");
		break;
	}
	case OP_JLT: {
		JUMP_IF(<);
		printf("[JLT]\n");
		break;
	}
	case OP_JFG: {
		uint8_t flag_id = read_pc();
		uint16_t destiny = read_memory();
		if(is_flag_setted_from_code(flag_id)) {
			jump_to(destiny);
		}
		printf("[JFG]\n");
		break;
	}
	case OP_LDR: {
		uint16_t direction = read_memory();
		uint8_t r = read_pc();
		set_register(r, cpu.memory[direction]);
		printf("[LDR]\n");
		break;
	}
	case OP_STR: {
		uint8_t r = read_pc();
		uint16_t direction = read_memory();
		cpu.memory[direction] = get_register(r);
		printf("[SIR]\n");
		break;
	}
	case OP_MOV: {
		uint8_t from = read_pc();
		uint8_t to = read_pc();
		set_register(to, get_register(from));
		printf("[MOV]\n");
		break;
	}
	case OP_MOVI: {
		uint8_t number = read_pc();
		uint8_t r = read_pc();
		set_register(r, number);
		printf("[MOVI]\n");
		break;
	}
	case OP_TAR: {
		uint8_t r = read_pc();
		set_register(r, cpu.acc);
		printf("[ADD]\n");
		break;
	}
	case OP_TRA: {
		cpu.acc = READ_REG();
		printf("[TRA]\n");
		break;
	}
	case OP_INR: {
		uint16_t direction = read_memory();
		uint8_t r = read_pc();
		set_register(r, read_indirection(direction));
		printf("[INR]\n");
		break;
	}
	case OP_INW: {
		uint8_t r = read_pc();
		uint16_t direction = read_memory();
		write_indirection(direction, get_register(r));
		printf("[INW]\n");
		break;
	}
	case OP_INT: {
		uint8_t interruption = read_pc();
		dispatch_interruption(interruption);
		printf("[INT]\n");
		break;
	}
	case OP_HLT:
		printf("[HLT]\n");
		return ErrExecEnd;
	case OP_PMD: {
		cpu.protected_mode = true;
		printf("[PMD]\n");
		break;
	}
	case OP_EIN: {
		cpu.enabled_interruptions = true;
		printf("[EIN]\n");
		break;
	}
	case OP_DIN: {
		cpu.enabled_interruptions = false;
		printf("[DIN]\n");
		break;
	}
	case OP_CFG: {
		uint8_t flag_id = read_pc();
		clear_flag_from_code(flag_id);
		printf("[CFG]\n");
		break;
	}
	case OP_PSA: {
		stack_push(cpu.acc);
		printf("[PSA]\n");
		break;
	}
	case OP_POA: {
		cpu.acc = stack_pop();
		printf("[POA]\n");
		break;
	}
	case OP_PSR: {
		uint8_t r = read_pc();
		stack_push(get_register(r));
		printf("[PSR]\n");
		break;
	}
	case OP_POR: {
		uint8_t r = read_pc();
		set_register(r, stack_pop());
		printf("[POR]\n");
		break;
	}
	default:
		return ErrExecInstruction;
	}
	return ErrNone;

#undef JUMP_IF
#undef READ_REG
}