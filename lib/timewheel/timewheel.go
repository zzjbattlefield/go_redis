package timewheel

import (
	"container/list"
	"log"
	"time"
)

type Task struct {
	job    func()
	delay  time.Duration
	key    string
	cicrly int
}

type location struct {
	etask *list.Element
	slot  int
}

type timeWheel struct {
	interval    time.Duration //每个轮槽的时间间隔
	slots       []*list.List  //轮槽的切片
	ticker      *time.Ticker
	currentSlot int //当前转动到的轮槽
	slotNum     int //轮槽的数量
	addTaskChan chan *Task
	rmTaskchan  chan string
	stopChan    chan struct{}
	timer       map[string]*location
}

func NewTimeWheel(interval time.Duration, slotsNum int) (wheel *timeWheel) {
	wheel = &timeWheel{
		interval:    interval,
		slots:       make([]*list.List, slotsNum),
		slotNum:     slotsNum,
		addTaskChan: make(chan *Task),
		rmTaskchan:  make(chan string),
		stopChan:    make(chan struct{}),
		timer:       make(map[string]*location),
	}
	for i := 0; i < slotsNum; i++ {
		wheel.slots[i] = list.New()
	}
	return
}

func (t *timeWheel) Start() {
	t.ticker = time.NewTicker(t.interval)
	go t.start()
}

func (t *timeWheel) Stop() {
	t.stopChan <- struct{}{}
}

func (t *timeWheel) RemoveJob(key string) {
	t.rmTaskchan <- key
}

func (t *timeWheel) AddJob(delay time.Duration, key string, job func()) {
	if delay < 0 {
		return
	}
	t.addTaskChan <- &Task{job: job, delay: delay, key: key}
}

func (t *timeWheel) start() {
	for {
		select {
		case <-t.ticker.C:
			t.tickHandler()
		case key := <-t.rmTaskchan:
			t.removeTask(key)
		case task := <-t.addTaskChan:
			t.addTask(task)
		case <-t.stopChan:
			t.ticker.Stop()
			return
		}
	}
}

func (t *timeWheel) removeTask(key string) {
	location, ok := t.timer[key]
	if !ok {
		return
	}
	l := t.slots[location.slot]
	l.Remove(location.etask)
	delete(t.timer, key)
}

// tickHandler 主要是处理轮槽下的链表
func (t *timeWheel) tickHandler() {

	l := t.slots[t.currentSlot]

	if t.currentSlot < t.slotNum-1 {
		t.currentSlot++
	} else {
		t.currentSlot = 0
	}
	go t.scanAndRunTask(l)
}

func (t *timeWheel) scanAndRunTask(l *list.List) {
	for e := l.Front(); e != nil; {
		task := e.Value.(*Task)
		if task.cicrly > 0 {
			task.cicrly--
			e = e.Next()
			continue
		}
		go func() {
			defer func() {
				if err := recover(); err != nil {
					log.Println(err)
				}
			}()
			task.job()
		}()
		next := e.Next()
		l.Remove(e)
		if task.key != "" {
			delete(t.timer, task.key)
		}
		e = next
	}
}

func (t *timeWheel) addTask(task *Task) {
	position, circle := t.getPositionAndCircle(task.delay)
	task.cicrly = circle
	e := t.slots[position].PushBack(task)
	loc := &location{
		etask: e,
		slot:  position,
	}
	if task.key != "" {
		//把这个key之前存在的task的删掉 相当于覆盖
		if _, ok := t.timer[task.key]; ok {
			t.removeTask(task.key)
		}
	}
	t.timer[task.key] = loc
}

func (t *timeWheel) getPositionAndCircle(d time.Duration) (position int, circle int) {
	delaySeconds := int(d.Seconds())
	intervalSeconds := int(t.interval.Seconds())
	circle = int(delaySeconds / intervalSeconds / t.slotNum)
	position = (delaySeconds/intervalSeconds + t.currentSlot) % t.slotNum
	return
}
