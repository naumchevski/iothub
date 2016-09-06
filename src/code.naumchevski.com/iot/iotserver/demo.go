package main

import ()

const (
	demoPrivateCounter1 = "xxxx-iot0-xxxx-1111"
	demoPublicCounter1  = "xxxx-iot0-1111"

	demoPrivateCounter2 = "yyyy-iot0-yyyy-1111"
	demoPublicCounter2  = "yyyy-iot0-1111"
)

func startDemo() {
	privateCounters.Add(demoPrivateCounter1, 0)
	privateCounters.Add(demoPrivateCounter2, 0)
	
	publicCounters.Add(demoPublicCounter1, demoPrivateCounter1)
	publicCounters.Add(demoPublicCounter2, demoPrivateCounter2)
}
