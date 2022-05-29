package register

import (
	"encoding/json"
	"mongo-discovery/config"
	"mongo-discovery/etcd"
)

func GetMembers(reg *etcd.Registrar) ([]config.Member, error) {
	entries, err := reg.GetEntries()
	if err != nil {
		return nil, err
	}
	members := make([]config.Member, len(entries))

	for i, entry := range entries {
		var member config.Member
		err := json.Unmarshal([]byte(entry), &member)
		if err != nil {
			members = members[:len(members)-1]
			continue
		}
		members[i] = member
	}
	return members, nil
}
