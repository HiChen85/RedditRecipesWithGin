package main

import "fmt"

const (
	PROTOCOL              = "amqp"
	RABBITMQ_DEFAULT_USER = "user"
	RABBITMQ_DEFAULT_PASS = "password"
	RABBITMQ_HOST         = "localhost"
	RABBITMQ_PORT         = "5672"
)

var RABBITMQ_URI = fmt.Sprintf("%v://%v:%v@%v:%v/", PROTOCOL, RABBITMQ_DEFAULT_USER, RABBITMQ_DEFAULT_PASS, RABBITMQ_HOST, RABBITMQ_PORT)

var RABBITMQ_QUEUE = "rss_urls"
