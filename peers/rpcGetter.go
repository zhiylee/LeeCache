package peers

import (
	pb "LeeCache/leecachepb"
	"context"
	"fmt"
	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/protocol"
)

// RPCGetter : get cache from peers by RPC, implemented Getter interface
type RPCGetter struct {
	host string
}

func (g *RPCGetter) Get(group string, key string) ([]byte, error) {
	option := client.DefaultOption
	option.SerializeType = protocol.ProtoBuffer

	c := client.NewClient(option)
	err := c.Connect("tcp", g.host)
	if err != nil {
		return nil, fmt.Errorf("rpc failed to connect: %v", err)
	}
	defer c.Close()

	args := &pb.GetRequest{
		Group: group,
		Key:   key,
	}
	reply := &pb.GetResponse{}

	err = c.Call(context.Background(), "leecache", "Get", args, reply)
	if err != nil {
		return nil, fmt.Errorf("rpc failed to call: %v", err)
	}

	return reply.Value, nil
}

func NewRPCGetter(addr string) *RPCGetter {
	return &RPCGetter{addr}
}
