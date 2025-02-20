@echo off
setlocal

set APP_NAME=SilverBullet
set OUTPUT_DIR=output
set EXE_NAME=%APP_NAME%.exe
set INSTALLER_NAME=%APP_NAME%-installer.exe

:: Ensure required tools are installed
where makensis >nul 2>nul
if %ERRORLEVEL% neq 0 (
    echo [Error] NSIS (Nullsoft Scriptable Install System) is not installed.
    exit /b 1
)

:: Create output directory
mkdir %OUTPUT_DIR% 2>nul

:: Copy SPV resolver binary
copy spvproc.exe %OUTPUT_DIR%\%EXE_NAME%

:: Create NSIS script
set NSIS_SCRIPT=%OUTPUT_DIR%\%APP_NAME%.nsi
(echo !define APPNAME "%APP_NAME%"
echo !define OUTPUTFILE "%INSTALLER_NAME%"
echo !define EXEFILE "%EXE_NAME%"
echo OutFile %OUTPUT_DIR%\%INSTALLER_NAME%
echo InstallDir $PROGRAMFILES\%APP_NAME%
echo Page directory
echo Page instfiles
echo Section "Install"
echo SetOutPath $INSTDIR
echo File %OUTPUT_DIR%\%EXE_NAME%
echo SectionEnd) > %NSIS_SCRIPT%

:: Generate installer
makensis %NSIS_SCRIPT%

echo [Success] %INSTALLER_NAME% created successfully.
endlocal