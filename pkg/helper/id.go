package helper

import (
	"math/rand"
)

var (
	CLUSTER_ID_PREFIX  = "C"
	CONTRACT_ID_PREFIX = "P"
	ID_LENGTH          = 9
)

func GenerateClusterId() string {
	return CLUSTER_ID_PREFIX + RandStringBytesRmndr(ID_LENGTH)
}

func GenerateContractId() string {
	return CONTRACT_ID_PREFIX + RandStringBytesRmndr(ID_LENGTH)
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
	if len(id) != ID_LENGTH {
		return false
	}
	return true
}

func RandStringBytesRmndr(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n-1)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

/*
func randomBase64String(l int) string {
	buff := make([]byte, int(math.Ceil(float64(l)/float64(1.33333333333))))
	rand.Read(buff)
	str := base64.RawURLEncoding.EncodeToString(buff)
	return str[:l-1]
}
*/
