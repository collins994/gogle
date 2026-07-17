@echo off
setlocal enabledelayedexpansion

set n=1

for %%f in (build\gl2\*.xhtml) do (
	ren %%f !n!.xhtml
	set /A n+=1
)
