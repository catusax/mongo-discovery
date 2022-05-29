package register

import (
	"crypto/tls"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
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

		var tlsConfig *tls.Config

		if config.GetConfig().Etcd.KeyPath != "" && config.GetConfig().Etcd.CertPath != "" {
			cert, err := tls.LoadX509KeyPair(config.GetConfig().Etcd.CertPath, config.GetConfig().Etcd.KeyPath)
			if err != nil {
				panic(fmt.Errorf("load cert error : %w", err))
			}
			tlsConfig = &tls.Config{
				Certificates: []tls.Certificate{cert},
			}
		}

		//var err error
		etcdClient, err = clientv3.New(
			clientv3.Config{
				Endpoints:         config.GetConfig().Etcd.Endpoints,
				DialTimeout:       ttl,
				DialKeepAliveTime: ttl,
				Username:          config.GetConfig().Etcd.UserName,
				Password:          config.GetConfig().Etcd.Password,
				TLS:               tlsConfig,
			})
		if err != nil {
			panic(fmt.Errorf("cant connect to etcd: %s: %w", config.GetConfig().Etcd.Endpoints, err))
		}
		log.Println("connected to etcd")
	})
	return etcdClient
}
