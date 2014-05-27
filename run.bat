::===========================================================
:: GNVM: Node.exe version manage by GO
::===========================================================

@echo off

IF "%1" == "icon" goto icon
IF "%1" == "go" goto go

:icon
@echo rsrc -ico gnvm.ico -o gnvm.syso
rsrc -ico gnvm.ico -o gnvm.syso
goto exit

:go
@echo go install -ldflags "-w -s"
go install -ldflags "-w -s"
goto exit

:exit
@echo create complete.