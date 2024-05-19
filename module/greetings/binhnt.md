# build package 
go build -buildmode=archive .


# View package 
go tool objdump greetings.a

go tool nm  -type T greetings.a