build:
	protoc -I. --go_out=plugins=micro:$(GOPATH)/src/github.com/chauhanr/shipcon-consignment-service proto/consignment/consignment.proto
	GOOS=linux GOARCH=amd64 go build -o shipcon-consignment-service
	docker build -t shipcon-consignment-service .

run:
	docker run -d --net="host" \
    		-p 50052 \
    		-e MICRO_SERVER_ADDRESS=:50052 \
    		-e MICRO_REGISTRY=mdns \
    		-e DISABLE_AUTH=true \
            shipcon-consignment-service
