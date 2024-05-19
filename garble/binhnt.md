export PATH=$PATH:~/go/bin

# Using go
go build -o plugin.so -buildmode=plugin ./plugin
go build -o main main.go
./main


# Using garble in build main 
go build -buildmode=c-shared -buildmode=plugin ./plugin





LD_PRELOAD=./plugin.so garble build -trimpath=true -o main main.go


LD_PRELOAD=./plugin.so ./main
fatal error: runtime: no plugin module data



# Using garble build plugin 
garble build  -trimpath -o plugin.so  -buildmode=plugin ./plugin

err: 
panic: plugin.Open("plugin"): plugin was built with a different version of package sync/atomic

# Using -buildmode=c-shared and LD_PRELOAD=./cbridge.so