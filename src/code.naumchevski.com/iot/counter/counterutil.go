package counter

import (

)

const PrivateCounterLen = 19
const PublicCounterLen = 14

func GetPrivateCounterId(id string, publicCounters counters) (string, bool) {
	if IsPrivateCounter(id) {
		return id, true	
	}
	if IsPublicCounter(id) {
		privId := publicCounters.Get(id)
		if privId != nil {
			return privId.(string), true 	
		}
		
	}
	return "", false
}

func IsPrivateCounter(id string) bool {
	return len(id) == PrivateCounterLen
}

func IsPublicCounter(id string) bool {
	return len(id) == PublicCounterLen
}
