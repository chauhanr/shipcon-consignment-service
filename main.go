package main


import (

	pb "github.com/chauhanr/shipcon-consignment-service/proto/consignment"
	vesselProto "github.com/chauhanr/shipcon-vessel-service/proto/vessel"
		"github.com/micro/go-micro"
	"github.com/micro/go-micro/server"
	"fmt"
	"log"
	"os"
	"context"
	"github.com/micro/go-micro/metadata"
	"github.com/pkg/errors"
	userService "github.com/chauhanr/shipcon-user-service/proto/user"
	)


const(
	defaultHost = "localhost:27017"
)

var (
	srv micro.Service
)

func main(){
	host := os.Getenv("DB_HOST")
	if host == ""{
		host = defaultHost
	}

	session, err := CreateSession(host)
	defer session.Close()

	if err != nil{
		log.Panicf("Could not connect to data store with host %s - %v", host, err)
	}

	srv = micro.NewService(
		micro.Name("go.micro.srv.consignment"),
		micro.Version("latest"),
	)

	vesselClient := vesselProto.NewVesselServiceClient("go.micro.srv.vessel", srv.Client())
	srv.Init(micro.WrapHandler(AuthWrapper))
	pb.RegisterShippingServiceHandler(srv.Server(),&service{session, vesselClient})

	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}

}

func AuthWrapper(fn server.HandlerFunc) server.HandlerFunc{
	return func(ctx context.Context, req server.Request, resp interface{}) error{
		if os.Getenv("DISABLE_AUTH") == "true" {
			return  fn(ctx,req,resp)//fn(ctx, req, ctx)
		}
		meta, ok := metadata.FromContext(ctx)
		if !ok{
			return errors.New("no auth meta data found in request")
		}
		token := meta["Token"]
		log.Println("Authenticating with token: ", token)
		// Auth here
		if srv != nil{
			authClient := userService.NewUserServiceClient("", srv.Client())
			if authClient != nil{
				_, err := authClient.ValidateToken(ctx, &userService.Token{
					Token: token,
				})
				if err != nil{
					log.Printf("Error authenticating : %s\n", err)
					return err
				}
				return fn(ctx, req, resp)
			}else{
				log.Println("Authentication client could not be initialized")
				return errors.New("Authentication client could not be initialized")
			}
		}else {
			log.Println("The micro.Service instance is nil so going unauthenticated")
		}
		return fn(ctx,req,resp)
	}
}

