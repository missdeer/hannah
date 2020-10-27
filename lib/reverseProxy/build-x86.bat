gcc -v
go version
go env
set GOARCH=386
set CGO_ENABLED=1
go build -buildmode=c-archive -o rp.a 
gcc rp.def rp.a -shared -lwinmm -lWs2_32 -o rp.dll -Wl,--out-implib,rp.dll.a
lib /def:rp.def /name:rp.dll /out:rp.lib /MACHINE:X86
del /q *.exp
