/*
 * Copyright 2019 Nalej
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
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
const ttlExpired = time.Duration(24) * time.Hour
const rowNotFound = "not found"

const limitTime = time.Duration(5) * time.Minute

type ScyllaProvider struct {
	Address  string
	Port     int
	Keyspace string
	Session  *gocql.Session
	sync.Mutex
}

func NewScyllaProvider(address string, port int, keyspace string) *ScyllaProvider {
	provider := ScyllaProvider{Address: address, Port: port, Keyspace: keyspace, Session: nil}
	provider.connect()
	return &provider
}

func (sp *ScyllaProvider) connect() derrors.Error {
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

func (sp *ScyllaProvider) checkAndConnect() derrors.Error {

	if sp.Session == nil {
		log.Info().Msg("session no created, trying to reconnect...")
		// try to reconnect
		err := sp.connect()
		if err != nil {
			return err
		}
	}
	return nil
}

func (sp *ScyllaProvider) Disconnect() {

	sp.Lock()
	defer sp.Unlock()

	if sp.Session != nil {
		sp.Session.Close()
		sp.Session = nil
	}
}

// -- Latency
func (sp *ScyllaProvider) AddPingLatency(latency entities.Latency) derrors.Error {

	sp.Lock()
	defer sp.Unlock()

	// check connection
	err := sp.checkAndConnect()
	if err != nil {
		return err
	}

	// insert the application instance
	stmt, names := qb.Insert("latency").Columns("organization_id", "device_group_id", "device_id",
		"inserted", "latency").TTL(ttlExpired).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindStruct(latency)
	cqlErr := q.ExecRelease()

	if cqlErr != nil {
		return derrors.AsError(cqlErr, "cannot add latency")
	}

	return nil
}

func (sp *ScyllaProvider) GetLatency(organizationID string, deviceGroupID string, deviceID string) ([]*entities.Latency, derrors.Error) {
	sp.Lock()
	defer sp.Unlock()

	// check connection
	err := sp.checkAndConnect()
	if err != nil {
		return nil, err
	}

	latencyList := make([]*entities.Latency, 0)
	stmt, names := qb.Select("latency").Where(qb.Eq("organization_id")).
		Where(qb.Eq("device_group_id")).Where(qb.Eq("device_id")).
		OrderBy("device_id", qb.DESC).OrderBy("inserted", qb.DESC).ToCql()

	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		"organization_id": organizationID,
		"device_group_id": deviceGroupID,
		"device_id":       deviceID,
	})

	cqlErr := gocqlx.Select(&latencyList, q.Query)

	if cqlErr != nil {
		if cqlErr.Error() == rowNotFound {
			return latencyList, nil
		} else {
			return nil, derrors.AsError(cqlErr, "cannot list device latencies")
		}
	}

	return latencyList, nil
}

func (sp *ScyllaProvider) RemoveLatency(organizationID string, deviceGroupID string, deviceID string) derrors.Error {

	sp.Lock()
	defer sp.Unlock()

	// check connection
	err := sp.checkAndConnect()
	if err != nil {
		return err
	}

	stmt, _ := qb.Delete("latency").Where(qb.Eq("organization_id")).Where(qb.Eq("device_group_id")).Where(qb.Eq("device_id")).ToCql()
	cqlErr := sp.Session.Query(stmt, organizationID, deviceGroupID, deviceID).Exec()

	if cqlErr != nil {
		return derrors.AsError(cqlErr, "cannot delete device group")
	}

	return nil

}

// -- Last Latency
// GetLastPingLatency get the las latency measure of a device
func (sp *ScyllaProvider) GetLastLatency(organizationID string, deviceGroupID string, deviceID string) (*entities.Latency, derrors.Error) {

	sp.Lock()
	defer sp.Unlock()

	// check connection
	err := sp.checkAndConnect()
	if err != nil {
		return nil, err
	}

	var latency entities.Latency

	// get the lastPingLatency from LastLatency table
	stmt, names := qb.Select("lastlatency").Where(qb.Eq("organization_id")).
		Where(qb.Eq("device_group_id")).Where(qb.Eq("device_id")).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		"organization_id": organizationID,
		"device_group_id": deviceGroupID,
		"device_id":       deviceID,
	})

	cqlErr := q.GetRelease(&latency)
	if cqlErr != nil {
		if cqlErr.Error() == rowNotFound {
			return entities.NewEmptyLatency(), nil
		} else {
			return nil, derrors.AsError(err, "cannot Cannot retrieve last latency")
		}
	}

	return &latency, nil

}

func (sp *ScyllaProvider) AddLastLatency(latency entities.Latency) derrors.Error {

	sp.Lock()
	defer sp.Unlock()

	// check connection
	err := sp.checkAndConnect()
	if err != nil {
		return err
	}

	// insert the application instance
	stmt, names := qb.Insert("lastlatency").Columns("organization_id", "device_group_id", "device_id",
		"inserted", "latency").TTL(ttlExpired).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindStruct(latency)
	cqlErr := q.ExecRelease()

	if cqlErr != nil {
		return derrors.AsError(cqlErr, "cannot add last latency")
	}

	return nil
}

func (sp *ScyllaProvider) GetGroupLastLatencies(organizationID string, deviceGroupID string) ([]*entities.Latency, derrors.Error) {
	sp.Lock()
	defer sp.Unlock()

	// check connection
	err := sp.checkAndConnect()
	if err != nil {
		return nil, err
	}

	latencyList := make([]*entities.Latency, 0)
	stmt, names := qb.Select("lastlatency").Where(qb.Eq("organization_id")).
		Where(qb.Eq("device_group_id")).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		"organization_id": organizationID,
		"device_group_id": deviceGroupID,
	})

	cqlErr := gocqlx.Select(&latencyList, q.Query)

	if cqlErr != nil {
		if cqlErr.Error() == rowNotFound {
			return latencyList, nil
		} else {
			return nil, derrors.AsError(cqlErr, "cannot list group latencies")
		}
	}

	return latencyList, nil
}
