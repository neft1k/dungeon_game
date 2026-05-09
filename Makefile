build:
	go build -o dungeon_game ./cmd/

test:
	go test ./...

lint:
	golangci-lint run ./...

bench:
	go test ./... -bench=. -benchmem -run=^$

clean:
	rm -f dungeon_game
