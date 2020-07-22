#include "error.h"

char* errors[] = {
	"No error",
	"Error while reading from ROM",
	"Error while loading ROM: bad format",
	"Execution reached end",
	"Undefined instruction",
	"Program tried to read a memory out of bounds",
};

char* tis_error_string(TisErr error) {
	return errors[error];
}