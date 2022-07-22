package pkg

import (
	"fmt"
	"sync"
)

type (
	logStatusArray   []*int8
	logIDStatusIndex map[string]int
)

var (
	RequestHistory *logIDCache // 请求历史，对请求做幂等处理
)

type logIDCache struct {
	partitions []*logIDstatusPartition
	// cap(partitions)
	cap int
}

func (c *logIDCache) getPartition(id string) *logIDstatusPartition {
	x := int(id[len(id)-1])
	pid := x % c.cap
	return c.partitions[pid]
}

func (c *logIDCache) Find(logID string) int64 {
	if logID == "" {
		return 0
	}
	p := c.getPartition(logID)
	ch := make(chan *int8)
	p.find(logID, ch)
	s := <-ch
	if s == nil {
		return 0
	}

	return int64(*s)
}

func (c *logIDCache) Add(logID string, status int64) {
	if logID == "" {
		return
	}
	p := c.getPartition(logID)
	p.add(logID, int8(status))
}

func NewLogIDCache(partitions, partitionsCap, partitionCap int) *logIDCache {
	var cache logIDCache
	cache.cap = partitions
	for i := cache.cap; i > 0; i-- {
		cache.partitions = append(cache.partitions, newLogIDStatusPartition(partitionsCap, partitionCap))
	}

	fmt.Printf("初始化LogIDCache完成")

	return &cache
}

const (
	LogIdStatusNone int64 = iota
	LogIdStatusPending
	LogIdStatusFail
	LogIdStatusSuccess
	LogIdStatusRollback
)

type logIDStatus struct {
	status     logStatusArray
	index      logIDStatusIndex
	queryChan  chan map[string]chan<- *int8
	updateChan chan map[string]int8
	done       chan struct{}
}

func (l *logIDStatus) info() {
	fmt.Printf("logIDs:%v \n", l.index)
}

func (l *logIDStatus) find_(id string) *int8 {
	i, ok := l.index[id]
	if !ok {
		return nil
	}

	s := l.status[i]
	return s
}

func (l *logIDStatus) find(id string, ch chan<- *int8) {
	l.queryChan <- map[string]chan<- *int8{id: ch}
}

func (l *logIDStatus) add(id string, status int8) {
	l.updateChan <- map[string]int8{id: status}
}

func (l *logIDStatus) add_(id string, status int8) {
	s := l.find_(id)
	if s != nil {
		if *s == status {
			return
		}

		*s = status
		return
	}

	i := len(l.status)
	l.status = append(l.status, &status)
	l.index[id] = i
}

func (l *logIDStatus) isFull() bool {
	return len(l.status) >= cap(l.status)
}

func (l *logIDStatus) run() {
	for {
		select {
		// 查询
		case q := <-l.queryChan:
			for k, v := range q {
				v <- l.find_(k)
				close(v)
				break
			}

		case u := <-l.updateChan:
			for k, v := range u {
				l.add_(k, v)
			}

		case <-l.done:
			return
		}
	}
}

func (l *logIDStatus) close() {
	select {
	case <-l.done:
	default:
		close(l.done)
	}
}

func newLogIDStatus(c int) *logIDStatus {
	var s logIDStatus
	s.done = make(chan struct{})
	s.status = make([]*int8, 0, c)
	s.index = make(logIDStatusIndex)
	s.queryChan = make(chan map[string]chan<- *int8, 100)
	s.updateChan = make(chan map[string]int8, 100)
	go s.run()
	return &s
}

type logIDstatusPartition struct {
	data []*logIDStatus
	// cap(data)
	cap int
	// cap(logIDStatus)
	limit      int
	queryChan  chan map[string]chan<- *int8
	updateChan chan map[string]int8
	done       chan struct{}
}

func (p *logIDstatusPartition) drop() {
	if p.cap == 1 {
		p.data[0].close()
		p.data = nil
		return
	}

	p.data[0].close()
	data := p.data
	p.data = data[1:]
}

func (p *logIDstatusPartition) expand() {
	partition := newLogIDStatus(p.limit)
	p.data = append(p.data, partition)
}

func (p *logIDstatusPartition) find_(id string) *int8 {
	var wg sync.WaitGroup
	wg.Add(len(p.data))
	var status *int8
	for _, p := range p.data {
		go func(p *logIDStatus) {
			ch := make(chan *int8)
			p.find(id, ch)
			s := <-ch
			if s != nil {
				status = s
			}
			wg.Done()
		}(p)
	}
	wg.Wait()
	return status
}

func (p *logIDstatusPartition) isFull() bool {
	return len(p.data) >= p.cap
}

func (p *logIDstatusPartition) add_(id string, status int8) {
	l := len(p.data)
	d := p.data[l-1]
	if d.isFull() {
		// 如果满了，打印
		// d.info()
		if p.isFull() {
			p.drop()
		}

		p.expand()
		p.add_(id, status)
		return
	}

	d.add(id, status)
}

func (p *logIDstatusPartition) run() {
	for {
		select {
		case q := <-p.queryChan:
			for k, v := range q {
				v <- p.find_(k)
				close(v)
				break
			}

		case u := <-p.updateChan:
			for k, v := range u {
				p.add_(k, v)
			}

		case <-p.done:
			return
		}
	}
}

func (p *logIDstatusPartition) find(id string, ch chan<- *int8) {
	p.queryChan <- map[string]chan<- *int8{id: ch}
}

func (p *logIDstatusPartition) add(id string, status int8) {
	p.updateChan <- map[string]int8{id: status}
}

func newLogIDStatusPartition(cap, subCap int) *logIDstatusPartition {
	var p logIDstatusPartition
	p.cap = cap
	p.limit = subCap
	p.done = make(chan struct{})
	p.queryChan = make(chan map[string]chan<- *int8)
	p.updateChan = make(chan map[string]int8)
	for i := cap; i > 0; i-- {
		p.data = append(p.data, newLogIDStatus(p.limit))
	}
	go p.run()
	return &p
}