# use alpine cause it light weigth
FROM golang:1.12-alpine as build
WORKDIR /go/src/keygen
COPY . .
# we should turn off the cgo to run from scratch
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o api cmd/keygenapi/main.go

FROM scratch
EXPOSE 80
COPY --from=build /go/src/keygen/api .
CMD [ "./api", "-p", "80" ]
