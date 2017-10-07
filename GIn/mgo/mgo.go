package main

import (
	"GIn/mongodb" //这里要改成mongodb.go的文件存放路径
	"fmt"
	"gopkg.in/mgo.v2/bson"

	"time"
)

type User struct {
	Id_       bson.ObjectId `bson:"_id"`
	Name      string        `bson:"name"`
	Age       int           `bson:"age"`
	JoinedAt  time.Time     `bson:"joned_at"`
	Interests []string      `bson:"interests"`
}

func main() {
	session := mongodb.CloneSession() //调用这个获得session
	defer session.Close()             //一定要记得释放

	c := session.DB("test1").C("people")
	err := c.Insert(&User{Id_: bson.NewObjectId(),
		Name:      "Jimmy Kuu",
		Age:       33,
		JoinedAt:  time.Now(),
		Interests: []string{"Develop", "Movie"}})

	if err != nil {
		panic(err)
	}

	var users []User
	err = c.Find(nil).Limit(5).All(&users)
	if err != nil {
		panic(err)
	}
	fmt.Println(users)

}
