package stahl4

import (
	"errors"
	"fmt"
	"strconv"
)

func MovementSleep() {
	return
	// TODO: the following code is untested, hence we won't enable it for the
	//  ctf and has led to errors in the past
	//generator, err := rand.Int(rand.Reader, big.NewInt(75))
	//if err != nil {
	//	// default value if something goes wrong with random
	//	time.Sleep(30 * time.Second)
	//} else {
	//	seconds := uint8(generator.Int64() + 15)
	//	time.Sleep(time.Duration(seconds) * time.Second)
	//}
}

func CompareBroadcasts(broadcast1 BroadcastCommand, broadcast2 BroadcastCommand) bool {
	// compares two broadcasts using the ID.
	id1 := broadcast1.GetID()
	id2 := broadcast2.GetID()
	return id1 == id2
}

func GetIntValueFromMap(data map[string]string, key string) (int, error) {
	// get the value of the key
	value := data[key]
	// if the value is empty
	if value == "" {
		// return an error
		return 0, errors.New(fmt.Sprintf("key: %s, is not a valid integer but empty", key))
	}
	// return the value as integer
	return strconv.Atoi(value)
}

func GetStringValueFromMap(data map[string]string, key string) (string, error) {
	// get the value of the key
	value := data[key]
	// if the value is empty
	if value == "" {
		// return an error
		return "", errors.New(fmt.Sprintf("key: %s, is not a valid string but empty", key))
	} else {
		// return the value as string
		return value, nil
	}
}
