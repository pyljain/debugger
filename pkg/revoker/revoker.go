package revoker

import (
	"context"
	"debugger/pkg/db"
	"log"
	"time"

	"golang.org/x/sync/errgroup"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

func Start(dbConn db.DB, cs *kubernetes.Clientset) {
	for {
		time.Sleep(1 * time.Second)
		expiredLeases, err := dbConn.GetExpiredLeases()
		if err != nil {
			log.Printf("Unable to fetch expired leases %s", err)
		}

		log.Printf("Expried lease %v", expiredLeases)

		// Loop through leases to process and start error group
		eg := errgroup.Group{}
		ctx := context.Background()
		for _, l := range expiredLeases {
			lease := l
			eg.Go(func() error {
				deployment, err := cs.AppsV1().Deployments(lease.Namespace).Get(ctx, lease.Deployment, metav1.GetOptions{})
				if err != nil {
					return err
				}

				lbls := deployment.Spec.Selector.MatchLabels
				matchedPods, err := cs.CoreV1().Pods(lease.Namespace).List(ctx, metav1.ListOptions{
					LabelSelector: labels.FormatLabels(lbls),
				})
				if err != nil {
					return err
				}

				for _, p := range matchedPods.Items {

					err := cs.CoreV1().Pods(lease.Namespace).Delete(ctx, p.Name, metav1.DeleteOptions{})
					if err != nil {
						return err
					}
				}

				// Update status in the database
				err = dbConn.UpdateLeaseStatus(lease.LeaseId, "Expired", time.Now().UTC())
				if err != nil {
					return err
				}

				return nil
			})
		}

		err = eg.Wait()
		if err != nil {
			log.Printf("Error occured while deleting ephemeral containers %s", err)
		}

		// Call k8s APIs to revoke the ephemeral container

		// If error notify

		// Update status of lease once completed
	}
}
