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

 /*
 Steps to test:
 -------------
 1) launch docker image
 docker run --name scylla -p 9042:9042 -d scylladb/scylla

 2) create keyspace and table
 docker exec -it scylla cqlsh

 create KEYSPACE IF NOT EXISTS measure WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};
 create table IF NOT EXISTS measure.latency (organization_id text, device_group_id text, device_id text, inserted bigint, latency int, PRIMARY KEY ((organization_id, device_group_id, device_id), inserted) );

 3)environment variables:
 RUN_INTEGRATION_TEST=true
 IT_SCYLLA_HOST=127.0.0.1
 IT_SCYLLA_PORT=9042
 IT_NALEJ_KEYSPACE=measure

  */
package latency

import (
	"github.com/onsi/ginkgo"
	"github.com/rs/zerolog/log"
	"github.com/nalej/device-manager/internal/pkg/utils"

	"os"
	"strconv"
)

var _ = ginkgo.Describe("Scylla latency provider", func(){


	if ! utils.RunIntegrationTests() {
		log.Warn().Msg("Integration tests are skipped")
		return
	}

	var scyllaHost = os.Getenv("IT_SCYLLA_HOST")
	if scyllaHost == "" {
		ginkgo.Fail("missing environment variables")
	}

	var nalejKeySpace = os.Getenv("IT_NALEJ_KEYSPACE")
	if nalejKeySpace == "" {
		ginkgo.Fail("missing environment variables")

	}
	scyllaPort, _ := strconv.Atoi(os.Getenv("IT_SCYLLA_PORT"))
	if scyllaPort <= 0 {
		ginkgo.Fail("missing environment variables")

	}

	// create a provider and connect it
	sp := NewScyllaProvider(scyllaHost, scyllaPort, nalejKeySpace)

	// disconnect
	ginkgo.AfterSuite(func() {
		sp.Disconnect()
	})

	RunTest(sp)

})
