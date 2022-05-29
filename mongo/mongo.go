package mongo

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"mongo-discovery/config"
	"reflect"
)

func ConfigureMongo(ctx context.Context, members []config.Member) error {

	//connect to current server
	uri := fmt.Sprintf("mongodb://%s/?replicaSet=%s", config.GetConfig().Member.Host, config.GetConfig().Mongo.ReplicaName)

	if config.GetConfig().Mongo.Params != "" {
		uri = fmt.Sprintf("%s&%s", uri, config.GetConfig().Mongo.Params)
	}

	//TODO:TLS Config

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri).SetAuth(
		options.Credential{
			AuthMechanism:           "",
			AuthMechanismProperties: nil,
			AuthSource:              "",
			Username:                config.GetConfig().Mongo.UserName,
			Password:                config.GetConfig().Mongo.Password,
			PasswordSet:             false,
		},
	))
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)

	coll := client.Database("local").Collection("system.replset")

	rep := coll.FindOne(ctx, bson.M{"_id": config.GetConfig().Mongo.ReplicaName})

	if errors.Is(rep.Err(), mongo.ErrNoDocuments) { //do not exist,create rep config

		membersBytes, err := bson.Marshal(members)
		if err != nil {
			return fmt.Errorf("marshal members %w", err)
		}

		err = client.Database("local").RunCommand(ctx,
			fmt.Sprintf("rs.initiate({_id:\"%s\", members: %s})",
				config.GetConfig().Mongo.ReplicaName, membersBytes),
		).Err()
		if err != nil {
			return fmt.Errorf("init replicaset %w", err)
		}
		return nil
	} else if err != nil {
		return fmt.Errorf("query replicaset %w", err)
	}

	bsonRaw, err := rep.DecodeBytes()
	if err != nil {
		return fmt.Errorf("decode replica conf %w", err)
	}
	membersRaw := bsonRaw.Lookup("members")

	var oldMembers []config.Member
	err = membersRaw.Unmarshal(&oldMembers)
	if err != nil {
		return fmt.Errorf("unmarshal members %w", err)
	}

	//add members
	for _, newMember := range members {

		oldMember, exist := memberExist(newMember, oldMembers)
		if !exist {
			membersBytes, err := bson.Marshal(newMember)
			if err != nil {
				return fmt.Errorf("marshal members %w", err)
			}
			err = client.Database("local").RunCommand(ctx,
				fmt.Sprintf("rs.add(%s)",
					membersBytes),
			).Err()
			if err != nil {
				return fmt.Errorf("add member %w", err)
			}
		} else {
			//TODO: change member config
			if !reflect.DeepEqual(newMember, oldMember) {

			}
		}
	}

	//remove members
	if config.GetConfig().Delete {
		for _, oldMember := range oldMembers {
			_, exist := memberExist(oldMember, members)
			if !exist {
				err = client.Database("local").RunCommand(ctx,
					fmt.Sprintf("rs.remove(\"%s\")",
						oldMember.Host),
				).Err()
				if err != nil {
					return fmt.Errorf("remove member %w", err)
				}
			}
		}
	}

	return nil
}

func memberExist(member config.Member, members []config.Member) (*config.Member, bool) {
	for _, m := range members {
		if m.Host == member.Host {
			return &m, true
		}
	}
	return nil, false
}
