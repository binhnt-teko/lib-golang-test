# Go Execution modes 
https://docs.google.com/document/d/1nr-TQHw_er6GOQRsF6T43GGhFDelrAP0NqSS_00RgZQ/edit?pli=1#heading=h.fwmrrio0df0i

# Compile to shared library 
- In order to create library, we first need to compile & install runtime shared libraries:
  go install -buildmode=shared runtime sync/atomic
  go install -buildmode=shared runtime std
  This will install the create "linux_amd64_dynlink" under $GOROOT/pkg. Here is what directory looks like

ls -R /usr/local/go/pkg/linux_amd64_dynlink/runtime

- Now lets build our dummy library, we need to pass -linkshared flag in order to use the .so file created above

go build -linkshared -buildmode=shared dummy

-  Now lets install this library

go install -linkshared -buildmode=shared dummy


ls -al /go/pkg/linux_amd64_dynlink

ldd libdummy.so 

#  Now lets use this library
go build -linkshared main.go

ldd main
./main