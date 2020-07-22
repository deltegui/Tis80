#ifndef tiscpu_err_h
#define tiscpu_err_h

typedef enum {
	ErrNone,
	ErrRomRead,
	ErrRomFormat,
	ErrExecEnd,
	ErrExecInstruction,
	ErrMemOutBounds,
} TisErr;

char* tis_error_string(TisErr error);

#endif