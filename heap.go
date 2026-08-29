package main

type Item struct {
	IP string
	Ms int64
}

// 大根堆
type MaxHeap []*Item

func (h MaxHeap) Len() int {
	return len(h)
}

// 关键：大的优先
func (h MaxHeap) Less(i, j int) bool {
	return h[i].Ms > h[j].Ms
}

func (h MaxHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *MaxHeap) Push(x any) {
	*h = append(*h, x.(*Item))
}

func (h *MaxHeap) Pop() any {
	old := *h
	n := len(old)

	x := old[n-1]
	*h = old[:n-1]

	return x
}
