package utils

// A Snowflake ID is composed of
//
//	39 bits for time in units of 10 msec = 549755813887
//	 8 bits for a sequence number = 255
//	 6 bits for a machine id  = 63

import (
	"errors"
	"sync"
	"time"
)

const (
	BitLenTime      = 39                               // bit length of time
	BitLenSequence  = 8                                // bit length of sequence number
	BitLenMachineID = 53 - BitLenTime - BitLenSequence // bit length of machine id
)

type Settings struct {
	StartTime time.Time
	MachineID uint16
}

// Snowflake is a distributed unique ID generator.
type Snowflake struct {
	mutex       *sync.Mutex
	startTime   int64
	elapsedTime int64
	sequence    uint16
	machineID   uint16
}

func NewSnowflake(machineID uint16, t ...time.Time) (*Snowflake, error) {
	sf := new(Snowflake)
	sf.mutex = new(sync.Mutex)
	sf.sequence = uint16(1<<BitLenSequence - 1)
	if len(t) > 0 {
		if t[0].After(time.Now()) {
			return nil, errors.New("start time is less than current time")
		}
		sf.startTime = toSnowflakeTime(t[0])
	} else {
		sf.startTime = toSnowflakeTime(time.Date(2023, 4, 5, 0, 0, 0, 0, time.UTC))
	}
	if machineID <= 0 || machineID > 63 {
		return nil, errors.New("machine id range must be between 1 and 63")
	}
	sf.machineID = machineID
	return sf, nil
}

// NextID generates a next unique ID.
// After the Snowflake time overflows, NextID returns an error.
func (sf *Snowflake) NextID() (uint64, error) {
	const maskSequence = uint16(1<<BitLenSequence - 1)
	sf.mutex.Lock()
	defer sf.mutex.Unlock()

	current := currentElapsedTime(sf.startTime)
	if sf.elapsedTime < current {
		sf.elapsedTime = current
		sf.sequence = 0
	} else { // sf.elapsedTime >= current
		sf.sequence = (sf.sequence + 1) & maskSequence
		if sf.sequence == 0 {
			sf.elapsedTime++
			overtime := sf.elapsedTime - current
			time.Sleep(sleepTime((overtime)))
		}
	}

	return sf.toID()
}

const snowflakeTimeUnit = 1e7 // nsec, i.e. 10 msec

func toSnowflakeTime(t time.Time) int64 {
	return t.UTC().UnixNano() / snowflakeTimeUnit
}

func currentElapsedTime(startTime int64) int64 {
	return toSnowflakeTime(time.Now()) - startTime
}

func sleepTime(overtime int64) time.Duration {
	return time.Duration(overtime*snowflakeTimeUnit) -
		time.Duration(time.Now().UTC().UnixNano()%snowflakeTimeUnit)
}

func (sf *Snowflake) toID() (uint64, error) {
	if sf.elapsedTime >= 1<<BitLenTime {
		return 0, errors.New("over the time limit")
	}

	return uint64(sf.elapsedTime)<<(BitLenSequence+BitLenMachineID) |
		uint64(sf.sequence)<<BitLenMachineID |
		uint64(sf.machineID), nil
}
