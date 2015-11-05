EXE=flappy

default: bin/$(EXE)


run: bin/$(EXE)
	./bin/$(EXE)

bin/$(EXE):	main.go gen/assets.go
	go build main.go
	mkdir -p bin
	mv main bin/$(EXE)

gen/assets.go: $(shell find assets -type f)
	mkdir -p gen
	go-bindata -pkg gen -o gen/assets.go assets/

clean:
	rm -rf bin gen
