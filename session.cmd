@echo off

::===========================================================
:: Initialize
::===========================================================
if not defined NODE_HOME (
    set "NODE_HOME=%~dp0"
)

:: Add %GNVM_SESSION_HOME% to path
set GNVM_SESSION_HOME="%NODE_HOME%"
set path=%GNVM_SESSION_HOME%;%path%

::===========================================================
:: Logic
::===========================================================
if "%1" == ""        goto help
if "%1" == "help"    goto help
if "%1" == "run"     goto run
if "%1" == "exit"    goto exit
if "%1" == "version" goto version
if "%1" == "reg"     goto reg

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
echo   reg                   Set env GNVM_SESSION_HOME to regedit.
echo   exit                  Quit session node.exe version.
echo   version               Show version.
echo;
echo Example:
echo   session help          Show session cli command help.
echo   session run 0.10.24   Set 0.10.24 is session node.exe verison.
echo   session reg           Set GNVM_SESSION_HOME to environment variable.
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
:: regedit : Add GNVM_SESSION_HOME to regedit
::===========================================================
:reg

@echo off
:: runas /user:<admin use> xxx, e.g. runas /user:Kenshin cmd
:: HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Control\Session Manager\Environment
:: HKEY_CURRENT_USER\Environment

:: Get origin path
for /f "tokens=1,2,3,4,*" %%i in ('reg query "HKEY_CURRENT_USER\Environment" ^| find /i "path"') do SET "OldPath=%%k"

:: Add GNVM_SESSION_HOME to regedit
set RegPath=HKEY_CURRENT_USER\Environment
reg add "%RegPath%" /v GNVM_SESSION_HOME /t REG_SZ /d "%NODE_HOME%"
reg add "%RegPath%" /v path /t REG_SZ /d "%NODE_HOME%";%OldPath%

@echo on
@echo Regedit GNVM_SESSION_HOME=%GNVM_SESSION_HOME%.
@echo Regedit successful.
@echo off
goto quit

::===========================================================
:: quit : Quit batch script.
::===========================================================
:quit
exit /b 0
