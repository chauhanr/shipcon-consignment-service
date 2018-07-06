package main


import (

	pb "github.com/chauhanr/shipcon-consignment-service/proto/consignment"
	vesselProto "github.com/chauhanr/shipcon-vessel-service/proto/vessel"
		"github.com/micro/go-micro"
	"fmt"
	"log"
	"os"
)


const(
	defaultHost = "localhost:27017"
)

func main(){
	host := os.Getenv("DB_HOST")
	if host == ""{
		host = defaultHost
	}

	session, err := CreateSession(host)
	defer session.Close()

	if err != nil{
		log.Panic("Could not connect to data store with host %s - %v", host, err)
	}

	srv := micro.NewService(
		micro.Name("go.micro.srv.consignment"),
		micro.Version("latest"),
	)

	vesselClient := vesselProto.NewVesselServiceClient("go.micro.srv.vessel", srv.Client())
	srv.Init()
	pb.RegisterShippingServiceHandler(srv.Server(),&service{session, vesselClient})

	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}

}

