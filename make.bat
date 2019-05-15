@echo off

if exist make.bat goto ok
echo Must run make.bat from Metapod's src directory. 1>&2
goto end
:ok

:: Setup
set METAPOD_BK_GOPATH=%GOPATH%
set METAPOD_BK_GOARCH=%GOARCH%
set GOPATH=%~dp0
set CGO_ENABLED=1

if "%GOOS%" == "" set GOOS=windows

mkdir bin\%GOOS%\x64 >NUL 2>&1
mkdir bin\%GOOS%\x32 >NUL 2>&1

:: 64-bit
echo Building 64-bit
set GOARCH=amd64
pushd bin\%GOOS%\x64
go build -buildmode=c-shared metapod
popd

:: 32-bit
echo Building 32-bit
set GOARCH=386
pushd bin\%GOOS%\x32
go build -buildmode=c-shared metapod
popd

:: Cleanup
set GOPATH=%METAPOD_BK_GOPATH%
set GOARCH=%METAPOD_BK_GOARCH%
set METAPOD_BK_GOPATH=
set METAPOD_BK_GOARCH=

:end
