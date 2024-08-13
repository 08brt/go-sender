@echo off
cd /d D:\Developer\GoSender

REM Build the Go project (optional, if you haven't built it yet)
go build -o go-sender.exe

REM Run the executable
go-sender.exe

pause