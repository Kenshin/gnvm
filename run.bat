::===========================================================
:: GNVM  : Node.js version manager by GO
:: HOST  : https://github.com/kenshin/gnvm
:: Author: Kenshin<kenshin@ksria.com>
::===========================================================

@ECHO off

IF "%1" == "icon" GOTO icon
IF "%1" == "install" GOTO install
IF "%1" == "test" GOTO test

:icon
@ECHO run rsrc.exe build syso
rsrc -ico gnvm.ico -o gnvm.syso
IF "%1" == "icon" GOTO exit

:install
@ECHO run go install
go install -ldflags "-w -s"
GOTO exit
IF "%1" == "install" GOTO exit

:test
@ECHO go test
go test
GOTO exit

:exit
@ECHO create complete.