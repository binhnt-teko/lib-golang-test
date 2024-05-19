 docker run -it -v $(pwd):/app -w /app golang:1.19.2 bash

# Compile plugin 
go build -buildmode=plugin -o eng/eng.so eng/greeter.go
go build -buildmode=plugin -o chi/chi.so chi/greeter.go
go build -buildmode=plugin -o ru/ru.so ru/greeter.go

# run code 
go run greeter.go english
go run greeter.go chinese
