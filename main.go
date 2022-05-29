package main

import (
	"context"
	"log"
	"mongo-discovery/config"
	"mongo-discovery/mongo"
	"mongo-discovery/register"
	"time"
)

func main() {
	registrar, err := register.NewRegistrar(context.TODO())
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			members, err := register.GetMembers(registrar)
			if err != nil {
				log.Println(err)
				time.Sleep(time.Second * 5)
				return
			}

			if len(members) >= config.GetConfig().ClusterSize {

				err = mongo.ConfigureMongo(context.TODO(), members)
				if err != nil {
					log.Println(err)
					return
				}
				time.Sleep(time.Second * 10)
			}
		}
	}()

	registrar.Register()

}
