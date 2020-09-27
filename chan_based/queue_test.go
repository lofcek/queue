package chan_based

import (
	"context"
	"testing"
	"time"

	"github.com/matryer/is"
)

func TestPushPop(t *testing.T) {
	is := is.New(t)

	pushC := make(chan string)
	defer close(pushC)
	popC := NewQueue(context.Background(), pushC)

	pushC <- "1"
	pushC <- "2"

	is.Equal(<-popC, "1")
	is.Equal(<-popC, "2")

	pushC <- "3"
	is.Equal(<-popC, "3")

	
}

func TestBlockEmpty(t *testing.T) {
	is := is.New(t)

	pushC := make(chan string)
	defer close(pushC)
	popC := NewQueue(context.Background(), pushC)

	go func() {
		time.Sleep(10 * time.Millisecond)
		pushC <- "1"
	}()

	is.Equal(<-popC, "1")
}

func TestFinishUnblock(t *testing.T) {
	is := is.New(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pushC := make(chan string)
	defer close(pushC)
	popC := NewQueue(ctx, pushC)

	cancel()
	_, ok := <-popC
	is.Equal(ok, false)
}

func TestClosePushCUnblock(t *testing.T) {
	is := is.New(t)

	pushC := make(chan string)
	popC := NewQueue(context.Background(), pushC)
	close(pushC)

	_, ok := <-popC
	is.Equal(ok, false)
}