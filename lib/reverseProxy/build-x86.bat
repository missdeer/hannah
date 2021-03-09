gcc -v
ren librp.a rp.a
windres --input rp.dll.rc --output rp.dll.res --output-format=coff
gcc rp.def rp.a rp.dll.res -shared -lwinmm -lWs2_32 -o rp.dll -s -Wl,--subsystem,windows,--out-implib,rp.dll.a
lib /def:rp.def /name:rp.dll /out:rp.lib /MACHINE:X86
del /q *.exp
