// Copyright (c) 2021 Nikos Leivadaris
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

// hlclock provides functions for Hybrid Logical Clocks.
package hlclock

import (
	"encoding/json"
	"fmt"
)

// Now implementation provides the current wall clock time.
type SysClock interface {
	Now() int64
}

// HTimestamp contains the timestamp information.
type HTimestamp struct {
	timestamp int64
	counter   uint16
}

// NewHTimestamp takes the current wall clock time and
// initial counter and returns a new HTimestamp.
func NewHTimestamp(timestamp int64, counter uint16) HTimestamp {
	return HTimestamp{timestamp: timestamp, counter: counter}
}

// Increment takes the current wall clock time and
// updates the timestamp on local event.
func (ht *HTimestamp) Increment(pt int64) {
	if pt > ht.timestamp {
		ht.timestamp = pt
		ht.counter = 0
	} else {
		ht.counter += 1
	}
}

// Merge takes the current wall clock time and a remote timestamp and
// updates the current timestamp.
func (ht *HTimestamp) Merge(pt int64, msg *HTimestamp) {
	switch {
	case pt > ht.timestamp && pt > msg.timestamp:
		ht.timestamp = pt
		ht.counter = 0
	case ht.timestamp == msg.timestamp:
		ht.counter = max(ht.counter, msg.counter) + 1
	case msg.timestamp > ht.timestamp:
		ht.timestamp = msg.timestamp
		ht.counter = msg.counter + 1
	default:
		ht.counter += 1
	}
}

func (ht *HTimestamp) Copy() HTimestamp {
	return HTimestamp{
		ht.timestamp,
		ht.counter,
	}
}

func (ht *HTimestamp) String() string {
	return fmt.Sprintf("Timestamp:{clock: %d, counter: %d}", ht.timestamp, ht.counter)
}

func (ht *HTimestamp) Equal(other *HTimestamp) bool {
	return ht.timestamp == other.timestamp && ht.counter == other.counter
}

func (ht *HTimestamp) Compare(other *HTimestamp) int {
	switch {
	case ht.Equal(other):
		return 0
	case ht.timestamp == other.timestamp:
		if ht.counter < other.counter {
			return -1
		}
		return 1
	case ht.timestamp < other.timestamp:
		return -1
	default:
		return 1
	}
}

func (ht *HTimestamp) MarshalJSON() ([]byte, error) {
	j, err := json.Marshal(struct {
		Timestamp int64  `json:"timestamp"`
		Counter   uint16 `json:"counter"`
	}{
		Timestamp: ht.timestamp,
		Counter:   ht.counter,
	})

	if err != nil {
		return nil, err
	}

	return j, nil
}

func (ht *HTimestamp) UnmarshalJSON(data []byte) error {
	v := struct {
		Timestamp int64  `json:"timestamp"`
		Counter   uint16 `json:"counter"`
	}{}

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	ht.timestamp = v.Timestamp
	ht.counter = v.Counter

	return nil
}

func (ht *HTimestamp) Timestamp() int64 {
	return ht.timestamp
}

func (ht *HTimestamp) Counter() uint16 {
	return ht.counter
}

type HCLock struct {
	sysClock SysClock
	latest   *HTimestamp
}

func New(nodeID string, sysClock SysClock) HCLock {
	tm := NewHTimestamp(sysClock.Now(), 0)
	
	return HCLock{sysClock, &tm}
}

func (hcl *HCLock) Increment() {
	hcl.latest.Increment(hcl.sysClock.Now())
}

func (hcl *HCLock) Merge(e *HTimestamp) {
	hcl.latest.Merge(hcl.sysClock.Now(), e)
}

func (hcl *HCLock) String() string {
	return hcl.latest.String()
}

func (hcl *HCLock) CopyTimestamp() HTimestamp {
	return hcl.latest.Copy()
}

func max(x, y uint16) uint16 {
	if x > y {
		return x
	}
	return y
}
