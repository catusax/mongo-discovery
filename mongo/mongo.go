package mongo

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"mongo-discovery/config"
	"sort"
)

func ConfigureMongo(ctx context.Context, members []config.Member) error {

	//connect to current server
	uri := fmt.Sprintf("mongodb://%s/admin?connect=direct", config.GetConfig().Mongo.Host)

	if config.GetConfig().Mongo.Params != "" {
		uri = fmt.Sprintf("%s&%s", uri, config.GetConfig().Mongo.Params)
	}

	//TODO:TLS Config

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri).SetAuth(
		options.Credential{
			Username: config.GetConfig().Mongo.UserName,
			Password: config.GetConfig().Mongo.Password,
		},
	))

	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)
	log.Println("connected mongodb")

	db := client.Database("admin")

	var rep bson.M
	err = db.RunCommand(ctx, bson.D{{"replSetGetConfig", 1}}).Decode(&rep)
	if err != nil {
		fmt.Errorf("error get replSetGetConfig %w", err)
	}

	okValue, _ := rep["ok"].(float64)

	if err != nil || okValue == 0.0 { //do not exist,create rep config
		fmt.Println("init replicaset...")

		for i := 0; i < len(members); i++ {
			members[i].ID = i + 1
		}

		var initresp bson.M
		err = db.RunCommand(ctx, bson.D{{"replSetInitiate", bson.M{
			"_id":     config.GetConfig().Mongo.ReplicaName,
			"members": members,
		}}}).Decode(&initresp)

		if err != nil {
			return fmt.Errorf("init replicaset %w", err)
		}

		bytes, err := bson.Marshal(initresp)
		fmt.Sprintf("%s,%v", bytes, err)

		return nil
	} //else if err != nil {
	//	return fmt.Errorf("query replicaset %w", err)
	//}

	repConfig, ok := rep["config"].(bson.M)
	if !ok {
		return errors.New("replicaset config do not exist")
	}

	oldRaw, ok := repConfig["members"].(bson.A)
	if !ok {
		return errors.New("replicaset config members do not exist")
	}

	oldBytes, err := bson.Marshal(oldRaw)
	if err != nil {
		return fmt.Errorf("marshal old members")
	}

	var oldMembers []config.Member

	err = bson.Unmarshal(oldBytes, &oldMembers)
	if err != nil {
		return fmt.Errorf("unmarshal members %w", err)
	}

	//remove members
	if config.GetConfig().Delete {
		for i, newMember := range members {
			j, exist := memberExist(newMember, oldMembers)
			//migrate IDs
			if exist {
				members[i].ID = oldMembers[j].ID
			} else {
				members[i].ID = getNewID(members)
			}
		}
	}

	//only add new members
	for _, newMember := range members {

		_, exist := memberExist(newMember, oldMembers)
		if !exist {
			newMember.ID = getNewID(oldMembers)
			members = append(members, newMember)
		}
	}

	//sort members by id
	sort.Slice(members, func(i, j int) bool {
		return members[i].ID < members[j].ID
	})

	repConfig["members"] = members
	var reconfRes bson.M
	err = db.RunCommand(ctx, bson.D{{"replSetReconfig", repConfig}}).Decode(&reconfRes)
	if err != nil {
		return fmt.Errorf("replSetReconfig %w", err)
	}
	//if reconfRes["ok"].(int) != 0 { //TODO: handle code
	//
	//}

	bytes, err := bson.Marshal(reconfRes)
	fmt.Sprintf("%s,%v", bytes, err)
	return nil
}

func memberExist(member config.Member, members []config.Member) (int, bool) {
	for i, m := range members {
		if m.Host == member.Host {
			return i, true
		}
	}
	return 0, false
}

//
//func replicaSetExist(ctx context.Context, db *mongo.Database) {
//	collections, err := db.ListCollectionNames()
//
//}

//func hasReplicaConfig(ctx context.Context, db *mongo.Database) bool {
//
//}

func getNewID(members []config.Member) int {
	var membersCpy = make([]config.Member, len(members))
	copy(membersCpy, members)

	if len(members) == 0 {
		return 1
	}
	sort.Slice(membersCpy, func(i, j int) bool {
		return membersCpy[i].ID > membersCpy[j].ID
	})
	return membersCpy[0].ID + 1
}
