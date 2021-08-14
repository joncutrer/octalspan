@echo off

go-winres make
echo compiling build/octalspan.syslogd.exe
go build -o build/octalspan.syslogd.exe

echo cleanup syso files
del *.syso /F
echo copy svc files
cp --f svc/service.xml build/octalspan.svc.xml
cp --f svc/WinSW.exe build/octalspan.svc.exe
echo copy other files
echo build/octalspan.yml 
cp --f octalspan-default.yml build/octalspan.yml 
echo build/LICENSE
cp --f LICENSE build/LICENSE 
echo build complete
