all: assembler diassembler tis

assembler: folder
	cd ./asm && go build -o ../build/tisasm ./cmd/assembler/main.go && cd ..

diassembler: folder
	cd ./asm && go build -o ../build/tisdiasm ./cmd/diassembler/main.go && cd ..

folder:
	mkdir build

clean:
	rm -rf build

tis: folder
	gcc ./console/*.c ./cpu/*.c -o ./build/tis
