export GOPATH="/tmp/go"

binary="baasapi"

mkdir -p dist
mkdir -p ${GOPATH}/src/github.com/baasapi/baasapi

cp -R api ${GOPATH}/src/github.com/baasapi/baasapi/api

cd 'api/cmd/baasapi'

go get -t -d -v ./...
GOOS=$1 GOARCH=$2 CGO_ENABLED=0 go build -a --installsuffix cgo --ldflags '-s'

mv "$BUILD_SOURCESDIRECTORY/api/cmd/baasapi/$binary" "$BUILD_SOURCESDIRECTORY/dist/baasapi"
