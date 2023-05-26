package event_channel

import (
	"fmt"
	"reflect"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

func TestNewEventChannel(t *testing.T) {
	num := func() {}
	log.Debug(reflect.TypeOf(num).Kind() == reflect.Func)
}

var wait *sync.WaitGroup

func TestEventMessage_Subscribe(t *testing.T) {
	wait = &sync.WaitGroup{}
	eventCenter := NewMessageCenter()

	for i := 0; i < 5; i++ {
		consumer := NewChannelConsumer(strconv.Itoa(i))

		eventCenter.Subscribe("devinfo", consumer)

		go func(consumer Consumer, index int) {
			target := 0
			for data := range consumer.GetTransport() {
				log.Debug(index, " read: ", data)
				target++
				//发送失败
				if target > 2*index {
					consumer.CloseTransport()
					fmt.Println(index, "设备消息通知接口:消息发送结束--------------------------------")
					wait.Done()
				}
			}
		}(consumer, i+1)
	}

	notifyChan := make(chan bool)
	go func() {
	LOOP:
		for {
			select {
			case <-notifyChan:
				break LOOP
			default:
				eventCenter.Publish("devinfo", "test")
				time.Sleep(time.Second * 1)
			}
		}
		log.Info("end")
	}()

	wait.Add(5)
	wait.Wait()
	notifyChan <- true
}

func TestMessageCenter_Subscribe(t *testing.T) {
	consumer1 := NewChannelConsumer("1")
	consumer2 := NewChannelConsumer("2")
	center := NewMessageCenter()
	center.Subscribe("devinfo", consumer1)
	center.Subscribe("devinfo", consumer2)
	center.Subscribe("devinfo", consumer2)
	center.UnSubScribe("devinfo", "2")
}
