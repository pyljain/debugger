package client

import (
	"debugger/proto"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	proto.DebuggerClient
}

func New(address string) (*Client, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}

	log.Println("In client.New")

	// defer conn.Close()
	clientConn := proto.NewDebuggerClient(conn)

	log.Println("In client.New after clientConn")

	return &Client{
		clientConn,
	}, nil
}
