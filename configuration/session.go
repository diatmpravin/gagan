package configuration

import (
	"encoding/json"
	"github.com/diatmpravin/gagan/models"
	"github.com/garyburd/redigo/redis"
	"log"
	"math/rand"
	"strconv"
	"time"
)

type Session struct {
	SessionId int `json:"sessionid"`
	AccessToken  string    `json:"accesstoken"`
	Timestamp    time.Time `json:"timestamp"`
	Organization models.Organization
	Space        models.Space
	Application  models.Application
}

func HandleError(err error) {
	if err != nil {
		panic(err)
	}
}

func RedisConnect() redis.Conn {
	c, err := redis.Dial("tcp", "127.0.0.1:6379")
	HandleError(err)
	return c
}

func random(min, max int) int {
	return rand.Intn(max-min) + min
}

// FIXME, need to rework, it's not standered one
func NewId() int {
	rand.Seed(time.Now().Unix())
	myrand := random(1, 99)
	return myrand
}

func CreateSession(config *Configuration) (session Session) {
	sessionId := NewId()

	session.SessionId, session.Timestamp = sessionId, time.Now()

	config.SessionId = sessionId

	c := RedisConnect()
	defer c.Close()

	b, err := json.Marshal(config)

	HandleError(err)

	// Save JSON blob to Redis
	reply, err := c.Do("SET", "user:"+strconv.Itoa(session.SessionId), b)

	HandleError(err)
	log.Println("GET ", reply)
	return
}

func DeleteSession(id int) {

	c := RedisConnect()
	defer c.Close()

	reply, err := c.Do("DEL", "user:"+strconv.Itoa(id))
	HandleError(err)

	if reply.(int64) != 1 {
		log.Println("No session removed")
	} else {
		log.Println("Session removed")
	}
}
