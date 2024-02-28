run:
	go run ./cmd/fandm/ .

build:
	mkdir ./dist && go build -o ./dist cmd/fandm/main.go

clean:
	rm -rf ./dist
