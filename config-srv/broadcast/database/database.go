package database

import (
	"container/list"
	"sync"
	"time"

	"XConf/config-srv/broadcast"
	"XConf/config-srv/dao"
	"XConf/proto/config"
	"github.com/micro/go-micro/util/log"
)

type Broker struct {
	sync.RWMutex

	currentID int
	watchers  list.List
}

func New() (broadcast.Broadcast, error) {
	var (
		b   Broker
		err error
	)
	b.currentID, err = dao.GetDao().GetNewestMessageID()
	if err != nil {
		return nil, err
	}

	// 用来监控messageid是否有更新，也就是判断是否有新发布的配置
	go b.scanAndNotify()

	return &b, nil
}

func (b *Broker) Send(namespace *config.ConfigResponse) error {
	return dao.GetDao().CreateReleaseMessage(
		namespace.GetAppName(),
		namespace.GetClusterName(),
		namespace.GetNamespaceName(),
		namespace.GetFormat(),
		namespace.GetValue())
}

func (b *Broker) Watch() broadcast.Watcher {
	w := &Watcher{
		exit:    make(chan interface{}),
		updates: make(chan *config.ConfigResponse, 2), // TODO 1 ?? 2 ?? or config
	}

	b.Lock()
	e := b.watchers.PushBack(w)
	b.Unlock()

	go func() {
		<-w.exit
		b.Lock()
		b.watchers.Remove(e)
		b.Unlock()
	}()

	return w
}

func (b *Broker) scanAndNotify() {
	for {
		time.Sleep(time.Second)

		// 获取最新的发布配置的messageid
		newestID, err := dao.GetDao().GetNewestMessageID()
		if err != nil {
			log.Error(err)
			continue
		}

		if newestID > b.currentID {
			// 获取从上次发布配置到现在新产生的发布信息
			msgs, err := dao.GetDao().GetReleaseMessage(b.currentID)
			if err != nil {
				log.Error(err)
				continue
			}

			for _, msg := range msgs {
				event := &config.ConfigResponse{
					AppName:       msg.AppName,
					ClusterName:   msg.ClusterName,
					NamespaceName: msg.NamespaceName,
					Format:        msg.Format,
					Value:         msg.Value,
				}

				watchers := make([]*Watcher, 0, b.watchers.Len())
				b.RLock()
				for e := b.watchers.Front(); e != nil; e = e.Next() {
					watchers = append(watchers, e.Value.(*Watcher))
				}

				b.RUnlock()

				for _, w := range watchers {
					select {
					case w.updates <- event:
					default:
					}
				}

				if int(msg.ID) > b.currentID {
					b.currentID = int(msg.ID)
				}
			}
		}
	}
}
