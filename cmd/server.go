package main

import (
	"context"
	"flag"
	"github.com/jmoiron/sqlx"
	proto "github.com/m0cchi/gfalcon-internal-grpc/pb/proto"
	"github.com/m0cchi/gfalcon/complex"
	"github.com/m0cchi/gfalcon/model"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

var db *sqlx.DB

type gfalcon struct {
}

func (s *gfalcon) SignIn(ctx context.Context, req *proto.SignInRequest) (*proto.SignInResponse, error) {

	team, err := model.GetTeam(db, req.Team)
	if err != nil {
		return nil, err
	}
	user, err := model.GetUser(db, team.IID, req.Id)
	if err != nil {
		return nil, err
	}
	session, err := complex.AuthenticateWithPassword(db, user, req.Password)
	if err != nil {
		return nil, err
	}
	if err = session.Validate(); err != nil {
		return nil, err
	}

	res := &proto.SignInResponse{
		Ok:      true,
		Iid:     user.IID,
		Session: session.SessionID,
	}
	return res, nil
}

func (s *gfalcon) Check(ctx context.Context, req *proto.CheckRequest) (*proto.CheckResponse, error) {
	session, err := model.GetSession(db, req.Iid, req.Session)
	if err != nil {
		return nil, err
	}
	if err = session.Validate(); err != nil {
		return nil, err
	}

	res := &proto.CheckResponse{
		Ok: true,
	}
	return res, nil
}

func main() {
	log.Println("start")

	var dbhost string
	var err error
	flag.StringVar(&dbhost, "dbhost", "", "gfalcon's DB")
	flag.Parse()

	if dbhost == "" {
		log.Println("required --dbhost [host]")
		os.Exit(1)
	}

	db, err = sqlx.Connect("mysql", dbhost)
	if err != nil {
		log.Printf("err: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	listenPort, err := net.Listen("tcp", ":19003")
	if err != nil {
		log.Fatalln(err)
	}
	server := grpc.NewServer()
	proto.RegisterGfalconServer(server, &gfalcon{})
	server.Serve(listenPort)
}
