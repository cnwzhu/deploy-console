@echo off
echo build platform %1
set GOOS=%1
go build -tags netgo console.go
