package queue

import (
	"testing"
	"time"

	"github.com/matryer/is"
)


func TestPushPop(t *testing.T) {
	is := is.New(t)

	q := New()
	q.Push("1")
	q.Push("2")

	s, err := q.Pop()
	is.NoErr(err)
	is.Equal(s, "1")

	s, err = q.Pop()
	is.NoErr(err)
	is.Equal(s, "2")

	q.Push("3")
	s, err = q.Pop()
	is.NoErr(err)
	is.Equal(s, "3")

	q.Finish()
	s, err = q.Pop()
	is.Equal(s, "")
	is.Equal(err, ErrDone)
}
func TestBlockEmpty(t *testing.T) {
	is := is.New(t)

	q := New()
	go func() {
		time.Sleep(10 * time.Millisecond)
		q.Push("1")
	}()

	s, err := q.Pop()
	is.NoErr(err)
	is.Equal(s, "1")
}

func TestFinishUnblock(t *testing.T) {
	q := New()
	afterFinish := false
	go func() {
		time.Sleep(10 * time.Millisecond)
		afterFinish = true
		q.Finish()
		q.Push("1")
	}()

	for i := 0; i < 10; i++ {
		s, err := q.Pop()
		if s != "" || err != ErrDone || !afterFinish {
			t.Fatalf("wrong result %q %s", s, err)
		}
	}
}
