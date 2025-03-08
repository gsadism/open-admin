package snowflake

import (
	"sync"
	"time"
)

type Worker struct {
	sync.Mutex
	lastStamp  int64
	machineID  int64 //机器id,0~31
	serviceID  int64 //服务id,0~31
	sequenceID int64

	machineBits  int64
	serviceBits  int64
	sequenceBits int64

	maxMachineID  int64
	maxServiceID  int64
	maxSequenceID int64

	timeLeft    uint8
	machineLeft uint8
	serviceLeft uint8

	twepoch int64
}

func New(machineID, serviceID int64) *Worker {
	w := new(Worker)
	w.machineID = machineID
	w.serviceID = serviceID
	w.lastStamp = 0
	w.sequenceID = 0

	w.machineBits = int64(5)   // 机器ID位数
	w.serviceBits = int64(5)   // 服务ID位数
	w.sequenceBits = int64(12) // 序列ID位数

	w.maxMachineID = int64(-1) ^ (int64(-1) << w.machineBits)   // 最大机器ID
	w.maxServiceID = int64(-1) ^ (int64(-1) << w.serviceBits)   // 最大服务ID
	w.maxSequenceID = int64(-1) ^ (int64(-1) << w.sequenceBits) // 最大序列ID

	w.timeLeft = uint8(22)    // 时间ID向左移位的量
	w.machineLeft = uint8(17) // 机器ID向左移位的量
	w.serviceLeft = uint8(12) // 服务ID向左移位的量

	w.twepoch = int64(1667972427000) //初始毫秒,时间是: Wed Nov  9 13:40:27 CST 2022

	return w
}

func (w *Worker) GetID() int64 {
	//多线程互斥
	w.Lock()
	defer w.Unlock()

	mill := time.Now().UnixMilli()

	if mill == w.lastStamp {
		w.sequenceID = (w.sequenceID + 1) & w.maxSequenceID
		//当一个毫秒内分配的id数>4096个时，只能等待到下一毫秒去分配。
		if w.sequenceID == 0 {
			for mill > w.lastStamp {
				mill = time.Now().UnixMilli()
			}
		}
	} else {
		w.sequenceID = 0
	}

	w.lastStamp = mill

	id := (w.lastStamp-w.twepoch)<<w.timeLeft | w.machineID<<w.machineLeft | w.serviceID<<w.serviceLeft | w.sequenceID
	return id
}
