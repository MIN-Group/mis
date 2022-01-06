package Node

import (
	"MIS-BC/MetaData"
	"log"
)

type Ticket struct {
	counter map[string]int
}

func (node *Node) VotingStatistics(item MetaData.BlockGroup) (MetaData.BlockGroup, bool) {
	var THRESHOLD = node.config.VotedNum * 2 / 3
	var THRESHOLD2 = node.config.VotedNum / 2

	counter_yes := make([]Ticket, node.config.WorkerNum)
	counter_no := make([]Ticket, node.config.WorkerNum)
	result := make([]Ticket, node.config.WorkerNum)

	CounterInit(counter_yes)
	CounterInit(counter_no)
	CounterInit(result)

	for _, x := range item.VoteTickets {
		if x.BlockHashes == nil || x.VoteResult == nil {
			continue
		}
		if len(x.BlockHashes) > len(x.VoteResult) {
			log.Println(item)
			log.Println("VotingStatistics : len(x.BlockHashes) > len(x.VoteResult):")
			return item, false
		}
		for i := 0; i < len(x.BlockHashes); i++ {
			_, is_exist := counter_yes[i].counter[x.BlockHashes[i]]
			if !is_exist {
				counter_yes[i].counter[x.BlockHashes[i]] = 0
			}

			_, is_exist = counter_no[i].counter[x.BlockHashes[i]]
			if !is_exist {
				counter_no[i].counter[x.BlockHashes[i]] = 0
			}
			if x.VoteResult[i] == 1 {
				counter_yes[i].counter[x.BlockHashes[i]] += 1
			} else if x.VoteResult[i] == -1 {
				counter_no[i].counter[x.BlockHashes[i]] += 1
			}
		}
	}

	check := 0
	if node.round < 2 {
		for i, x := range counter_yes {
			for k, v := range x.counter {
				if v > THRESHOLD {
					result[i].counter[k] = 1
					check += 1
					break
				}
				if counter_no[i].counter[k] > THRESHOLD {
					result[i].counter[k] = -1
					check += 1
					break
				}
			}
		}
		if check != node.config.WorkerNum {
			return item, false
		}
	} else {
		for i, x := range counter_yes {
			for k, v := range x.counter {
				if v > THRESHOLD2 {
					result[i].counter[k] = 1
					check += 1
					break
				}
				if counter_no[i].counter[k] > THRESHOLD2 {
					result[i].counter[k] = -1
					check += 1
					break
				}
			}
		}
		if check == 0 {
			return item, false
		}
	}

	item.BlockHashes = make([]string, node.config.WorkerNum)
	item.VoteResult = make([]int, node.config.WorkerNum)
	for i, x := range result {
		for k, v := range x.counter {
			item.BlockHashes[i] = k
			item.VoteResult[i] = v
			/*			item.BlockHashes = append(item.BlockHashes, k)
						item.VoteResult = append(item.VoteResult, v)*/
		}
	}
	return item, true
}

func CounterInit(counter []Ticket) {
	for i := 0; i < len(counter); i++ {
		counter[i].counter = make(map[string]int)
	}
}
