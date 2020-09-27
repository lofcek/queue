package chan_based

import "context"

// Again queue of string, however based on channels

// NewQueue creates a new queue. String could be sent into push and received in returned channel
func NewQueue(ctx context.Context, pushC <-chan string) (popC <-chan string) {
	resultC := make(chan string)

	go func() {
		defer close(resultC)
		var queue []string
		for {
			select {
			case <-ctx.Done():
				return
			case s, ok := <-pushC:
				if !ok {
					return
				}
				queue = append(queue, s)
			case maybe(resultC, len(queue) > 0) <- first(queue):
				queue = queue[1:]
			}
		}
	}()
	return resultC
}

func first(q []string) string {
	if len(q) > 0 {
		return q[0]
	}
	return ""
}

func maybe(c chan<- string, cond bool) chan<- string {
	if cond {
		return c
	}
	return nil
}
