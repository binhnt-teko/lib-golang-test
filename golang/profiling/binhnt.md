# docs:
   - https://blog.stackademic.com/profiling-go-applications-in-the-right-way-with-examples-e784526e9481
   - https://100go.co/98-profiling-execution-tracing/
   - https://github.com/DataDog/go-profiler-notes/blob/main/block.md#relationship-with-mutex-profiling
   - 
# Profiling 
   +  Goroutine: stack traces of all current Goroutines
   +  CPU: stack traces of CPU returned by runtime
   +  Heap: a sampling of memory allocations of live objects
   +  Allocation: a sampling of all past memory allocations
   +  Thread: stack traces that led to the creation of new OS threads
   +  Block: stack traces that led to blocking on synchronization primitives
   +  Mutex: stack traces of holders of contended mutexes
# Start 
cd all 
go test -cpuprofile cpu.prof -memprofile mem.prof -bench .

go run main.go -cpuprofile cpu.prof -memprofile mem.prof 

go tool pprof cpu.prof
# Cpu profiling 
cd cpu 
go run main.go 
go tool pprof profile.prof

go tool pprof mem.prof 


go tool trace trace.out


# Mutex profiling 
go run main.go

go tool pprof  http://localhost:8889/debug/pprof/mutex?debug=1

