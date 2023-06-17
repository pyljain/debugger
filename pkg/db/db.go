package db

import (
	"debugger/proto"
	"time"
)

type DB interface {
	CreateLease(*proto.CreateLeaseRequest) (*proto.Lease, error)
	UpdateLeaseStatus(leaseId int32, status string, revocationDate time.Time) error
	ListLeases() ([]*proto.Lease, error)
	GetLease(leaseId int32) (*proto.Lease, error)
	GetExpiredLeases() ([]*proto.Lease, error)
}
