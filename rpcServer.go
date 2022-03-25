package leecache

import (
	pb "LeeCache/leecachepb"
	"LeeCache/status"
	"context"
	"fmt"
	"github.com/smallnest/rpcx/server"
	"log"
)

type RPCServer struct{}

func (s *RPCServer) Get(ctx context.Context, args *pb.GetRequest, reply *pb.GetResponse) error {
	group := GetGroup(args.Group)
	if group == nil {
		return fmt.Errorf("no such group %s", args.Group)
	}

	view, err := group.Get(args.Key)
	if err != err {
		return fmt.Errorf("get key=%s from group %s failed: %v", args.Key, args.Group, err)
	}
	reply.Value = view.ByteCopy()

	return nil
}

func NewRPCServer(addr string) {
	if addr == "" {
		addr = ":9998"
	}

	s := server.NewServer()
	err := s.RegisterName("leecache", new(RPCServer), "")
	if err != nil {
		log.Fatal("rpc server: ", err)
	}
	err = s.Serve("tcp", addr)
	if err != nil {
		log.Fatal("rpc server: ", err)
	}
	status.Log("rpc server is running at: %s", addr)
}
