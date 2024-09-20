package bytes

import (
	"sync"
)

// Buffer buffer.
type Buffer struct {
	buf  []byte
	next *Buffer // next free buffer
}

// Bytes bytes.
func (b *Buffer) Bytes() []byte {
	return b.buf
}

// Pool is a buffer pool.
type Pool struct {
	lock sync.Mutex
	free *Buffer
	max  int
	num  int
	size int
}

// NewPool new a memory buffer pool struct.
func NewPool(num, size int) (p *Pool) {
	p = new(Pool)
	p.init(num, size)
	return
}

// Init init the memory buffer.
func (p *Pool) Init(num, size int) {
	p.init(num, size)
}

// init init the memory buffer.
func (p *Pool) init(num, size int) {
	p.num = num
	p.size = size
	p.max = num * size
	p.grow()
}

// grow grow the memory buffer size, and update free pointer.
func (p *Pool) grow() {
	var (
		i   int
		b   *Buffer
		bs  []Buffer
		buf []byte
	)
	buf = make([]byte, p.max)  // 内存池的总大小，一次创建好减少对象的创建
	bs = make([]Buffer, p.num) // 内存池的数量
	p.free = &bs[0]            // pool 记录第一个可用内存的指针
	// 初始化 bs 中的全部 Buffer
	b = p.free
	for i = 1; i < p.num; i++ {
		b.buf = buf[(i-1)*p.size : i*p.size]
		b.next = &bs[i]
		b = b.next
	}
	// 将最后一个 Buffer 的 next 指针设为 nil
	b.buf = buf[(i-1)*p.size : i*p.size]
	b.next = nil
}

// Get get a free memory buffer.
func (p *Pool) Get() (b *Buffer) {
	p.lock.Lock()
	if b = p.free; b == nil {
		p.grow()
		b = p.free
	}
	p.free = b.next
	p.lock.Unlock()
	return
}

// Put put back a memory buffer to free.
func (p *Pool) Put(b *Buffer) {
	p.lock.Lock()
	b.next = p.free
	p.free = b
	p.lock.Unlock()
}
