package sync

import (
	"testing"
	"time"
)

func Test_MultiWait(t *testing.T) {
	wg := NewWaitGroupChannel()

	done := make(chan struct{})
	timeout := time.NewTimer(2 * time.Second)

	go func() {
		wg.Add(5)
		go func() {
			for i := 0; i < 5; i++ {
				time.Sleep(100 * time.Millisecond)
				wg.Done()
			}
		}()

		<-wg.Channel()
		t.Log("Returned 1")

		wg.Add(5)
		go func() {
			for i := 0; i < 5; i++ {
				time.Sleep(100 * time.Millisecond)
				wg.Done()
			}
		}()

		<-wg.Channel()
		t.Log("Returned 2")

		// Expected to return right away
		<-wg.Channel()
		t.Log("Returned 3")

		done <- struct{}{}
	}()

	select {
	case <-done:
		t.Log("Success")
	case <-timeout.C:
		t.Error("timeout")
	}
}
