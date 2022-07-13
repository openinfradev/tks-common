package helper

import (
	"math/rand"
)

var (
	CLUSTER_ID_PREFIX  = "c"
	CONTRACT_ID_PREFIX = "p"
	ID_LENGTH          = 9
)

// lowercase RFC 1123
const letterBytes = "abcdefghijklmnopqrstuvwxyz0123456789"

func GenerateClusterId() string {
	return CLUSTER_ID_PREFIX + randStringBytesRmndr(ID_LENGTH-1)
}

func GenerateContractId() string {
	return CONTRACT_ID_PREFIX + randStringBytesRmndr(ID_LENGTH-1)
}

func ValidateClusterId(id string) bool {
	if id == "" {
		return false
	}

	if id[0:1] != CLUSTER_ID_PREFIX {
		return false
	}
	return validateId(id)
}

func ValidateContractId(id string) bool {
	if id == "" {
		return false
	}

	if id[0:1] != CONTRACT_ID_PREFIX {
		return false
	}
	return validateId(id)
}

func validateId(id string) bool {
	return len(id) == ID_LENGTH
}

func randStringBytesRmndr(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}
