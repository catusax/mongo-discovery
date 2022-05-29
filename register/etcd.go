package register

import (
	"crypto/tls"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"mongo-discovery/config"
	"sync"
	"time"
)

var (
	ttl        = 5 * time.Second
	etcdOnce   sync.Once
	etcdClient *clientv3.Client
)

type EtcdClient struct {
	cli *clientv3.Client
	kv  clientv3.KV

	// leaseID will be 0 (clientv3.NoLease) if a lease was not created
	leaseID clientv3.LeaseID
}

func GetEtcdClient() *clientv3.Client {
	etcdOnce.Do(func() {
		var err error
		var cert tls.Certificate

		cert, err = tls.LoadX509KeyPair("/certs/client.pem", "/certs/client-key.pem")
		if err != nil {
			panic(fmt.Errorf("load cert error : %w", err))
		}

		//var err error
		etcdClient, err = clientv3.New(
			clientv3.Config{
				Endpoints:         config.GetConfig().Etcd.Endpoints,
				DialTimeout:       ttl,
				DialKeepAliveTime: ttl,
				Username:          config.GetConfig().Etcd.UserName,
				Password:          config.GetConfig().Etcd.Password,
				TLS: &tls.Config{
					Certificates: []tls.Certificate{cert},
				},
			})
		if err != nil {
			panic(fmt.Errorf("cant connect to etcd: %w", err))
		}
	})
	return etcdClient
}
