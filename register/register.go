package register

import (
	"context"
	"encoding/json"
	"mongo-discovery/config"
	"mongo-discovery/etcd"
)

func NewRegistrar(ctx context.Context) (*etcd.Registrar, error) {
	service := config.GetConfig().Member.Name

	instance, err := json.Marshal(config.GetConfig().Member)
	if err != nil {
		return nil, err
	}

	// etcd
	etcdCli := etcd.NewClient(ctx, GetEtcdClient())

	registrar := etcd.NewRegistrar(etcdCli, etcd.Service{
		Prefix: service,
		Key:    config.GetConfig().Member.Name,
		Value:  string(instance),
	})
	return registrar, nil
}
