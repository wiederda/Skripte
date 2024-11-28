$env:GOOS = "linux"
$env:GOARCH = "amd64"
go build -o ./linux/get-ext-ip main.go

$env:GOOS = "windows"
$env:GOARCH = "amd64"
go build -o ./windows/get-ext-ip.exe main.go

$env:GOOS = "darwin"
$env:GOARCH = "amd64"
go build -o ./macos/amd64/get-ext-ip main.go

$env:GOOS = "darwin"
$env:GOARCH = "arm64"
go build -o ./macos/arm64/get-ext-ip main.go