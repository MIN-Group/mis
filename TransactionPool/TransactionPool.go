package TransactionPool

import (
	"MIS-BC/MetaData"
	"sync"
)

//事务池
type TransactionPool struct {
	txsPool  map[int][][]byte //int 映射到 二维字节数组
	txsSize  int              //
	amout    int              //事务池中当前事务的数量
	capacity int              //事务池塘的容量
	writeNum int
	readNum  int
	lock     sync.Mutex //互斥锁
}

func (tp *TransactionPool) Init(txsSize int, capacity int) {
	var k int = 0
	if txsSize <= 0 || capacity <= 0 {
		tp.txsSize = 0
		tp.capacity = 0
		k = 1
	} else {
		tp.lock.Lock()
		tp.txsPool = make(map[int][][]byte)
		if tp.txsSize > tp.capacity {
			tp.txsSize = capacity
		} else {
			tp.txsSize = txsSize
		}
		tp.capacity = capacity
		if capacity%txsSize > 0 {
			k = capacity/txsSize + 1
		} else {
			k = capacity / txsSize
		}
		for i := 0; i < k; i++ {
			txs := make([][]byte, 0, tp.txsSize)
			tp.txsPool[i] = txs
		}
		tp.lock.Unlock()
	}
}

func (tp *TransactionPool) PushbackTransaction(header MetaData.TransactionHeader, transactionInterface MetaData.TransactionInterface) int {
	//事务池数量 大于容量 返回-1
	if tp.amout >= tp.capacity {
		return -1
	}
	for {
		//锁住 互斥
		tp.lock.Lock()
		//一个事务池对应一个二维字节数组 也就是一个int 对应一个事务组
		//writeNum指示当前正在写的事务组标号 如果当前事务组的大小 大于 规定的大小
		//则进入到下一个事务组
		if len(tp.txsPool[tp.writeNum]) >= tp.txsSize {
			tp.writeNum = (tp.writeNum + 1) % len(tp.txsPool)
		} else {
			tp.lock.Unlock()
			break
		}
	}
	tp.lock.Lock()
	tx := MetaData.EncodeTransaction(header, transactionInterface)
	tp.txsPool[tp.writeNum] = append(tp.txsPool[tp.writeNum], tx)
	tp.amout++
	tp.lock.Unlock()
	return 1
}

func (tp *TransactionPool) GetCurrentTxsList() (txs [][]byte) {
	tp.lock.Lock()
	txs = tp.txsPool[tp.readNum]
	tp.readNum = (tp.readNum + 1) % len(tp.txsPool)
	tp.lock.Unlock()
	return
}

func (tp *TransactionPool) GetCurrentTxsListDelete() (txs [][]byte) {
	tp.lock.Lock()
	txs = tp.txsPool[tp.readNum]
	tp.amout -= len(txs)
	tp.txsPool[tp.readNum] = make([][]byte, 0, tp.txsSize)
	tp.readNum = (tp.readNum + 1) % len(tp.txsPool)
	tp.lock.Unlock()
	return
}

func (tp *TransactionPool) IsFull() bool {
	if tp.amout >= tp.capacity {
		return true
	} else {
		return false
	}
}
