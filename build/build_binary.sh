binary="baasapi"
mkdir -p dist

cd 'api/cmd/baasapi'

go get -t -d -v ./...
GOOS=$1 GOARCH=$2 CGO_ENABLED=0 go build -a --installsuffix cgo --ldflags '-s'

mv "${binary}" "../../../dist/baasapi"