package manager

import "github.com/WuKongIM/WuKongIM/internal/types"

type inFlightPqueue []*types.RetryMessage // 队列消息量大，这里直接使用具体类型，防止大量类型转换影响性能

func newInFlightPqueue(capacity int) inFlightPqueue {
	return make(inFlightPqueue, 0, capacity)
}

func (pq inFlightPqueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

func (pq *inFlightPqueue) Push(x *types.RetryMessage) {
	n := len(*pq)
	c := cap(*pq)
	if n+1 > c {
		npq := make(inFlightPqueue, n, c*2)
		copy(npq, *pq)
		*pq = npq
	}
	*pq = (*pq)[0 : n+1]
	x.Index = n
	(*pq)[n] = x
	pq.up(n)
}

func (pq *inFlightPqueue) Pop() *types.RetryMessage {
	n := len(*pq)
	c := cap(*pq)
	pq.Swap(0, n-1)
	pq.down(0, n-1)
	if n < (c/2) && c > 25 {
		npq := make(inFlightPqueue, n, c/2)
		copy(npq, *pq)
		*pq = npq
	}
	x := (*pq)[n-1]
	x.Index = -1
	*pq = (*pq)[0 : n-1]
	return x
}

func (pq *inFlightPqueue) Remove(i int) *types.RetryMessage {
	n := len(*pq)
	if n-1 != i {
		pq.Swap(i, n-1)
		pq.down(i, n-1)
		pq.up(i)
	}
	x := (*pq)[n-1]
	x.Index = -1
	*pq = (*pq)[0 : n-1]
	return x
}

func (pq *inFlightPqueue) PeekAndShift(max int64) (*types.RetryMessage, int64) {
	if len(*pq) == 0 {
		return nil, 0
	}

	x := (*pq)[0]
	if x.Pri > max {
		return nil, x.Pri - max
	}
	pq.Pop()

	return x, 0
}

func (pq *inFlightPqueue) up(j int) {
	for {
		i := (j - 1) / 2 // parent
		if i == j || (*pq)[j].Pri >= (*pq)[i].Pri {
			break
		}
		pq.Swap(i, j)
		j = i
	}
}

func (pq *inFlightPqueue) down(i, n int) {
	for {
		j1 := 2*i + 1
		if j1 >= n || j1 < 0 { // j1 < 0 after int overflow
			break
		}
		j := j1 // left child
		if j2 := j1 + 1; j2 < n && (*pq)[j1].Pri >= (*pq)[j2].Pri {
			j = j2 // = 2*i + 2  // right child
		}
		if (*pq)[j].Pri >= (*pq)[i].Pri {
			break
		}
		pq.Swap(i, j)
		i = j
	}
}
