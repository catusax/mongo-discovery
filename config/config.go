package config

import (
	"fmt"
	"github.com/caarlos0/env/v6"
)

type Config struct {
	Mongo Mongo `envPrefix:"MONGO_"`
	Etcd  Etcd  `envPrefix:"ETCD_"`
	// ClusterSize is the minimum member nums required to init replicaSet
	ClusterSize int `env:"CLUSTER_SIZE" envDefault:"2"`
	// Delete set true to automatically remove member when its down
	Delete bool   `env:"DELETE" envDefault:"false"`
	Member Member `envPrefix:"MEMBER_"`
}

// Member is mongo replica-configuration member, doc: https://www.mongodb.com/docs/manual/reference/replica-configuration/
type Member struct {
	ID int `bson:"_id" json:"id,omitempty"`
	//Name               string ` env:"NAME" envDefault:"${MEMBER_HOST}" envExpand:"true" bson:"name"`
	// Host hostname of this member ip:port or dns:port
	Host               string `bson:"host" json:"host,omitempty" env:"HOST"`
	ArbiterOnly        bool   `bson:"arbiterOnly" json:"arbiter_only,omitempty" env:"ARBITER_ONLY" envDefault:"false" `
	BuildIndexes       bool   `bson:"buildIndexes" json:"build_indexes,omitempty" env:"BUILD_INDEXES" envDefault:"true" `
	Hidden             bool   `bson:"hidden" json:"hidden,omitempty" env:"HIDDEN" envDefault:"false" `
	Priority           int    `bson:"priority" json:"priority,omitempty" env:"PRIORITY" envDefault:"1" `
	SecondaryDelaySecs int    `bson:"secondaryDelaySecs" json:"secondary_delay_secs,omitempty" env:"SECONDARY_DELAY_SECS" envDefault:"0" `
	Votes              int    `bson:"votes" json:"votes,omitempty" env:"VOTES" envDefault:"1" `
	//Tags               map[string]string
}

// Mongo config
type Mongo struct {
	// Host to init replicaset, ip:port or dns:port. by default, it's set to read MEMBER_HOST
	Host string `env:"HOST" envDefault:"${MEMBER_HOST}" envExpand:"true"`
	//ReplicaName name of replicaset
	ReplicaName string `env:"REPLICA_NAME" envDefault:"rs"`
	UserName    string `env:"USER_NAME" envDefault:"root"`
	Password    string `env:"PASSWORD" envDefault:"example"`
	//uri params for mongodb, eg: w=majority&wtimeoutMS=2000, you can use Params to set certificate
	Params string `env:"PARAMS"`
}

//Etcd registry
type Etcd struct {
	//Endpoints are ETCD Endpoints, env: ETCD_ENDPOINTS ,multiple endpoints can be divided by comma, eg: http://172.28.232.235:3379
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
