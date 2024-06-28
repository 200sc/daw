// Package timing provides utilities for time.
package timing

import (
	"math"
	"time"
)

const (
	nanoPerSecond = 1000000000

	maximumFPS = 1200
)

// FPS returns the number of frames being processed per second,
// supposing a time interval from lastTime to now.
func FPS(lastTime, now time.Time) float64 {
	fps := 1 / now.Sub(lastTime).Seconds()
	// This indicates that we recorded two times within
	// the inaccuracy of the OS's system clock, so the values
	// were the same. 1200 is chosen because on windows,
	// fps will be 1200 instead of a negative value.
	if int(fps) < 0 {
		return maximumFPS
	}
	return fps
}

// FPSToNano converts a framesPerSecond value to the number of
// nanoseconds that should take place for each frame.
func FPSToNano(fps float64) int64 {
	if fps == 0.0 {
		return math.MaxInt64
	}
	return int64(nanoPerSecond / fps)
}

// FPSToFrameDelay converts a frameRate like 60fps into a delay time between frames
func FPSToFrameDelay(frameRate int) time.Duration {
	if frameRate == 0 {
		return time.Duration(math.MaxInt64)
	}
	return time.Second / time.Duration(int64(frameRate))
}

// FrameDelayToFPS converts a duration of delay between frames into a frames per second count
func FrameDelayToFPS(dur time.Duration) float64 {
	if dur == 0 {
		return math.MaxFloat64
	}
	return float64(time.Second) / float64(dur)
}
