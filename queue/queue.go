// Copyright 2022 nnsgmsone
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package queue

import (
	"sync/atomic"
	"unsafe"
)

// New constructs a fixed size lockfree queue based on array
// 	size must be a multiple of 2
func New(size int) *Queue {
	return &Queue{
		size:  uint32(size),
		items: make([]item, size),
	}
}

// Push a value to the queue
func (q *Queue) Push(val any) {
	ht := atomic.LoadUint64(&q.headTail)
	head, tail := unpack(ht)
	if tail+q.size == head { // full
		return
	}
	it := &q.items[head&uint32(len(q.items)-1)]
	if typ := atomic.LoadPointer(&it.typ); typ != nil {
		return
	}
	*(*any)(unsafe.Pointer(it)) = val
	atomic.AddUint64(&q.headTail, 1<<32)
}

// Pop a value from the queue
func (q *Queue) Pop() any {
	var it *item

	for {
		ht := atomic.LoadUint64(&q.headTail)
		head, tail := unpack(ht)
		if tail == head { // empty
			return nil
		}
		nht := pack(head, tail+1) // new ptr
		if atomic.CompareAndSwapUint64(&q.headTail, ht, nht) {
			it = &q.items[tail&uint32(len(q.items)-1)]
			break
		}
	}
	val := *(*any)(unsafe.Pointer(it))
	it.val = nil
	atomic.StorePointer(&it.typ, nil)
	return val
}

func unpack(ht uint64) (head, tail uint32) {
	return uint32(ht>>32) & mask, uint32(ht & mask)
}

func pack(head, tail uint32) uint64 {
	return (uint64(head) << 32) |
		uint64(tail&mask)
}
