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

	log.Println("registering...")
	registrar.Register()

	log.Println("running...")
	for {
		log.Println("requesting member list")
		members, err := register.GetMembers(registrar)
		if err != nil {
			log.Println(err)
			time.Sleep(time.Second * 5)
		}
		log.Println("got member list", "size", len(members))

		if len(members) >= config.GetConfig().ClusterSize {

			err = mongo.ConfigureMongo(context.TODO(), members)
			if err != nil {
				log.Println(err)
			}
			time.Sleep(time.Second * 10)
		} else {
			log.Println("cluster not ready", "size", len(members))
			time.Sleep(time.Second * 1)
		}
	}

}
