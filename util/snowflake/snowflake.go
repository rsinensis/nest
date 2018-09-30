package snowflake

import (
	"fmt"
	"sync"
	"time"
)

/*
核心代码为其Id这个类实现，其原理结构如下，分别用一个0表示一位，用—分割开部分的作用：
0 --- 0000000000 0000000000 0000000000 0000000000 0 --- 00000 --- 00000 --- 000000000000
在上面的字符串中，第一位为未使用（实际上也可作为long的符号位），接下来的41位为毫秒级时间，然后5位datacenter标识位，5位机器ID（并不算标识符，实际是为线程标识），然后12位该毫秒内的当前毫秒内的计数，加起来刚好64位，为一个Long型。
这样的好处是，整体上按照时间自增排序，并且整个分布式系统内不会产生ID碰撞（由datacenter和机器ID作区分），并且效率较高，经测试，snowflake每秒能够产生26万ID左右，完全满足需要。
*/

const (
	twepoch            = int64(1483200000000) // 2017/01/01 00:00:00
	workerIdBits       = uint(5)
	datacenterIdBits   = uint(5)
	maxWorkerId        = -1 ^ (-1 << workerIdBits)
	maxDatacenterId    = -1 ^ (-1 << datacenterIdBits)
	sequenceBits       = uint(12)
	workerIdShift      = sequenceBits
	datacenterIdShift  = sequenceBits + workerIdBits
	timestampLeftShift = sequenceBits + workerIdBits + datacenterIdBits
	sequenceMask       = -1 ^ (-1 << sequenceBits)
	maxNextIdsNum      = 1000 //小于4096
)

type Id struct {
	sequence      int64
	lastTimestamp int64
	workerId      int64
	twepoch       int64
	datacenterId  int64
	mutex         sync.Mutex
}

func GetIdTwepoch() int64 {
	return twepoch
}

// NewId new a snowflake id generator object.
func NewId(datacenterId, workerId int64, twepoch int64) (*Id, error) {
	id := &Id{}

	if workerId > maxWorkerId || workerId < 0 {
		return nil, fmt.Errorf("worker Id: %d error", workerId)
	}

	if datacenterId > maxDatacenterId || datacenterId < 0 {
		return nil, fmt.Errorf("datacenter Id: %d error", datacenterId)
	}

	id.workerId = workerId
	id.datacenterId = datacenterId
	id.lastTimestamp = -1
	id.sequence = 0
	id.twepoch = twepoch
	id.mutex = sync.Mutex{}

	return id, nil
}

// timeGen generate a unix millisecond.
func timeGen() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

// tilNextMillis spin wait till next millisecond.
func tilNextMillis(lastTimestamp int64) int64 {
	timestamp := timeGen()

	for timestamp <= lastTimestamp {
		timestamp = timeGen()
	}

	return timestamp
}

func (id *Id) unsafeId() (int64, error) {
	timestamp := timeGen() //当前毫秒数

	if timestamp < id.lastTimestamp { // 时间戳有问题
		return 0, fmt.Errorf("Clock moved backwards.  Refusing to generate id for %d milliseconds", id.lastTimestamp-timestamp)
	}

	if id.lastTimestamp == timestamp { //如果时间戳一致
		id.sequence = (id.sequence + 1) & sequenceMask //更新步进
		if id.sequence == 0 {                          //如果超过毫秒内步进，则更新时间戳
			timestamp = tilNextMillis(id.lastTimestamp)
		}
	} else { //如果时间戳不一致，重置步进
		id.sequence = 0
	}

	id.lastTimestamp = timestamp //更新时间戳

	return ((timestamp - id.twepoch) << timestampLeftShift) | (id.datacenterId << datacenterIdShift) | (id.workerId << workerIdShift) | id.sequence, nil //生成id
}

// NextId get a snowflake id.
func (id *Id) NextId() (int64, error) {
	id.mutex.Lock() //加锁
	defer id.mutex.Unlock()

	return id.unsafeId()
}

// NextIds get snowflake ids.
func (id *Id) NextIds(num int) ([]int64, error) {
	if num > maxNextIdsNum || num < 0 {
		return nil, fmt.Errorf("NextIds num: %d error", num)
	}

	ids := make([]int64, num)

	id.mutex.Lock()
	defer id.mutex.Unlock()

	for i := 0; i < num; i++ {

		idNum, err := id.unsafeId()

		if err != nil {
			return nil, fmt.Errorf("NextIds num: %d error: %v", num, err)
		}

		ids[i] = idNum
	}
	return ids, nil
}
