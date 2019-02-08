/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package latency

import (
	"github.com/gocql/gocql"
	"github.com/nalej/derrors"
	"github.com/nalej/device-manager/internal/pkg/entities"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/rs/zerolog/log"
	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"
	"sync"
	"time"
)
// TTL -> 1 day. After 24 hours, the record will be deleted
const ttlExpired = time.Duration(24)*time.Hour

type ScyllaProvider struct {
	Address string
	Port int
	Keyspace string
	Session *gocql.Session
	sync.Mutex
}

func NewScyllaProvider (address string, port int, keyspace string) * ScyllaProvider {
	provider := ScyllaProvider{Address:address, Port: port, Keyspace: keyspace, Session: nil}
	provider.connect()
	return &provider
}

func(sp * ScyllaProvider) connect() derrors.Error{
	// connect to the cluster
	conf := gocql.NewCluster(sp.Address)
	conf.Keyspace = sp.Keyspace
	conf.Port = sp.Port

	session, err := conf.CreateSession()
	if err != nil {
		log.Error().Str("provider", "ScyllaProvider").Str("trace", conversions.ToDerror(err).DebugReport()).Msg("unable to connect")
		return derrors.AsError(err, "cannot connect")
	}

	sp.Session = session

	return nil
}

func (sp *ScyllaProvider) checkAndConnect () derrors.Error{

	if sp.Session == nil{
		log.Info().Msg("session no created, trying to reconnect...")
		// try to reconnect
		err := sp.connect()
		if err != nil  {
			return err
		}
	}
	return nil
}

func (sp *ScyllaProvider) Disconnect () {

	sp.Lock()
	defer sp.Unlock()

	if sp.Session != nil {
		sp.Session.Close()
		sp.Session = nil
	}
}

func (sp * ScyllaProvider) AddPingLatency(latency entities.Latency ) derrors.Error {

	sp.Lock()
	defer sp.Unlock()

	// check connection
	err := sp.checkAndConnect()
	if err != nil {
		return err
	}


	// insert the application instance
	stmt, names := qb.Insert("latency").Columns("organization_id","device_group_id", "device_id",
		"inserted","latency").TTL(ttlExpired).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindStruct(latency)
	cqlErr := q.ExecRelease()

	if cqlErr != nil {
		return derrors.AsError(cqlErr, "cannot add latency")
	}

	return nil
}



