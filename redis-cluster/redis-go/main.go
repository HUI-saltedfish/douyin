package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type User struct {
	Id   int
	Name string
}

type House struct {
	Id    int
	Users []User
}

func (u User) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

func (u *User) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, u)
}

func (h House) MarshalBinary() ([]byte, error) {
	return json.Marshal(h)
}

func (h *House) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, h)
}

func main() {
	// var user1 *User
	// var user2 = new(User)
	// fmt.Println(reflect.TypeOf(user1), reflect.ValueOf(user1))
	// fmt.Println(reflect.TypeOf(user2), reflect.ValueOf(user2))

	// connect to redis by cluster
	var ctx = context.Background()
	redisClient := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{"localhost:6374", "localhost:6375", "localhost:6376"},
	})

	// ping
	pong, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	println(pong)

	// // set a hash and expire it
	// user := &User{Id: 1, Name: "John"}
	// val := redisClient.HSet(ctx, "user", user.Name, user)
	// if val.Err() != nil {
	// 	panic(val.Err())
	// }
	// println(val.Val())
	// redisClient.Expire(ctx, "user", 10*time.Second)

	// // get the hash
	// var user2 = new(User)
	// user2.Name = "John"
	// err = redisClient.HGet(ctx, "user", user2.Name).Scan(user2)
	// if err != nil {
	// 	panic(err)
	// }
	// // json.Unmarshal([]byte(valStr), &user2)
	// fmt.Printf("%v\n", user2)

	// set the house and expire it
	// house := &House{Id: 1, Users: []User{{Id: 1, Name: "John"}, {Id: 2, Name: "Mary"}}}
	// redisClient.HSet(ctx, "house", strconv.Itoa(house.Id), house)
	// redisClient.Expire(ctx, "house", 10*time.Second)

	// get the house
	var house2 = new(House)
	house2.Id = 1
	err = redisClient.HGet(ctx, "house", strconv.Itoa(house2.Id)).Scan(house2)
	fmt.Println(err)
	fmt.Printf("%v\n", house2)

	// // get the key expire time
	// ttl, _ := redisClient.TTL(ctx, "user").Result()
	// fmt.Printf("%v\n", ttl)

	// close the connection
	defer redisClient.Close()
}
