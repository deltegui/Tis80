#include <stdio.h>
#include <stdbool.h>
#include <SDL2/SDL.h>
#include <SDL2/SDL_ttf.h>
#include "screen.h"
#include "drive.h"
#include "../cpu/tis.h"
#include "../cpu/error.h"

void print_status() {
	CpuStatus* status = get_status();
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
	printf("Stack top: $%04x\n", status->stack_top);
	printf("Program counter: $%04x\n", status->pc);
	printf("\n");
	for(int i = 0; i < MEMORY_LENGTH; i++) {
		if(i % 16 == 0) {
			printf("\n $%04x:", i);
		}
		printf(" %02x", status->memory[i]);
	}
	printf("\n");
	free_status(status);
}

char* get_video_string(CpuStatus* status) {
	uint8_t* init_screen = &status->memory[INIT_VID_MEM];
	status->memory[INIT_KEYBOARD_BUFFER - 1] = '\0';
	return (char*)init_screen;
}

bool begin_tis() {
	RomReader reader = {
		.open = &open_drive,
		.is_at_end = &is_drive_end,
		.read = &read_drive,
		.close = &close_drive,
	};
	TisErr err = init_tis(reader);
	if(err != ErrNone) {
		printf("Error while initializing Tis80: %s\n", tis_error_string(err));
		return false;
	}
	return true;
}

void end_tis() {
	free_tis();
}

void loop() {
	SDL_Event event;
	bool quit = false;
	bool not_first = true;
	TisErr err = ErrNone;
	while(!quit) {
		while(SDL_PollEvent(&event)){ 
            switch( event.type ){
            case SDL_QUIT:
                quit = true;
                break;
            default:
                break;
            }
		}
		if(err == ErrNone) {
			err = execute_instruction();
		}
		if(err != ErrNone && err != ErrExecEnd) {
			print_status();
			printf("Error while executing assembler: %s. \n", tis_error_string(err));
			quit = true;
			break;
		}
		if(err == ErrExecEnd && not_first) {
			CpuStatus* status = get_status();
			print_text_screen(get_video_string(status));
			free_status(status);
			not_first = false;
		}
	}
	print_status();
}

int main() {
	init_drive();
	if(!init_screen()) {
		return 1;
	}
	if(!begin_tis()) {
		free_screen();
		return 1;
	}

	loop();

	end_tis();
	free_screen();
	close_drive();
	return 0;
}