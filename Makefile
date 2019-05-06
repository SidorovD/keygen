all:
	go build -o api cmd/keygenapi/main.go

docker:
	docker build -t sidorovd/keygen .

clean:
	rm api