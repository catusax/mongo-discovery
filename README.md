# mongo-discovery

discover mongodb replicaSet members.

## usage

see [example](example)

1. start an etcd server
2. start mongodb alongside with mongo-discovery
3. when the members count is bigger than `CLUSTER_SIZE`,mongo-discovery will try to init mongodb replicaset with all
   members.
4. start/stop/reconfigure container to add/remove/reconfigure member

## config

see [configs](config/config.go) for all available envs

### required config

- `MEMBER_NAME` name of this replicaset member
- `MEMBER_HOST` host/ip:port of this replicaset member, eg:172.28.232.235:27017
- `ETCD_ENDPOINTS` etcd discovery endpoints eg: http://172.28.232.235:3379,http://172.28.232.235:3379
- `MONGO_REPLICA_NAME` name of replicaset
- `MONGO_USER_NAME` mongodb admin user name
- `MONGO_PASSWORD` mongodb admin user password