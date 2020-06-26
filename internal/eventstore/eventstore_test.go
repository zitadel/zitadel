package eventstore

import (
	"log"
	"sync"
	"testing"
)

type li struct {
	lock      sync.Mutex
	list      map[string][]*o
	listCount int
}

func (l *li) remove(ob *o) {
	l.lock.Lock()
	defer l.lock.Unlock()
	for _, name := range ob.names {
		os, ok := l.list[name]
		if !ok {
			continue
		}
		for i, obb := range os {
			if obb.idx != ob.idx {
				continue
			}

			os[i] = os[len(os)-1]
			os = os[:len(os)-1]
		}
		l.list[name] = os
	}
}

func (l *li) adds(names ...string) *o {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.listCount++
	ob := &o{idx: l.listCount, names: names}
	for _, name := range names {
		_, ok := l.list[name]
		if !ok {
			l.list[name] = make([]*o, 0)
		}

		l.list[name] = append(l.list[name], ob)
	}
	return ob
}

type o struct {
	idx   int
	names []string
}

func TestAdd(t *testing.T) {
	l := li{list: make(map[string][]*o)}
	l.adds("1", "2")
	l.adds("1", "2")
	l.adds("1", "2")
	log.Println(l)
}

func TestRemove(t *testing.T) {
	l := li{list: make(map[string][]*o)}
	l.adds("1", "2")
	l.adds("1", "2")
	ob := l.adds("1", "2")
	l.remove(ob)
	log.Println(l)
}
