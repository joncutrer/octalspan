@echo off

go-winres make
echo building bin/octalspan.syslogd.exe
go build -o bin/octalspan.syslogd.exe
echo copy config to bin/octalspan.yml 
cp --f octalspan-default.yml bin/octalspan.yml 
cp --f LICENSE bin/LICENSE 
echo build complete