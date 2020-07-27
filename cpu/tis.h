#ifndef tiscpu_h
#define tiscpu_h

#include <stdbool.h>
#include <stdlib.h>
#include "rom.h"
#include "error.h"
#include "cpu.h"

TisErr init_tis(RomReader reader);
void free_tis();

TisErr execute_instruction();

//void dispatch_key_down(int key);

CpuStatus* get_status();
void free_status(CpuStatus* status);

#endif