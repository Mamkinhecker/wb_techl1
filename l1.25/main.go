package main

import "time"

func Sleep(d time.Duration) {
	select {
	case <-time.After(d):
	}
}
