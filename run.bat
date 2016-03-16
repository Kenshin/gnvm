::===========================================================
:: GNVM   : Node.js version manager by GO
:: HOST   : https://github.com/kenshin/gnvm
:: Author : Kenshin<kenshin@ksria.com>
:: Version: 0.0.2
::===========================================================

@echo off

::===========================================================
:: Initialize
::===========================================================
if "%1" == "icon"    goto icon
if "%1" == "test"    goto test
if "%1" == "install" (
    if "%2" == ""    goto install
    if "%2" == "x86" (
        set GOARCH=386
        call :build %2
    )
    if "%2" == "x64" (
        set GOARCH=amd64
        call :build %2
    )
)

::===========================================================
:: icon : Set icon.
::===========================================================
:icon
@echo run rsrc.exe build syso
rsrc -ico gnvm.ico -o gnvm.syso
goto quit

::===========================================================
:: install : Install current os runtime
::===========================================================
:install
@echo run go install
go install -ldflags "-w -s"
goto quit

::===========================================================
:: build : Cross compile, support x86 and x64 on Windows
::===========================================================
:build
@echo run go build %1
go build -ldflags "-w -s"
goto quit

::===========================================================
:: test : Test gnvm
::===========================================================
:test
@echo go test
go test
goto quit

::===========================================================
:: quit : Quit batch script.
::===========================================================
:quit
@echo gnvm.exe compile success.
exit /b 0
