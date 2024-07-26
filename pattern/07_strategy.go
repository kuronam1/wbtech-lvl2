package main

import (
	"fmt"
	"slices"
	"sort"
)

type SortStrategy interface {
	Sort([]int)
}

type QuickSort struct{}

func (q *QuickSort) Sort(data []int) {
	if len(data) < 2 {
		return
	}
	left, right := 0, len(data)-1
	pivot := len(data) / 2
	data[pivot], data[right] = data[right], data[pivot]
	for i := range data {
		if data[i] < data[right] {
			data[i], data[left] = data[left], data[i]
			left++
		}
	}
	data[left], data[right] = data[right], data[left]
	q.Sort(data[:left])
	q.Sort(data[left+1:])
}

type ReverseSort struct{}

func (r *ReverseSort) Sort(data []int) {
	sort.Ints(data)
	slices.Reverse(data)
}

type Context struct {
	strategy SortStrategy
}

func (c *Context) SetStrategy(strategy SortStrategy) {
	c.strategy = strategy
}

func (c *Context) Sort(data []int) {
	c.strategy.Sort(data)
}
func main() {
	data := []int{3, 1, 4, 1, 5, 9, 2, 6, 5, 3, 5}

	context := &Context{}

	fmt.Println("Original:", data)

	context.SetStrategy(&QuickSort{})
	context.Sort(data)
	fmt.Println("Quick sort:", data)

	context.SetStrategy(&ReverseSort{})
	context.Sort(data)
	fmt.Println("Reverse sort:", data)
}
