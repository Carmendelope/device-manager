create KEYSPACE IF NOT EXISTS measure WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};

Create table IF NOT EXISTS measure.latency (organization_id text, device_group_id text, device_id text, inserted bigint, latency int, PRIMARY KEY ((organization_id, device_group_id), device_id, inserted) );

Create materialized view IF NOT EXISTS measure.deviceGrouplatency as  select * from measure.latency where organization_id is not null and device_group_id is not null and inserted is not null and device_id is not null primary key ((organization_id, device_group_id), inserted, device_id);
