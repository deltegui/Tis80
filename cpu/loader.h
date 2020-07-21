#ifndef tiscpu_loader_h
#define tiscpu_loader_h

#include "rom.h"

void init_loader(RomReader reader);

bool load_rom(const char* rom_name);

#endif