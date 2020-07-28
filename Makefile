all: assembler diassembler console tis

assembler: folder
	cd ./asm && go build -o ../build/tisasm ./cmd/assembler/main.go && cd ..

diassembler: folder
	cd ./asm && go build -o ../build/tisdiasm ./cmd/diassembler/main.go && cd ..

folder:
	mkdir build

clean:
	rm -rf build

console: folder
	gcc ./console/*.c ./cpu/*.c -o ./build/tisconsole

tis: folder
	gcc ./screen/*.c ./cpu/*.c -o ./build/tis -L/usr/local/lib -lSDL2 -lSDL2_ttf
