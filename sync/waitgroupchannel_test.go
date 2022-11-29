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
				time.Sleep(100 * time.Second)
				wg.Done()
			}
		}()

		<-wg.Channel()

		wg.Add(5)
		go func() {
			for i := 0; i < 5; i++ {
				time.Sleep(100 * time.Second)
				wg.Done()
			}
		}()

		<-wg.Channel()

		// Expected to return right away
		<-wg.Channel()

		done <- struct{}{}
	}()

	select {
	case <-done:
		t.Log("Success")
	case <-timeout.C:
		t.Error("timeout")
	}
}
