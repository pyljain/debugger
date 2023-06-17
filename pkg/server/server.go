package server

import (
	"debugger/pkg/db"
	"debugger/pkg/revoker"
	"debugger/proto"
	"net"

	"google.golang.org/grpc"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type DebuggerService struct {
	proto.UnimplementedDebuggerServer
	dbConn    db.DB
	clientset *kubernetes.Clientset
}

func Start(connectionString string, kubeconfig string) error {
	listener, err := net.Listen("tcp", ":9091")
	if err != nil {
		return err
	}

	dbConn, err := db.NewPostgres(connectionString)
	if err != nil {
		return err
	}

	s := grpc.NewServer()

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}
	go revoker.Start(dbConn, clientset)

	d := DebuggerService{
		dbConn:    dbConn,
		clientset: clientset,
	}
	proto.RegisterDebuggerServer(s, &d)

	err = s.Serve(listener)
	if err != nil {
		return err
	}

	return nil
}

// CreateLease(context.Context, *CreateLeaseRequest) (*Lease, error)
// ApproveLease(context.Context, *ApproveLeaseRequest) (*Lease, error)
