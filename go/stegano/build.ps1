$env:GOOS = "linux"
$env:GOARCH = "amd64"
go build -o ./linux/stegano main.go

$env:GOOS = "windows"
$env:GOARCH = "amd64"
go build -o ./windows/stegano.exe main.go

$env:GOOS = "darwin"
$env:GOARCH = "amd64"
go build -o ./macos/amd64/stegano main.go

$env:GOOS = "darwin"
$env:GOARCH = "arm64"
go build -o ./macos/arm64/stegano main.go