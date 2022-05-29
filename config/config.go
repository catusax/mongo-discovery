package config

import (
	"fmt"
	"github.com/caarlos0/env/v6"
)

type Config struct {
	Mongo       Mongo  `envPrefix:"MONGO_"`
	Etcd        Etcd   `envPrefix:"ETCD_"`
	ClusterSize int    `env:"CLUSTER_SIZE" envDefault:"2"`
	Delete      bool   `env:"DELETE" envDefault:"false"`
	Member      Member `envPrefix:"MEMBER_"`
}

type Member struct {
	Name               string            `env:"NAME" envDefault:"${MEMBER_HOST}" envExpand:"true"`
	Host               string            `env:"HOST"`
	ArbiterOnly        bool              `env:"ARBITER_ONLY" envDefault:"false"`
	BuildIndexes       bool              `env:"BUILD_INDEXES" envDefault:"true"`
	Hidden             bool              `env:"HIDDEN" envDefault:"false"`
	Priority           int               `env:"PRIORITY" envDefault:"1"`
	SecondaryDelaySecs int               `env:"SECONDARY_DELAY_SECS" envDefault:"0"`
	Votes              int               `env:"VOTES" envDefault:"1"`
	Tags               map[string]string `envPrefix:"TAGS_"`
}

type Mongo struct {
	ReplicaName string `env:"REPLICA_NAME" envDefault:"rs"`
	UserName    string `env:"USER_NAME" envDefault:"root"`
	Password    string `env:"PASSWORD" envDefault:"example"`
	Params      string `env:"PARAMS"`
}

//Etcd describe etcd discovery server: https://github.com/etcd-io/discovery.etcd.io
type Etcd struct {
	//Endpoints are ETCD Endpoints, env: ETCD_ENDPOINTS ,multiple endpoints can be divided by comma
	Endpoints []string `env:"ENDPOINTS"`
	UserName  string   `env:"USER_NAME"`
	Password  string   `env:"PASSWORD"`
	CertPath  string   `env:"CERT_PATH"`
	KeyPath   string   `env:"KEY_PATH"`
}

var c *Config

func GetConfig() *Config {
	if c == nil {
		readConfig()
	}
	return c
}

func readConfig() { //TODO: impl
	c = &Config{}
	if err := env.Parse(c); err != nil {
		fmt.Printf("%+v\n", err)
	}
}
