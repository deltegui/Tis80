package tisasm

import (
	"os"
)

func OpenFile(path string) *os.File {
	data, err := os.Open(path)
	if err != nil {
		ShowError("File not found")
	}
	return data
}

func CreateFile(path string) *os.File {
	file, err := os.Create(path)
	if err != nil {
		ShowErrorf("%s", err)
	}
	return file
}
