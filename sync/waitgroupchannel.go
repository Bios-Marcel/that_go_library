package sync

import (
	"sync"
)

// WaitGroupChannel is a wrapper around sync.WaitGroup, that allow you to
// wait for its completion via a channel instead of calling WaitGroup.Wait().
//
// This allows you to do the following:
//
//	select {
//	case <-waitGroupChannel.Channel():
//	case timeout.C:
//	}
//
// In order to achieve this, a go routine is started and waits for the wrapped
// waitgroup to be done.
type WaitGroupChannel struct {
	wg      *sync.WaitGroup
	mutex   *sync.Mutex
	waiting bool
	channel chan struct{}
}

// NewWaitGroupChannel creates a new WaitGroupChannel instance. Simply call
// WaitGroupChannel.Add() like you would on sync.WaitGroup.
func NewWaitGroupChannel() *WaitGroupChannel {
	wgc := &WaitGroupChannel{
		wg:      &sync.WaitGroup{},
		channel: make(chan struct{}, 1),
		mutex:   &sync.Mutex{},
	}
	// Initially, the channel is non-nil, but closed, so you won't block when
	// waiting and also no cause a nil pointer dereference.
	close(wgc.channel)

	return wgc
}

func (wgc *WaitGroupChannel) Add(delta int) {
	wgc.mutex.Lock()
	defer wgc.mutex.Unlock()

	wgc.wg.Add(delta)
	if delta > 0 && !wgc.waiting {
		wgc.feedChannel()
	}
}

func (wgc *WaitGroupChannel) Done() {
	wgc.wg.Done()
}

func (wgc *WaitGroupChannel) feedChannel() {
	wgc.waiting = true
	wgc.channel = make(chan struct{}, 1)

	go func() {
		wgc.wg.Wait()

		wgc.mutex.Lock()
		defer wgc.mutex.Unlock()

		wgc.waiting = false
		wgc.channel <- struct{}{}
		// Close channel, so sthat follow up waits neither block, nor cause a nil
		// pointer dereference. See initialisation.
		close(wgc.channel)
	}()
}

// Channel returns the channel you can wait on. The channel will never be
// closed, as you can technically reuse the waitgroup. However, if you
// listen to the channel even tho the waitgroup is in a done state and the
// previous channel content has already been consumed or WaitGroup.Add has
// never been called before, you'll wait forever.
func (wgc *WaitGroupChannel) Channel() chan struct{} {
	wgc.mutex.Lock()
	defer wgc.mutex.Unlock()
	return wgc.channel
}
