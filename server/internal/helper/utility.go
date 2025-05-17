package helper

import (
	"crypto/rand"
	"strconv"
)

func RandomNumber(length int) (int, error) {

	const numbers = "1234567890"

	buffer := make([]byte, length)
	_, err := rand.Read(buffer)
	if err != nil {
		return 0, err
	}

	numLength := len(numbers)

	for i := range length {
		buffer[i] = numbers[int(buffer[i])%numLength]
	}

	return strconv.Atoi(string(buffer))
}
