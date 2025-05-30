SET LDFLAGS="-X 'main.buildDate=$(date)' -X 'main.gitHash=$(git rev-parse HEAD)' -X 'main.buildOn=$(go version)' -w -s -H=windowsgui -extldflags=-static"

SET CGO_ENABLED=1
SET GOOS=linux
SET GOARCH=amd64
go build -o pd2mm-linux -trimpath -ldflags %LDFLAGS% ./frontend/gui

SET CGO_ENABLED=1
SET GOOS=linux
SET GOARCH=arm64
go build -o pd2mm-linux-arm64 -trimpath -ldflags %LDFLAGS% ./frontend/gui

SET CGO_ENABLED=1
SET GOOS=darwin
SET GOARCH=amd64
go build -o pd2mm-darwin -trimpath -ldflags %LDFLAGS% ./frontend/gui

SET CGO_ENABLED=1
SET GOOS=darwin
SET GOARCH=arm64
go build -o pd2mm-darwin-arm64 -trimpath -ldflags %LDFLAGS% ./frontend/gui

SET CGO_ENABLED=1
SET GOOS=windows
SET GOARCH=amd64
go build -o pd2mm-windows.exe -trimpath -ldflags %LDFLAGS% ./frontend/gui