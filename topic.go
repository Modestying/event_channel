package event_channel

import (
	"sync"

	"github.com/go-kratos/kratos/v2/log"
)

// Topic 主题
type Topic struct {
	Name    string
	clients sync.Map
	//clients map[string]Consumer // k=clientName v=client
	sync.Mutex
}

// NewTopic 创建主题
func NewTopic(name string) *Topic {
	return &Topic{
		Name: name,
	}
}

// AddClient 添加客户端
func (t *Topic) AddClient(consumer Consumer) {
	log.Debug("添加客户端,Topcic:", t.Name, " 客户端ID:", consumer.Description())

	if _, exist := t.clients.Load(consumer.Description()); exist {
		log.Debug("该客户端ID已存在,进行替换........")
	}
	t.clients.Store(consumer.Description(), consumer)
}

// RemoveClient 移除客户端
func (t *Topic) RemoveClient(clientName string) {
	log.Debug("删除客户端,Topcic:", t.Name, " 客户端ID:", clientName)
	if consumer, exist := t.clients.Load(clientName); exist {
		log.Debug("该客户端ID存在,进行删除........")
		consumer.(Consumer).CloseTransport()
		t.clients.Delete(clientName)
	}
}

// Notify 通知
func (t *Topic) Notify(message any) {
	t.clients.Range(func(key, value any) bool {
		value.(Consumer).ConsumeCallBack(message)
		return true
	})
}
