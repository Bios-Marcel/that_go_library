package sync

import "sync"

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
	sync.WaitGroup
	channel chan struct{}
}

// NewWaitGroupChannel creates a new WaitGroupChannel instance. Simply call
// WaitGroupChannel.Add() like you would on sync.WaitGroup.
func NewWaitGroupChannel() *WaitGroupChannel {
	wgc := &WaitGroupChannel{
		channel: make(chan struct{}),
	}
	go func() {
		wgc.Wait()
		wgc.channel <- struct{}{}
	}()

	return wgc
}

// Channel returns the channel you can wait on. The channel will never be
// closed, as you can technically reuse the waitgroup. However, if you
// listen to the channel even tho the waitgroup is in a done state and the
// previous channel content has already been consumed or WaitGroup.Add has
// never been called before, you'll wait forever.
func (wgc *WaitGroupChannel) Channel() chan struct{} {
	return wgc.channel
}
