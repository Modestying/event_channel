package event_channel

import (
	"sync"

	"github.com/go-kratos/kratos/v2/log"
)

// MessageCenter 消息中心
type MessageCenter struct {
	topics sync.Map
}

// NewMessageCenter 创建消息中心
func NewMessageCenter() *MessageCenter {
	instance := &MessageCenter{}
	return instance
}

// Subscribe 订阅
func (m *MessageCenter) Subscribe(topicName string, client Consumer) error {
	if _, exist := m.topics.Load(topicName); !exist {
		// 只有第一个订阅者才会进入，其他goroutine则在第一个goroutine执行完之前等待（subscribe function的调用时同步的）
		log.Debug("创建新的主题: ", topicName)
		topic := NewTopic(topicName)
		topic.AddClient(client)
		m.topics.Store(topicName, topic)
	} else {
		// 其他goroutine等待第一个goroutine执行完之后就可以继续了，不同goroutine的处理顺序可能不同，但是只需要确保只创建了一个topic
		topicConsumers, _ := m.topics.Load(topicName)
		topicConsumers.(*Topic).AddClient(client)
	}
	log.Debug(client.Description(), "订阅 ", topicName, " 成功!")
	return nil
}

// Publish 发布
func (m *MessageCenter) Publish(topicName string, data any) {
	if topicConsumers, exist := m.topics.Load(topicName); exist {
		topicConsumers.(*Topic).Notify(data)
	}
}

// UnSubScribe 取消订阅
func (m *MessageCenter) UnSubScribe(topicName, channelName string) {
	if topicConsumers, exist := m.topics.Load(topicName); exist {
		topicConsumers.(*Topic).RemoveClient(channelName)
	}
}
