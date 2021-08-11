build:
	go mod download && CGO_ENABLED=0 GOOS=linux go build -o .bin/app ./cmd 

run:
	.bin/./app
	
br:
	make build && make run

clean:
	rm -rf .bin

default:
	build