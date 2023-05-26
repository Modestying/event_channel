package event_channel

import (
	"sync"
	"sync/atomic"

	"github.com/fatih/color"
	"github.com/go-kratos/kratos/v2/log"
)

// Consumer 消费者接口
type Consumer interface {
	ConsumeCallBack(data any)
	Description() string
	GetTransport() chan any
	CloseTransport()
}

// channelConsumer 通道消费者
type channelConsumer struct {
	handler     chan any
	description string //通道标识
	close       *sync.Once
	done        uint32
	sync.Mutex
}

// NewChannelConsumer 创建通道消费者
func NewChannelConsumer(topic string) *channelConsumer {
	instance := &channelConsumer{
		handler:     make(chan any, 50),
		description: topic,
		close:       &sync.Once{},
		done:        0,
	}
	return instance
}

// ConsumeCallBack 消费回调
func (f *channelConsumer) ConsumeCallBack(data any) {
	//fmt.Println(f.description, ": ", data)
	defer func() {
		if err := recover(); err != nil {
			log.Errorf(color.RedString("通道写入数据失败,错误信息：%s", err, " 用户标识:%s", f.description))
		}
	}()
	f.Mutex.Lock()
	if atomic.LoadUint32(&f.done) == 0 {
		f.handler <- data
	} else {
		log.Error("客户端已关闭，无法写入数据: ", f.description)
	}
	f.Mutex.Unlock()
}

// Description 获取通道标识
func (f *channelConsumer) Description() string {
	return f.description
}

// GetTransport 获取通道
func (f *channelConsumer) GetTransport() chan any {
	return f.handler
}

// CloseTransport 关闭通道
func (f *channelConsumer) CloseTransport() {
	f.close.Do(func() {
		atomic.StoreUint32(&f.done, 1)
		close(f.handler)
		log.Info("关闭通道: ", f.description)
	})
}
