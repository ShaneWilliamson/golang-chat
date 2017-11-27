package model

import (
	"container/heap"
	"fmt"
	"testing"
	"time"
)

func TestPriorityQueueOrdersByLastUsedCorrectly(t *testing.T) {
	pq := GetPriorityQueueInstance()
	// 3 rooms, all last used differs
	rm1Time, _ := time.Parse(time.RFC3339, "2017-11-08T23:01:00Z")
	rm2Time, _ := time.Parse(time.RFC3339, "2017-11-08T23:02:00Z")
	rm3Time, _ := time.Parse(time.RFC3339, "2017-11-08T23:03:00Z")
	room1 := &ChatRoom{LastUsed: rm1Time}
	room2 := &ChatRoom{LastUsed: rm2Time}
	room3 := &ChatRoom{LastUsed: rm3Time}

	// Push them out of order
	heap.Push(pq, room1)
	heap.Push(pq, room3)
	heap.Push(pq, room2)

	// Assert they're popped in the right order
	pop1 := heap.Pop(pq).(*ChatRoom)
	pop2 := heap.Pop(pq).(*ChatRoom)
	pop3 := heap.Pop(pq).(*ChatRoom)
	if pop1.LastUsed != rm1Time || pop2.LastUsed != rm2Time || pop3.LastUsed != rm3Time {
		fmt.Println("Last used times were not ordered properly.")
		fmt.Printf("1: %s\n", pop1.LastUsed.String())
		fmt.Printf("2: %s\n", pop2.LastUsed.String())
		fmt.Printf("3: %s\n", pop3.LastUsed.String())
		t.Fail()
	}
}

func TestPriorityQueueOrdersByLastUsedCorrectlyAfterUpdate(t *testing.T) {
	pq := GetPriorityQueueInstance()
	// 3 rooms, all last used differs
	rm1Time, _ := time.Parse(time.RFC3339, "2017-11-08T23:01:00Z")
	rm2Time, _ := time.Parse(time.RFC3339, "2017-11-08T23:02:00Z")
	rm3Time, _ := time.Parse(time.RFC3339, "2017-11-08T23:03:00Z")
	updatedTime, _ := time.Parse(time.RFC3339, "2017-11-08T23:04:00Z")
	room1 := &ChatRoom{LastUsed: rm1Time}
	room2 := &ChatRoom{LastUsed: rm2Time}
	room3 := &ChatRoom{LastUsed: rm3Time}

	// Push them out of order
	heap.Push(pq, room1)
	heap.Push(pq, room3)
	heap.Push(pq, room2)
	// Update room1's LastUsed time to be the most recent
	pq.Update(room1, updatedTime)

	// Assert they're popped in the right order
	pop1 := heap.Pop(pq).(*ChatRoom)
	pop2 := heap.Pop(pq).(*ChatRoom)
	pop3 := heap.Pop(pq).(*ChatRoom)
	if pop1.LastUsed != rm2Time || pop2.LastUsed != rm3Time || pop3.LastUsed != updatedTime {
		fmt.Println("Last used times were not ordered properly.")
		fmt.Printf("1: %s\n", pop1.LastUsed.String())
		fmt.Printf("2: %s\n", pop2.LastUsed.String())
		fmt.Printf("3: %s\n", pop3.LastUsed.String())
		t.Fail()
	}
}

func TestDurationManipulation(t *testing.T) {
	expectedTime := "2017-11-15 23:00:00 +0000 UTC"
	baseTime, _ := time.Parse(time.RFC3339, "2017-11-08T23:00:00Z")
	sleepUntil := baseTime.Add(time.Hour * time.Duration(7*24))
	if sleepUntil.String() != expectedTime {
		fmt.Printf("sleepUntil time incorrect.\nExpected: %s\nActual: %s\n", expectedTime, sleepUntil.String())
		t.Fail()
	}
}
