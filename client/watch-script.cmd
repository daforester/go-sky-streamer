@echo off

for /f "delims=" %%i in ('npm bin') do set NPM_BIN=%%i

SET NODE_ENV=dev

"%NPM_BIN%/webpack" --watch