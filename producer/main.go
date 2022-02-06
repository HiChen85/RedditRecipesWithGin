package main

import (
	"encoding/json"
	"github.com/HiChen85/RedditRecipesWithGin/utils/rabbitmq"
	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
	"log"
	"net/http"
)

/*
	producer 是将 rss_parser 拆分开来, 专门负责处理用户输入的各种类型的.rss url
	因为 reddit 的所有帖子的数据格式几乎一致,仅有 url 不同,所以可以用同一个生产者函数
	处理多个数据源
*/

var channelAmqp *amqp.Channel

func init() {
	amqpConn, err := amqp.Dial(rabbitmq.RABBITMQ_URI)
	if err != nil {
		log.Fatalln(err)
	}
	channelAmqp, _ = amqpConn.Channel()
}

type Request struct {
	URL string `json:"url"`
}

func ParseHandler(c *gin.Context) {
	var request Request
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	// 将接收过来绑定成 go 结构体对象的数据重新转成 json, 发送到消息队列中去
	data, _ := json.Marshal(request)
	/*
		解释一下 rabbitmq 的 Publish 的几个参数:
	
		exchange: 翻译成交换机也行,路由器也行, 是消息队列负责转发消息的部分,消息被生产者发送到消息队列时, 第一站就是 exchange.
			exchange 负责将收到的消息, 根据 routing key (binding key), 通过对应key 指定的 binding 发送的 binding 绑定的队列中去.
			rabbitmq 有几种不同的 exchange, 其中一种是默认的 exchange. 通常在 api 实现中,用空字符串表示.
		key: routing key, 指明了 exchange 要发往的队列绑定的 binding, 一般, 创建的队列会默认创建一个同名的 routing key (binding key) 然后
			绑定到默认的 exchange 上去. 当调用发布函数时,不指定 exchange, 那么rabbitmq就会根据指定的 key 把消息发送到同名的队列中去.
	*/
	err := channelAmqp.Publish("", rabbitmq.RABBITMQ_QUEUE, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        data,
	})
	
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while publishing to RabbitMQ"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
	
}

func main() {
	engine := gin.Default()
	engine.POST("/parse", ParseHandler)
	engine.Run(":8001")
}
