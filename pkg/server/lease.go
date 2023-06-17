package server

import (
	"context"
	"debugger/proto"
	"log"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

// CreateLease(context.Context, *CreateLeaseRequest) (*Lease, error)
func (d *DebuggerService) CreateLease(ctx context.Context, req *proto.CreateLeaseRequest) (*proto.Lease, error) {
	lease, err := d.dbConn.CreateLease(req)
	if err != nil {
		return nil, err
	}
	return lease, nil
}

// ApproveLease(context.Context, *ApproveLeaseRequest) (*Lease, error)
func (d *DebuggerService) ApproveLease(ctx context.Context, req *proto.ApproveLeaseRequest) (*proto.Lease, error) {

	// Pull the TTL from the lease
	lease, err := d.dbConn.GetLease(req.LeaseId)
	if err != nil {
		return nil, err
	}

	// Calculate expiration date
	ttl := lease.Ttl
	revocationDate := time.Now().Add(time.Second * time.Duration(ttl)).UTC()

	err = d.dbConn.UpdateLeaseStatus(req.LeaseId, "Approved", revocationDate)
	if err != nil {
		return nil, err
	}

	// Get the first pod in the deployment
	deployment, err := d.clientset.AppsV1().Deployments(lease.Namespace).Get(ctx, lease.Deployment, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	log.Println("Getting deployment")

	// Get the current pod spec
	lbls := deployment.Spec.Selector.MatchLabels
	matchedPods, err := d.clientset.CoreV1().Pods(lease.Namespace).List(ctx, metav1.ListOptions{
		LabelSelector: labels.FormatLabels(lbls),
	})
	if err != nil {
		return nil, err
	}

	log.Printf("Getting matched pod labels %+v", lbls)

	for _, targetPod := range matchedPods.Items {
		targetPod.Spec.EphemeralContainers = append(targetPod.Spec.EphemeralContainers, v1.EphemeralContainer{
			EphemeralContainerCommon: v1.EphemeralContainerCommon{
				Name:  "debug",
				Image: "busybox",
				Command: []string{
					"sleep",
					"infinity",
				},
			},
		})

		_, err = d.clientset.CoreV1().Pods(lease.Namespace).UpdateEphemeralContainers(ctx, targetPod.ObjectMeta.Name, &targetPod, metav1.UpdateOptions{})
		if err != nil {
			return nil, err
		}
	}

	log.Printf("Ephemeral container created")

	// Create Ephemeral container

	return &proto.Lease{
		LeaseId: req.LeaseId,
	}, nil
}

// ListLease(context.Context, *ListLeaseRequest) (*ListLeaseResponse, error)
func (d *DebuggerService) ListLease(ctx context.Context, req *proto.ListLeaseRequest) (*proto.ListLeaseResponse, error) {
	l, err := d.dbConn.ListLeases()
	if err != nil {
		return nil, err
	}

	return &proto.ListLeaseResponse{
		Leases: l,
	}, nil
}
