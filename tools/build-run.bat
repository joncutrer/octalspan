@echo off

call tools\build.bat
cd build
octalspan.syslogd.exe
