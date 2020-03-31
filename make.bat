
if exist make.bat goto ok
echo Must run make.bat from Metapod's src directory. 1>&2
goto end
:ok

set CURDIR=%~dp0

:: Setup
if "%GOOS%" == "" set GOOS=windows

mkdir bin\%GOOS%\x64 >NUL 2>&1
mkdir bin\%GOOS%\x32 >NUL 2>&1

pushd %CURDIR%metapod-c

:: 64-bit
echo Building 64-bit
set GOARCH=amd64
go build -o ..\bin\%GOOS%\x64 -buildmode=c-shared

:: 32-bit
echo Building 32-bit
set GOARCH=386
go build -o ..\bin\%GOOS%\x32 -buildmode=c-shared

popd


:: Cleanup
:end
