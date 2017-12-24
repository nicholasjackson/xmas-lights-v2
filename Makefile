generate_mocks:
	go generate ./...

build:
	GOOS=linux GOARCH=arm GOARM=6 go build -o lights_arm6
	GOOS=linux GOARCH=arm GOARM=7 go build -o lights_arm7
