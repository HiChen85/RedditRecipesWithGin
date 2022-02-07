package tests

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func TestChannel(t *testing.T) {
	// 验证了一个功能,当从一个无论是有缓冲还是无缓冲的 没有数据写入的 channel 向外读数据时, 程序会被阻塞,类似死循环的效果.
	forerve := make(chan bool, 1)
	
	go func() {
		for i := 0; i < 30; i++ {
			log.Println(i)
			time.Sleep(time.Second)
		}
		forerve <- true
	}()
	
	fmt.Println(<-forerve)
}
