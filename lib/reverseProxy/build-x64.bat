gcc -v
go version
go env
set GOARCH=amd64
set CGO_ENABLED=1
go build -buildmode=c-archive -o rp.a 
gcc rp.def rp.a -shared -lwinmm -lWs2_32 -o rp.dll -Wl,--out-implib,rp.dll.a
lib /def:rp.def /name:rp.dll /out:rp.lib /MACHINE:X64
mkdir x64
copy /y *.a x64\
del /q *.a
copy /y *.dll x64\
del /q *.dll
copy /y *.lib x64\
del /q *.lib
copy /y *.h x64\
del /q *.h
copy /y *.exp x64\
del /q *.exp
