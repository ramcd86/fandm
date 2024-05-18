run:
	go run ./cmd/fandm/ .

build:
	mkdir ./dist && go build -o ./dist cmd/fandm/main.go

format:
	gofmt -s -w .

clean:
	rm -rf ./dist

watch:
	nodemon --watch './**/*.go' --signal SIGKILL --exec 'go' run cmd/fandm/main.go