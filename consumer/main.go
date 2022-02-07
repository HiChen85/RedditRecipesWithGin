package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/HiChen85/RedditRecipesWithGin/utils"
	"github.com/HiChen85/RedditRecipesWithGin/utils/rabbitmq"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
)

type Request struct {
	URL string `json:"url"`
}

var mongoClient *mongo.Client
var ctx context.Context

func init() {
	ctx = context.Background()
	tempClient, err := mongo.Connect(ctx, options.Client().ApplyURI(utils.MONGO_URI))
	if err != nil {
		log.Fatalln(err)
	}
	mongoClient = tempClient
	if err = mongoClient.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatalln(err)
	}
	log.Println("MongoDB has been connected")
}

func main() {
	
	// 一个消费者程序, 在消费完消息之后应自动关闭与消息队列的链接
	amqpConn, err := amqp.Dial(rabbitmq.RABBITMQ_URI)
	if err != nil {
		log.Fatalln(err)
	}
	defer amqpConn.Close()
	
	// 创建一个获取消息的通道
	channelAmqp, _ := amqpConn.Channel()
	defer channelAmqp.Close()
	
	forever := make(chan bool)
	// 获取消息的过程
	msg, err := channelAmqp.Consume(rabbitmq.RABBITMQ_QUEUE, "", true, false, false, false, nil)
	
	// 另起一个 go routine 来让消费者轮询消息队列
	go func() {
		for deliver := range msg {
			//log.Println("Received messages :", string(deliver.Body))
			var req Request
			json.Unmarshal(deliver.Body, &req)
			// 现在消费者程序需要根据这个 URL 去获取数据并存储
			entries, _ := utils.GetDataFromReddit(req.URL)
			recipeCollection := mongoClient.Database(utils.MONGO_DATABASE).Collection(utils.MONGO_CONSUMER_COLLECTION)
			// 头两个数据不要
			for i := range entries[2:] {
				_, err := recipeCollection.InsertOne(ctx, bson.M{
					"title":     entries[i].Title,
					"link":      entries[i].Link.Href,
					"thumbnail": entries[i].Thumbnail.URL,
				})
				if err != nil {
					log.Println(err)
					return
				}
			}
			fmt.Println(entries)
		}
	}()
	
	// 从没有数据的 channel 中读数据, 程序会被阻塞住, 直到有消息写入到 channel 中.
	// 本程序由于要一致轮询建立的消息队列, 那么就需要保持主程序阻塞而 goroutine 去读取消息体中的数据
	
	log.Printf(" [*] Waiting for messages... To exit, press CTRL+C\n")
	<-forever
}
