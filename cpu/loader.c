#include "loader.h"
#include "cpu.h"

#define DATA_SECTION 0x00
#define CODE_SECTION 0x01

#define END_DATA_TYPE 0x00
#define NUMBER_TYPE 0x02
#define STRING_TYPE 0x01

#define END_STRING 0x00

void expect_section_header();
void expect_byte(uint8_t byte);
void read_data_section();
void read_code_section();
void read_string(uint16_t direction);
bool read_data_type(uint16_t direction);
uint16_t read_memory();

typedef struct {
	RomReader reader;
	bool success;
} Loader;

Loader loader;

void init_loader(RomReader reader) {
	loader.reader = reader;
	loader.success = true;
}

bool load_rom(const char* rom_name) {
	if(!loader.reader.open(rom_name)) {
		return false;
	}
	expect_section_header();
	switch(loader.reader.read()) {
	case DATA_SECTION:
		read_data_section();
		return loader.success;
	case CODE_SECTION:
		read_code_section();
		return loader.success;
	default:
		return false;
	}
	return loader.success;
}

void read_data_section() {
	data:
	while(!loader.reader.is_at_end()) {
		if(!loader.success) {
			return;
		}
		uint16_t direction = read_memory();
		if(!read_data_type(direction)) {
			break;
		}
	}
	expect_section_header();
	expect_byte(CODE_SECTION);
	if(!loader.success) {
		return;
	}
	read_code_section();
}

bool read_data_type(uint16_t direction) {
	switch(loader.reader.read()) {
	case END_DATA_TYPE:
		return false;
	case NUMBER_TYPE:
		write_byte(direction, loader.reader.read());
		return true;
	case STRING_TYPE:
		read_string(direction);
		return true;
	default:
		return false;
	}
}

void read_string(uint16_t direction) {
	uint8_t current_byte = loader.reader.read();
	uint16_t current_dir = direction;
	while(current_byte != END_STRING && !loader.reader.is_at_end()) {
		write_byte(current_dir, current_byte);
		current_byte = loader.reader.read();
		current_dir++;
	}
}

void read_code_section() {
	uint16_t start_code = read_memory();
	uint16_t offset = 0;
	while(!loader.reader.is_at_end() && loader.success) {
		uint8_t byte = loader.reader.read();
		write_byte(start_code+offset, byte);
		offset++;
	}
}

uint16_t read_memory() {
	uint16_t high = (uint16_t)loader.reader.read();
	uint16_t low = (uint16_t)loader.reader.read();
	high = (high << 8) & 0xff00;
	low = low & 0x00ff;
	return high + low;
}

void expect_section_header() {
	expect_byte(0xff);
	expect_byte(0xfe);
	expect_byte(0xfe);
	expect_byte(0xff);
}

void expect_byte(uint8_t byte) {
	loader.success = loader.reader.read() == byte;
}