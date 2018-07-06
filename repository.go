package main

import(
	pb "github.com/chauhanr/shipcon-consignment-service/proto/consignment"
    "gopkg.in/mgo.v2"
)

const(
	dbName = "shipcon"
	consignmentCollection = "consignments"
)

type Repository interface{
	Create(*pb.Consignment) error
	GetAll()([]*pb.Consignment, error)
	Close()
}

type ConsignmentRepository struct{
	session *mgo.Session
}

func (repo *ConsignmentRepository) GetAll() ([]*pb.Consignment, error) {
	var consignments []*pb.Consignment
	// query the data base and get the values for consigement
	err := repo.collection().Find(nil).All(&consignments)
	return consignments, err
}

/**
	Function to close the repository mongo session
	Mgo creates a master session at startup. A good practice is to create a new session from the master each time
    a new request needs to be servies. This is safer and more efficient to do so and management session. Therefore
	it becomes very important to close the sessions as well.
*/
func (repo *ConsignmentRepository) Close(){
	repo.session.Close()
}

func (repo *ConsignmentRepository) Create(consignment *pb.Consignment) error{
	return repo.collection().Insert(consignment)
}

func (repo *ConsignmentRepository) collection() *mgo.Collection{
	return repo.session.DB(dbName).C(consignmentCollection)
}
