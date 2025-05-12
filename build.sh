LDFLAGS="-X 'main.buildDate=$(date)' -X 'main.gitHash=$(git rev-parse HEAD)' -X 'main.buildOn=$(go version)' -w -s "

CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o pd2mm-windows.exe -trimpath -ldflags "${LDFLAGS}"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o pd2mm-linux -trimpath -ldflags "${LDFLAGS}"
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o pd2mm-linux-arm64 -trimpath -ldflags "${LDFLAGS}"
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o pd2mm-darwin -trimpath -ldflags "${LDFLAGS}"
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o pd2mm-darwin-arm64 -trimpath -ldflags "${LDFLAGS}"

# sha256
sha256sum pd2mm* > pd2mm-sha256
cat pd2mm-sha256

# chmod 
chmod +x pd2mm-*

# gzip
gzip --best pd2mm*