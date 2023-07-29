package state

import (
	"crypto/rand"
	"encoding/binary"
	"time"
)

var AppInitiatedAt time.Time

type State struct {
	State      bool
	TargetTime *time.Time
	Interval   time.Duration
}

type Exists bool

func (s *State) SetByChance(chanceOfSuccess float32) {
	buf := make([]byte, 4)
	_, err := rand.Read(buf)
	if err != nil {
		panic(err)
	}

	// Convert the bytes to a uint32
	num := binary.BigEndian.Uint32(buf)

	// Divide by the maximum value a uint32 can have to get a float32 between 0 and 1
	randomFloat := float32(num) / float32(^uint32(0))

	s.SetState(randomFloat <= chanceOfSuccess)
}

func (s *State) SetState(newState bool) {
	s.State = newState
}

func (s *State) IsGood() bool {
	return s.State
}

func (s *State) IsBad() bool {
	return !s.IsGood()
}

func (s *State) NextInterval() (Exists, *time.Duration) {
	if s.TargetTime == nil {
		return true, &s.Interval
	}

	diff := time.Since(*s.TargetTime)

	if diff < 0 {
		return false, nil
	}

	if diff < s.Interval {
		lastDiff := s.Interval - diff
		return true, &lastDiff
	}

	return true, &s.Interval
}
