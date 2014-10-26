package buffer

func Chan(in <-chan []byte, next chan<- []byte) {
	pending := NewBufferQueue(NewUnboundedBuffer(32*1024, 100*1024*1024))
	chanQueue(in, next, pending)
}

func MemChan(in <-chan []byte, next chan<- []byte) {
	chanQueue(in, next, NewSliceQueue())
}

func chanQueue(in <-chan []byte, next chan<- []byte, pending Queue) {
	defer close(next)

recv:

	for {

		if Empty(pending) {
			data, ok := <-in
			if !ok {
				break
			}

			pending.Push(data)
		}

		select {
		case data, ok := <-in:
			if !ok {
				break recv
			}
			pending.Push(data)

		case next <- pending.Peek():
			pending.Pop()
		}

	}

	for !Empty(pending) {
		next <- pending.Pop()
	}
}
