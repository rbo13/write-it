go:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-extldflags '-static'" -tags netgo -a -v -o write-it main.go


darwin:
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-extldflags '-static'" -tags netgo -a -v -o write-it main.go
