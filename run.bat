::===========================================================
:: GNVM  : Node.exe version manage by GO
:: HOST  : https://github.com/kenshin/gnvm
:: author: Kenshin<kenshin@ksria.com>
::===========================================================

@ECHO off

IF "%1" == "icon" GOTO icon
IF "%1" == "go" GOTO go

:icon
@ECHO run rsrc.exe build syso
rsrc -ico gnvm.ico -o gnvm.syso
IF "%1" == "icon" GOTO exit

:go
@ECHO run go install
go install -ldflags "-w -s"
GOTO exit

:exit
@ECHO create complete.