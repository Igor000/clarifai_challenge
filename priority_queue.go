// This example demonstrates a priority queue built using the heap interface.
package main

import (
	"container/heap"
	"fmt"
)

// An Item is something we manage in a priority queue.
type Item struct {
	value    string  // The value of the item; arbitrary.
	priority float64 // The priority of the item in the queue.
	// The index is needed by update and is maintained by the heap.Interface methods.
	index int // The index of the item in the heap.
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Print() {
	// fmt.Println("The whole pq ",  pq)

	fmt.Println("Priority queue: value,  priority,  index  ")
	for i := len(pq) - 1; i >= 0; i-- {
		fmt.Println("Priority queue ", pq[i].value, pq[i].priority, pq[i].index)
		// fmt.Println("Priority queue ", i, pq[i].value, pq[i].priority, pq[i].index)
	}
}

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	// return pq[i].priority > pq[j].priority
	// We want Pop to give us the lowest, not highest, priority so we use greater than here.
	return pq[i].priority < pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

/******
func (pq *PriorityQueue) Push(x interface{}, q_size int) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)

        n = len(*pq)
        if n > q_size {
	   old := *pq
           item := old[n-1]
           item.index = -1 // for safety
           _ = item
           *pq = old[0 : n-1]
        }
}
******/

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// update modifies the priority and value of an Item in the queue.
func (pq *PriorityQueue) update(item *Item, value string, priority float64) {
	item.value = value
	item.priority = priority
	heap.Fix(pq, item.index)
}

// This example creates a PriorityQueue with some items, adds and manipulates an item,
// and then removes the items in priority order.

// func main() {
func test() {

	q_size := 3
	// Some items and their priorities.
	items := map[string]float64{
		"banana": 3.3, "apple": 2.2, "pear": 4.4, "potato": 1.7, "another_veg": 2.1,
	}

	// Create a priority queue, put the items in it, and
	// establish the priority queue (heap) invariants.
	pq := make(PriorityQueue, len(items))
	i := 0
	for value, priority := range items {
		pq[i] = &Item{
			value:    value,
			priority: priority,
			index:    i,
		}
		i++
	}
	heap.Init(&pq)

	fmt.Println("After init")
	pq.Print()

	// downsize
	for pq.Len() > q_size {
		_ = heap.Pop(&pq).(*Item)
	}

	fmt.Println("After downsize")
	pq.Print()

	// Insert a new item and then modify its priority.
	item := &Item{
		value: "orange",
		// priority: 1.1,
		priority: 5.5,
	}
	heap.Push(&pq, item)
	// pq.update(item, item.value, 5.5)

	/////////
	/* */
	if pq.Len() > q_size {
		_ = heap.Pop(&pq).(*Item)
	}
	/* */

	item = &Item{
		value:    "orange7",
		priority: 7.1,
	}
	heap.Push(&pq, item)

	/////////
	/* */
	if pq.Len() > q_size {
		_ = heap.Pop(&pq).(*Item)
	}
	/* */

	pq.Print()

	// Take the items out; they arrive in decreasing priority order.
	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*Item)
		// fmt.Printf("%.2d:%s ", item.priority, item.value)
		fmt.Printf("%.2f:%s \n", item.priority, item.value)
	}

	pq.Print()
}
