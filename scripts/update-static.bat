@echo off
REM Script to update static JavaScript libraries
REM Downloads latest versions from CDN and updates local files

set STATIC_DIR=static
if not exist "%STATIC_DIR%" mkdir "%STATIC_DIR%"

echo Updating static JavaScript libraries...

REM Update HTMx
echo Checking %STATIC_DIR%\htmx.min.js...
set "temp_file=%STATIC_DIR%\htmx.min.js.tmp"
powershell -Command "Invoke-WebRequest -Uri 'https://cdn.jsdelivr.net/npm/htmx.org@latest/dist/htmx.min.js' -OutFile '%temp_file%'"
if exist "%temp_file%" (
    if exist "%STATIC_DIR%\htmx.min.js" (
        fc /b "%STATIC_DIR%\htmx.min.js" "%temp_file%" >nul 2>&1
        if errorlevel 1 (
            echo   Updating %STATIC_DIR%\htmx.min.js
            move /y "%temp_file%" "%STATIC_DIR%\htmx.min.js" >nul
        ) else (
            echo   %STATIC_DIR%\htmx.min.js is up to date
            del "%temp_file%"
        )
    ) else (
        echo   Creating %STATIC_DIR%\htmx.min.js
        move /y "%temp_file%" "%STATIC_DIR%\htmx.min.js" >nul
    )
) else (
    echo   Failed to download htmx.min.js
)

REM Update Alpine.js
echo Checking %STATIC_DIR%\alpine.js...
set "temp_file=%STATIC_DIR%\alpine.js.tmp"
powershell -Command "Invoke-WebRequest -Uri 'https://cdn.jsdelivr.net/npm/alpinejs@latest/dist/alpine.min.js' -OutFile '%temp_file%'"
if exist "%temp_file%" (
    if exist "%STATIC_DIR%\alpine.js" (
        fc /b "%STATIC_DIR%\alpine.js" "%temp_file%" >nul 2>&1
        if errorlevel 1 (
            echo   Updating %STATIC_DIR%\alpine.js
            move /y "%temp_file%" "%STATIC_DIR%\alpine.js" >nul
        ) else (
            echo   %STATIC_DIR%\alpine.js is up to date
            del "%temp_file%"
        )
    ) else (
        echo   Creating %STATIC_DIR%\alpine.js
        move /y "%temp_file%" "%STATIC_DIR%\alpine.js" >nul
    )
) else (
    echo   Failed to download alpine.js
)

REM Update xterm.js
echo Checking %STATIC_DIR%\xterm.js...
set "temp_file=%STATIC_DIR%\xterm.js.tmp"
powershell -Command "Invoke-WebRequest -Uri 'https://cdn.jsdelivr.net/npm/xterm@latest/lib/xterm.js' -OutFile '%temp_file%'"
if exist "%temp_file%" (
    if exist "%STATIC_DIR%\xterm.js" (
        fc /b "%STATIC_DIR%\xterm.js" "%temp_file%" >nul 2>&1
        if errorlevel 1 (
            echo   Updating %STATIC_DIR%\xterm.js
            move /y "%temp_file%" "%STATIC_DIR%\xterm.js" >nul
        ) else (
            echo   %STATIC_DIR%\xterm.js is up to date
            del "%temp_file%"
        )
    ) else (
        echo   Creating %STATIC_DIR%\xterm.js
        move /y "%temp_file%" "%STATIC_DIR%\xterm.js" >nul
    )
) else (
    echo   Failed to download xterm.js
)

REM Update xterm.css
echo Checking %STATIC_DIR%\xterm.css...
set "temp_file=%STATIC_DIR%\xterm.css.tmp"
powershell -Command "Invoke-WebRequest -Uri 'https://cdn.jsdelivr.net/npm/xterm@latest/css/xterm.css' -OutFile '%temp_file%'"
if exist "%temp_file%" (
    if exist "%STATIC_DIR%\xterm.css" (
        fc /b "%STATIC_DIR%\xterm.css" "%temp_file%" >nul 2>&1
        if errorlevel 1 (
            echo   Updating %STATIC_DIR%\xterm.css
            move /y "%temp_file%" "%STATIC_DIR%\xterm.css" >nul
        ) else (
            echo   %STATIC_DIR%\xterm.css is up to date
            del "%temp_file%"
        )
    ) else (
        echo   Creating %STATIC_DIR%\xterm.css
        move /y "%temp_file%" "%STATIC_DIR%\xterm.css" >nul
    )
) else (
    echo   Failed to download xterm.css
)

echo Done!
