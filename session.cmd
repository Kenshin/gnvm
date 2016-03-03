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
if "%1" == "clear"    goto clear
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
echo   clear                 Quit session node.exe version.
echo   version               Show version.
echo;
echo Example:
echo   session help          Show session cli command help.
echo   session run 0.10.24   Set 0.10.24 is session node.exe verison.
echo   session clear         Quit sesion node.exe, restore global node.exe version.
echo   session version       Show version.
goto exit

::===========================================================
:: version : Show session.bat version
::===========================================================
:version
echo Current version 0.0.1.
echo Copyright (C) 2014-2016 Kenshin Wang kenshin@ksria.com
echo See https://github.com/kenshin/gnvm for more information.
goto exit

::===========================================================
:: run : Set session node.exe
::===========================================================
:run

if "%2" == "" (
    @echo on
    @echo Parameter can't be empty.
    @echo Example: "session run 5.7.0"
    @echo off
    goto exit
)

:: if on the %NODE_HOME% directory, goto gnvm_session directory.
if  %cd% == %NODE_HOME% call :security

@echo off
set GNVM_SESSION_NODE_HOME=%NODE_HOME%\%2\
set path=%GNVM_SESSION_NODE_HOME%;%path%

@echo on
@echo Startup session node.exe version %2.
@echo Important:
@echo - if node.exe work on session version, "gnvm install -g", "gnvm update -g" "gnvm use x.xx.xx" can't be use.
@echo - if quit/remove session, you must use "session clear".
@echo - if on "%NODE_HOME%" directory, unable to "run %2".
@echo - if on "%NODE_HOME%" directory, auto goto "%NODE_HOME%\gnvm_session" directory.
@echo - if on "%NODE_HOME%\gnvm_session" directory, use "session clear" auto previous directory.
@echo off
goto exit

::===========================================================
:: security : Security directory.
::===========================================================
:security

:: Add %NODE_HOME% to path
set path=%NODE_HOME%;%path%

:: Add %GNVM_SESSION_HOME% to path
set GNVM_SESSION_HOME=%NODE_HOME%\gnvm_session
set path=%GNVM_SESSION_HOME%;%path%

:: Save current path
set ORI_GNVM_SESSION_PATH=%~dp0

:: Create and goto gnvm_session directory
md gnvm_session
attrib +h gnvm_session
cd %GNVM_SESSION_HOME%
goto exit

::===========================================================
:: clear : Quit/Remove session node.exe version
::===========================================================
:clear

if defined ORI_GNVM_SESSION_PATH (
    cd %ORI_GNVM_SESSION_PATH%
    rd /q /s gnvm_session
)

@echo off
set ORI_GNVM_SESSION_PATH=
set GNVM_SESSION_NODE_HOME=
set path=%NODE_HOME%;%path%

@echo on
@echo Session clear complete.
@echo off
goto exit

::===========================================================
:: exit : Quit batch script.
::===========================================================
:exit
exit /b 0
