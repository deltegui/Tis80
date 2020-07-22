#ifndef tiscpu_rom_h
#define tiscpu_rom_h

#include <stdbool.h>
#include <stdlib.h>

// Rom reader is an interface-like struct that lets TisCpu to
// handle external ROM files.
typedef struct {
	// Tells if bytes cannot be readed
	bool (*is_at_end)();

	// Reads a byte. It souce cant be readed should return
	// any byte. If that happends, is_at_end must return true.
	// It must never stop execution.
	uint8_t (*read)();

	// Opens a new file, and restart read pointer. If something wrong happens
	// it return false. Returns true otherwise.
	bool (*open)(const char* name);

	// Close data stream.
	void (*close)();
} RomReader;

#endif