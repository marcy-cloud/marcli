@echo off
REM Script to update static JavaScript libraries
REM Downloads latest versions from CDN and updates local files

set STATIC_DIR=static
if not exist "%STATIC_DIR%" mkdir "%STATIC_DIR%"

echo Updating static JavaScript libraries...

REM Function to download and update a file if it's different
:update_file
set url=%~1
set file=%~2
set temp_file=%file%.tmp

echo Checking %file%...

REM Download to temp file using PowerShell
powershell -Command "Invoke-WebRequest -Uri '%url%' -OutFile '%temp_file%'"

if exist "%temp_file%" (
    if exist "%file%" (
        REM Compare files
        fc /b "%file%" "%temp_file%" >nul 2>&1
        if errorlevel 1 (
            echo   Updating %file%
            move /y "%temp_file%" "%file%" >nul
        ) else (
            echo   %file% is up to date
            del "%temp_file%"
        )
    ) else (
        echo   Creating %file%
        move /y "%temp_file%" "%file%" >nul
    )
) else (
    echo   Failed to download %file%
)

goto :eof

REM Update HTMx
call :update_file "https://cdn.jsdelivr.net/npm/htmx.org@latest/dist/htmx.min.js" "%STATIC_DIR%\htmx.min.js"

REM Update Alpine.js
call :update_file "https://cdn.jsdelivr.net/npm/alpinejs@latest/dist/alpine.min.js" "%STATIC_DIR%\alpine.js"

REM Update xterm.js
call :update_file "https://cdn.jsdelivr.net/npm/xterm@latest/lib/xterm.js" "%STATIC_DIR%\xterm.js"

REM Update xterm.css
call :update_file "https://cdn.jsdelivr.net/npm/xterm@latest/css/xterm.css" "%STATIC_DIR%\xterm.css"

echo Done!

