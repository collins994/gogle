@echo off

if "%1" == "" (
	echo USAGE: make ^<build/run^>
)

if "%1" == "build" (
	pushd build
	go build -o gogle.exe ..\code
	popd
	exit /b
)


if "%1" == "run" (
	pushd build
	gogle.exe
	popd
	exit /b
)

if "%1" == "clean" (
	pushd build
	del *
	popd
	exit /b
)

if "%1" == "tags" (
	echo building tags
	ctags -R .
	EXIT /B
)

echo [ERROR]: unknown option: %1
