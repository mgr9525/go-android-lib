: mkdir bin
: SET dir=%cd%\bin
if not exist "bin" mkdir bin

gomobile bind -o %dir%/golib.aar go-android-lib
