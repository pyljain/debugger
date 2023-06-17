package db

import (
	"database/sql"
	"debugger/proto"
	"time"

	_ "github.com/lib/pq"
)

type Postgres struct {
	conn *sql.DB
}

func NewPostgres(connectionString string) (*Postgres, error) {
	conn, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	db := &Postgres{
		conn: conn,
	}

	err = db.createTables()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (p *Postgres) createTables() error {
	_, err := p.conn.Exec(`
	CREATE TABLE IF NOT EXISTS leases (
		id serial primary key,
		deployment VARCHAR(100) not null,
		ttl NUMERIC DEFAULT 240,
		revocation_date TIMESTAMP,
		status VARCHAR(40) DEFAULT 'Requested',
		namespace VARCHAR(100) not null
	);
	`)
	if err != nil {
		return err
	}

	return nil
}

func (pg *Postgres) CreateLease(req *proto.CreateLeaseRequest) (*proto.Lease, error) {

	var id int

	err := pg.conn.QueryRow(`
		INSERT INTO leases (deployment, namespace, ttl, status) VALUES ($1, $2, $3, $4)
		RETURNING id
	`, req.Deployment, req.Namespace, req.Ttl, "Requested").Scan(&id)
	if err != nil {
		return nil, err
	}

	return &proto.Lease{
		LeaseId:    int32(id),
		Deployment: req.Deployment,
		Ttl:        req.Ttl,
		Status:     "Requested",
	}, nil

}

func (pg *Postgres) UpdateLeaseStatus(leaseId int32, status string, revocationDate time.Time) error {
	row := pg.conn.QueryRow(`
		UPDATE leases SET STATUS=$1, REVOCATION_DATE=$3 WHERE ID=$2
	`, status, leaseId, revocationDate)
	if row.Err() != nil {
		return row.Err()
	}

	return nil
}

func (pg *Postgres) GetLease(leaseId int32) (*proto.Lease, error) {
	result := &proto.Lease{}
	row := pg.conn.QueryRow(`
		SELECT Id, Deployment, Namespace, TTL, Status from leases
		WHERE id = $1
	`, leaseId)
	if row.Err() != nil {
		return nil, row.Err()
	}

	row.Scan(&result.LeaseId, &result.Deployment, &result.Namespace, &result.Ttl, &result.Status)

	return result, nil
}

func (pg *Postgres) ListLeases() ([]*proto.Lease, error) {
	rows, err := pg.conn.Query(`
		SELECT ID, STATUS, DEPLOYMENT, NAMESPACE, TTL 
			FROM leases
	`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	leases := []*proto.Lease{}
	for rows.Next() {
		row := &proto.Lease{}
		rows.Scan(&row.LeaseId, &row.Status, &row.Deployment, &row.Namespace, &row.Ttl)
		leases = append(leases, row)
	}

	return leases, nil
}

func (pg *Postgres) GetExpiredLeases() ([]*proto.Lease, error) {
	rows, err := pg.conn.Query(`
		SELECT ID, STATUS, DEPLOYMENT, NAMESPACE, TTL
			FROM leases
			WHERE STATUS='Approved' 
			AND
			REVOCATION_DATE < NOW()
	`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	leases := []*proto.Lease{}
	for rows.Next() {
		row := &proto.Lease{}
		rows.Scan(&row.LeaseId, &row.Status, &row.Deployment, &row.Namespace, &row.Ttl)
		leases = append(leases, row)
	}

	return leases, nil
}
