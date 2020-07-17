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
