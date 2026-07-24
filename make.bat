@echo off

if "%1" == "" (
	echo USAGE: make ^<build/run^>
	exit /b
)

if "%1" == "build" (
	pushd build
	go build -gcflags="all=-N -l" -o gogle.exe ..\code
	popd
	exit /b
)

if "%1" == "test" (
	pushd code\index
		go test -v
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
	del gogle.exe
	popd
	exit /b
)

if "%1" == "tags" (
	echo building tags
	ctags -R .
	EXIT /B
)

if "%1" == "push" (
	git push https://github.com/collins994/gogle main
	EXIT /B
)

echo [ERROR]: unknown option: %1
