# Changes to the language
 - two changes to “for” loops
   - + the variables declared by a “for” loop were created once and updated by each iteration
   - “For” loops may now range over integers
  
# Tools 
- Go command
+ go get is no longer supported outside of a module in the legacy GOPATH mode 
+ go mod init no longer attempts to import module requirements from configuration files for other vendoring tools
+ go test -cover now prints coverage summaries for covered packages that do not have their own test files.
  - Trace: refreshed
  - Vet: 
  - + New warnings for missing values after append
  - + The vet tool now reports a non-deferred call to time.Since(t) within a defer statement
# Runtime
  - The runtime now keeps type-based garbage collection metadata nearer to each heap object

# Compiler

- Profile-guided Optimization (PGO) builds can now devirtualize a higher proportion of calls than previously possible
# Linker
The linker’s -s and -w flags are now behave more consistently across all platforms

# Bootstrap

# Core library
- New math/rand/v2 package

+ The Mitchell & Reeds LFSR generator provided by math/rand’s Source has been replaced by two more modern pseudo-random generator sources
- Enhanced routing patterns
+ A pattern that ends in “/” matches all paths that have it as a prefix, as always. To match the exact pattern including the trailing slash, end it with {$}, as in /exact/match/{$}.


- Minor changes to the library
+ archive/tar
+ archive/zip
+ bufio
+ cmp
+ crypto/tls
+ crypto/x509
+ database/sql
+ debug/elf
+ encoding
+ encoding/json
+ go/ast
+ go/types
+ html/template
+ io
+ log/slog
+ math/big
+ net
+ net/http
+ net/http/cgi
+ net/netip
+ os
+ os/exec
+ reflect
+ runtime/metrics
+ runtime/pprof
+ runtime/trace
+ slices
+ syscall
+ testing/slogtest
