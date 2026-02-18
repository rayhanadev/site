package render

import "time"

type streamer struct {
	speeds []int
}

func newStreamer(defaultSpeed int) *streamer {
	return &streamer{speeds: []int{defaultSpeed}}
}

func (s *streamer) push(speed int) {
	s.speeds = append(s.speeds, speed)
}

func (s *streamer) pop() {
	if len(s.speeds) > 1 {
		s.speeds = s.speeds[:len(s.speeds)-1]
	}
}

func (s *streamer) delay() {
	time.Sleep(time.Duration(s.speeds[len(s.speeds)-1]) * time.Millisecond)
}
