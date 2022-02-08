package main

import (
	"fmt"
	"github.com/go-redis/redis/v8"
)

var (
	MONGO_USERNAME              = "admin"
	MONGO_PASSWORD              = "password"
	MONGO_HOST                  = "localhost"
	MONGO_PORT                  = "27017"
	MONGO_DATABASE              = "demo"
	MONGO_RECIPES_COLLECTION    = "recipes"
	MONGO_USER_COLLECTION       = "users"
	MONGO_RSS_PARSER_COLLECTION = "parseRecipes"
	MONGO_CONSUMER_COLLECTION   = "consumerRecipes"
)

var MONGO_URI = fmt.Sprintf("mongodb://%v:%v@%v:%v/%v?authSource=%v", MONGO_USERNAME, MONGO_PASSWORD, MONGO_HOST, MONGO_PORT, MONGO_DATABASE, MONGO_USERNAME)

// Redis connection config
var (
	REDIS_HOST     = "localhost"
	REDIS_PORT     = "6379"
	REDIS_PASSWORD = ""
	REDIS_DB       = 0
)

var RedisOptions = redis.Options{
	Addr:     REDIS_HOST + ":" + REDIS_PORT,
	Password: REDIS_PASSWORD,
	DB:       REDIS_DB,
}
