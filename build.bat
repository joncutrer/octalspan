@echo off

go-winres make
echo compiling bin/octalspan.syslogd.exe
go build -o bin/octalspan.syslogd.exe

echo cleanup syso files
del *.syso /F
echo copy svc files
cp --f svc/service.xml bin/octalspan.svc.xml
cp --f svc/WinSW.exe bin/octalspan.svc.exe
echo copy other files
echo bin/octalspan.yml 
cp --f octalspan-default.yml bin/octalspan.yml 
echo bin/LiCENSE
cp --f LICENSE bin/LICENSE 
echo build complete