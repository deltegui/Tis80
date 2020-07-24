package tisasm

import (
	"fmt"
	"os"
)

func ShowError(msg string) {
	fmt.Printf("[ERROR] %s", msg)
	fmt.Println()
	os.Exit(1)
}

func ShowErrorf(format string, replaces ...interface{}) {
	ShowError(fmt.Sprintf(format, replaces...))
}

func ShowErrorToken(token Token, msg string) {
	fmt.Printf("[ERROR] %s in %s", msg, token)
	fmt.Println()
	os.Exit(1)
}

func ShowErrorTokenf(token Token, format string, replaces ...interface{}) {
	ShowErrorToken(token, fmt.Sprintf(format, replaces...))
}
