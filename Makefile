name := prnt.sc

all: clean build run

build:
	go build -v -o $(name) cmd/main.go

run:
	./$(name)

clean:
	rm -f $(name)
