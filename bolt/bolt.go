package bolt

import (
	"github.com/boltdb/bolt"
	"github.com/3stadt/GoTBot/structs"
	"encoding/json"
	"github.com/imdario/mergo"
	"github.com/3stadt/GoTBot/context"
	"fmt"
)

var db *bolt.DB

func CreateOrUpdateUser(updateUser structs.User) error {
	baseUser := GetUser(updateUser.Name)
	if baseUser != nil {
		if err := mergo.MergeWithOverwrite(baseUser, &updateUser); err != nil {
			panic(err)
		}
	} else {
		baseUser = &updateUser
	}
	open()
	dberr := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(context.UserBucketName))
		err := b.Put([]byte(baseUser.Name), marshalUser(*baseUser))
		return err
	})
	db.Close()
	return dberr
}

func GetUser(username string) *structs.User {
	open()
	var user *structs.User
	var v []byte
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(context.UserBucketName))
		if b == nil {
			return nil
		}
		v = b.Get([]byte(username))
		if len(v) == 0 {
			return nil
		}
		var err error
		user, err = unmarshalUser(v)
		return err
	})
	if err != nil {
		fmt.Println(string(v))
		panic(err)
	}
	db.Close()
	return user
}

func marshalUser(user structs.User) []byte {
	jUser, err := json.Marshal(user)
	if err != nil {
		panic(err)
	}
	return jUser
}

func unmarshalUser(bytes []byte) (*structs.User, error) {
	user := &structs.User{}
	if err := json.Unmarshal(bytes, user); err != nil {
		return nil, err
	}
	return user, nil
	/*var objMap map[string]*json.RawMessage
	var user structs.User
	if len(bytes) == 0 {
		return nil, nil
	}
	if err:= json.Unmarshal(bytes, &objMap); err != nil {
		return err, nil
	}

	json.Unmarshal(*objMap["Name"], user)
	json.Unmarshal(*objMap["LastJoin"], user)
	json.Unmarshal(*objMap["LastPart"], user)
	json.Unmarshal(*objMap["LastActive"], user)
	json.Unmarshal(*objMap["FirstSeen"], user)
	return nil, &user*/
}

func open() {
	var err error
	db, err = bolt.Open("gotbot.db", 0600, nil)
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(context.UserBucketName))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}
