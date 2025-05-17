clear

echo Building Windows edition
env GOOS=windows GOARCH=amd64 go build -o builds/windows/mfw-books-db.exe

echo Building Mac edition - Intel
env GOOS=darwin GOARCH=amd64 go build -o builds/macos-intel/mfw-books-db

echo Building Mac edition - ARM - Apple Silicon
env GOOS=darwin GOARCH=arm64 go build -o builds/macos-arm/mfw-books-db

echo Building Linux edition
env GOOS=linux GOARCH=amd64 go build -o builds/linux/mfw-books-db
