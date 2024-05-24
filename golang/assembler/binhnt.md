# Go's Assembler
- The assembler is based on the input style of the Plan 9 assemblers, which is documented in detail elsewhere
-  not a direct representation of the underlying machine. Some of the details map precisely to the machine, but some do not.
  
- View asm 
GOOS=linux GOARCH=amd64 go tool compile -S  x.go

- To see what gets put in the binary after linking, use go tool objdump
  
go build x.go 
go tool objdump -s main.main x\

# Constants
# Symbols
- set of pseudo-registers is the same for all architectures
+ FP: Frame pointer: arguments and locals.
+ PC: Program counter: jumps and branches.
+ SB: Static base pointer: global symbols.
+ SP: Stack pointer: the highest address within the local stack frame
