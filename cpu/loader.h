#ifndef tiscpu_loader_h
#define tiscpu_loader_h

#include "rom.h"
#include "error.h"

void init_loader(RomReader reader);

TisErr load_rom(const char* rom_name);

#endif