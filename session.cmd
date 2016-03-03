@echo off

::===========================================================
:: Initialize
::===========================================================
if not defined NODE_HOME (
    set "NODE_HOME=%~dp0"
)

::===========================================================
:: Logic
::===========================================================
if "%1" == ""        goto help
if "%1" == "help"    goto help
if "%1" == "run"     goto run
if "%1" == "exit"    goto exit
if "%1" == "version" goto version

::===========================================================
:: help : Show help message
::===========================================================
:help
echo;
echo GNVM - Session node.exe manage
echo;
echo Usage:
echo   session [command]
echo;
echo Commands:
echo   run                   Set session node.exe version.
echo   exit                  Quit session node.exe version.
echo   version               Show version.
echo;
echo Example:
echo   session help          Show session cli command help.
echo   session run 0.10.24   Set 0.10.24 is session node.exe verison.
echo   session exit          Quit sesion node.exe, restore global node.exe version.
echo   session version       Show version.
goto quit

::===========================================================
:: version : Show session.bat version( same as gnvm.exe version)
::===========================================================
:version
echo Current version 0.2.0.
echo Copyright (C) 2014-2016 Kenshin Wang kenshin@ksria.com
echo See https://github.com/kenshin/gnvm for more information.
goto quit

::===========================================================
:: run : Set session node.exe
::===========================================================
:run
if "%2" == "" (
    @echo on
    @echo Parameter can't be empty.
    @echo Example: "session run 5.7.0"
    @echo off
    goto quit
)

@echo off
set GNVM_SESSION_NODE_HOME=%NODE_HOME%\%2\
set path=%GNVM_SESSION_NODE_HOME%;%path%

@echo on
@echo Startup session node.exe version %2.
@echo Important:
@echo - if node.exe work on session version, 'gnvm install -g', 'gnvm update -g' 'gnvm use x.xx.xx' can't be use.
@echo - if quit/remove session, you must use 'session exit'.
@echo - if on "%NODE_HOME%" directory, unable to 'run %2'.
@echo off
goto quit

::===========================================================
:: exit : Quit/Remvoe session node.exe version
::===========================================================
:exit

@echo off
set GNVM_SESSION_NODE_HOME=
set path=%NODE_HOME%;%path%

@echo on
@echo Exit session successful.
@echo off
goto quit

::===========================================================
:: quit : Quit batch script.
::===========================================================
:quit
exit /b 0
