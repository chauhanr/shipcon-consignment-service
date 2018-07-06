FROM golang:1.9.0 as builder

WORKDIR /go/src/github.com/chauhanr/shipcon/consignment-service

COPY . .

RUN go get -u github.com/golang/dep/cmd/dep
RUN dep init && dep ensure
RUN CGO_ENABLED=0 GOOS=linux go build -o consignment-service -a -installsuffix cgo main.go repository.go handler.go datastore.go

FROM debian:latest

#RUN apk --no-cache add ca-certificates

RUN mkdir /app
WORKDIR /app
#ADD consignment-service /app/consignment-service
COPY --from=builder /go/src/github.com/chauhanr/shipcon/consignment-service/consignment-service .

CMD ["./consignment-service"]