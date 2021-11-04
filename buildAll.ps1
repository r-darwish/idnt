$env:GOOS = "darwin"
go build -o idnt-darwin

$env:GOOS = "linux"
go build -o idnt-linux

$env:GOOS = "windows"
go build -o idnt.exe
