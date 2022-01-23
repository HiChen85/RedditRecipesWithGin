package utils

import (
	"fmt"
)

var (
	MONGO_USERNAME   = "admin"
	MONGO_PASSWORD   = "password"
	MONGO_HOST       = "localhost"
	MONGO_PORT       = "27017"
	MONGO_DATABASE   = "demo"
	MONGO_COLLECTION = "recipes"
)

var MONGO_URI = fmt.Sprintf("mongodb://%v:%v@%v:%v/%v?authSource=%v", MONGO_USERNAME, MONGO_PASSWORD, MONGO_HOST, MONGO_PORT, MONGO_DATABASE, MONGO_USERNAME)
