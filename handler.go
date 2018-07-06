package main

import ("gopkg.in/mgo.v2"
	vesselProto "github.com/chauhanr/shipcon/vessel-service/proto/vessel"
	pb "github.com/chauhanr/shipcon/consignment-service/proto/consignment"
	"context"
	"log"
)

type service struct{
	session *mgo.Session
	vesselClient vesselProto.VesselServiceClient
}

func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment, res *pb.Response) error{
	defer s.GetRepo().Close()

	vesselResponse, err := s.vesselClient.FindAvailable(context.Background(), &vesselProto.Specification{
		MaxWeight: req.Weight,
		Capacity: int32(len(req.Containers)),
	})

	log.Printf("Found contaoiner %v\n", vesselResponse.Vessel)
	if err != nil{
		return err
	}

	req.VesselId = vesselResponse.Vessel.Id
	err = s.GetRepo().Create(req)
	if err != nil{
		return err
	}
	res.Created= true
	res.Consignment = req

	return nil
}

func (s *service) GetConsignments(ctx context.Context, req *pb.GetRequest, res *pb.Response) error{
	defer s.GetRepo().Close()
	consignments, err := s.GetRepo().GetAll()

	if err != nil{
		return err
	}
	res.Consignments = consignments
	return nil
}


func (s *service) GetRepo() Repository{
	return &ConsignmentRepository{s.session.Clone()}
}

