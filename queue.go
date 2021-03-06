/*
Package queue provides a fast, ring-buffer queue based on the version suggested by Dariusz Górecki.
Using this instead of other, simpler, queue implementations (slice+append or linked list) provides
substantial memory and time benefits, and fewer GC pauses.

The queue implemented here is as fast as it is for two additional reasons: it is *not* thread-safe, and it
intentionally does not follow go best-practices regarding errors - if you make a mistake with this
queue (such as trying to remove an element from an empty queue) then who knows what will happen.
*/
package queue

const minQueueLen = 16

// Queue represents a single instance of the queue data structure.
type Queue struct {
	buf               []interface{}
	head, tail, count int
}

// New constructs and returns a new Queue.
func New() *Queue {
	return &Queue{buf: make([]interface{}, minQueueLen)}
}

// Length returns the number of elements currently stored in the queue.
func (q *Queue) Length() int {
	return q.count
}

func (q *Queue) resize() {
	newBuf := make([]interface{}, q.count*2)

	if q.tail > q.head {
		copy(newBuf, q.buf[q.head:q.tail])
	} else {
		copy(newBuf, q.buf[q.head:len(q.buf)])
		copy(newBuf[len(q.buf)-q.head:], q.buf[:q.tail])
	}

	q.head = 0
	q.tail = q.count
	q.buf = newBuf
}

// Add puts an element on the end of the queue.
func (q *Queue) Add(elem interface{}) {
	if q.count == len(q.buf) {
		q.resize()
	}

	q.buf[q.tail] = elem
	q.tail = (q.tail + 1) % len(q.buf)
	q.count++
}

// Peek returns the element at the head of the queue. If the queue is empty (Length == 0),
// Peek does not panic, it simply returns garbage.
func (q *Queue) Peek() interface{} {
	return q.buf[q.head]
}

// Get returns the element at index i in the queue. If the index is invalid, the
// call will panic.
func (q *Queue) Get(i int) interface{} {
	if i >= q.Length() || i < 0 {
		panic("index out of range")
	}
	modi := (q.head + i) % len(q.buf)
	return q.buf[modi]
}

// Remove removes the element from the front of the queue. If you actually want the element,
// call Peek first. If the queue is empty (Length == 0), Remove will put the queue in a bad
// state and all further operations will be undefined.
func (q *Queue) Remove() {
	q.buf[q.head] = nil
	q.head = (q.head + 1) % len(q.buf)
	q.count--
	if len(q.buf) > minQueueLen && q.count*4 <= len(q.buf) {
		q.resize()
	}
}
