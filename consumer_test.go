package event_channel

import (
	"log"
	"reflect"
	"testing"
	"time"
)

func TestChannelConsumer_ConsumeCallBack(t *testing.T) {
	consumer := NewChannelConsumer("test")

	targetData := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	var testData []int
	go func() {
		for data := range consumer.GetTransport() {
			log.Println(" read: ", data)
			testData = append(testData, data.(int))
		}
		log.Println("通道关闭")
	}()

	for i := 0; i < 10; i++ {
		consumer.ConsumeCallBack(i)
		time.Sleep(time.Second * 1)
	}

	if !reflect.DeepEqual(targetData, testData) {
		t.Errorf("targetData: %v, testData: %v", targetData, testData)
	} else {
		t.Log("测试通过")
	}
}

func TestChannelConsumer_CloseTransport(t *testing.T) {
	consumer := NewChannelConsumer("testCloseTransport")

	targetData := []int{0, 1, 2, 3, 4}
	var testData []int

	time.AfterFunc(time.Second*5, func() {
		consumer.CloseTransport()
	})

	time.AfterFunc(time.Second*6, func() {
		consumer.CloseTransport()
	})

	go func() {
		for data := range consumer.GetTransport() {
			log.Println(" read: ", data)
			testData = append(testData, data.(int))
		}
		log.Println("通道关闭")
	}()

	for i := 0; i < 10; i++ {
		consumer.ConsumeCallBack(i)
		log.Println("发送: ", i)
		time.Sleep(time.Second * 1)
	}

	if !reflect.DeepEqual(targetData, testData) {
		t.Errorf("targetData: %v, testData: %v", targetData, testData)
	} else {
		t.Log("测试通过")
	}
}
