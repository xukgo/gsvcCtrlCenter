go build -ldflags -w -o out/client cmd/client/main.go;upx -9 out/*;