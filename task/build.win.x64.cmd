@echo off
cd /d %~dp0..

rem puts version to file for docker image building
git describe --tags > semver
set /p buildvers=<semver
set builddate="%date%"

xcopy .\config %GOPATH%\bin\config /f /d /i /e /k /y
go env -w GOOS=windows GOARCH=amd64
go build -v -o %GOPATH%/bin/omdbfind.x64.exe -ldflags="-X 'main.buildvers=%buildvers%' -X 'main.builddate=%builddate%'" .
