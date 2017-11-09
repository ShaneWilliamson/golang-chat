package model

import (
	"container/heap"
	"sync"
	"time"
)

// An Item is something we manage in a priority queue.
// type Item struct {
// 	value    string // The value of the item; arbitrary.
// 	priority int    // The priority of the item in the queue.
// 	// The Index is needed by update and is maintained by the heap.Interface methods.
// 	Index int // The Index of the item in the heap.
// }

// A PriorityQueue implements heap.Interface and holds ChatRooms.
type PriorityQueue []*ChatRoom

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the lowest number (earliest date), which is highest priority so we use less than here.
	return pq[i].LastUsed.Unix() < pq[j].LastUsed.Unix()
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

// Push pushes a chat room into the priority queue.
func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*ChatRoom)
	item.Index = n
	*pq = append(*pq, item)
}

// Pop pops the highest priority chat room.
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.Index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// Update modifies the priority and value of an Item in the queue.
func (pq *PriorityQueue) Update(item *ChatRoom, lastUsed time.Time) {
	item.LastUsed = lastUsed
	heap.Fix(pq, item.Index)
}

var pqinstance *PriorityQueue
var pqonce sync.Once

// GetPriorityQueueInstance returns a singleton instance of the program configuration
func GetPriorityQueueInstance() *PriorityQueue {
	pqonce.Do(func() {
		pqinstance = &PriorityQueue{}
	})
	return pqinstance
}
