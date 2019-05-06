# Key generator

KeyGen REST-API is a 4-symbol key generator consisting of uppercase and lowercase letters of the Latin alphabet, as well as numbers.

## Quick start

For build run

```sh
$ make
```
for launch use

```sh
$ ./api
```

by default, the server is started on port 8080, you can change this with the “-p” parameter

```sh
$ ./api -p 3000
```

Check that everything is ok

```
$ curl -XPOST localhost:8080/key
```

if in response you got a 4 symbols code, it means that your build is ok

### Launch from Docker

```sh
$ make docker
$ docker run --rm -p {your port}:80 sidorovd/keygen
```

## API

Key generation

```sh
$ curl -XPOST localhost:8080/key
```

after key generation, it can be submitted

```sh
$ curl -XPOST localhost:8080/keys/nEw1
```

after we can get key status

```sh
$ curl localhost:8080/keys/nEw1
```

Get counts of not issued keys

```sh
$ curl localhost:8080/count
```

Health check

```sh
$ curl localhost:8080/healthcheck
```
