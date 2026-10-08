package main

import (
	"github.com/Zyko0/go-sdl3/sdl"
)

type TimeManager struct {
	TargetFPS    float64
	DT           float32
	FPS          float64
	lastTimeNS   uint64
	frameStartNS uint64
	frameCount   int
	fpsTimerNS   uint64
}

// MARK: TIMERS
func uTimers() {
	if dropTimer > 0 {
		dropTimer--
	}

	if len(lev2d.tri) > 0 || len(lev2d.rec) > 0 {
		if selectFadeSwitch {
			if selectFade < 200 {
				selectFade += 10
			} else {
				selectFadeSwitch = false
			}
		} else {
			if selectFade > 50 {
				selectFade -= 10
			} else {
				selectFadeSwitch = true
			}
		}
	}

}

// MARK: FPS TIME MANAGER
func NewTimeManager(targetFPS float64) *TimeManager {
	now := sdl.TicksNS() // Uses sdl.TicksNS()
	return &TimeManager{
		TargetFPS:    targetFPS,
		lastTimeNS:   now,
		frameStartNS: now,
		fpsTimerNS:   now,
	}
}

func (tm *TimeManager) BeginFrame() {
	tm.frameStartNS = sdl.TicksNS()
	elapsedNS := tm.frameStartNS - tm.lastTimeNS
	tm.DT = float32(elapsedNS) / 1_000_000_000.0
	tm.lastTimeNS = tm.frameStartNS

	// Clamp max DT to prevent physics explosion during window dragging
	if tm.DT > 0.1 {
		tm.DT = 0.1
	}

	tm.frameCount++
	if tm.frameStartNS-tm.fpsTimerNS >= 1_000_000_000 {
		tm.FPS = float64(tm.frameCount)
		tm.frameCount = 0
		tm.fpsTimerNS = tm.frameStartNS
	}
}

func (tm *TimeManager) EndFrame() {
	if tm.TargetFPS <= 0 {
		return
	}
	targetDurationNS := uint64(1_000_000_000.0 / tm.TargetFPS)
	workTimeNS := sdl.TicksNS() - tm.frameStartNS

	if workTimeNS < targetDurationNS {
		sdl.DelayNS(targetDurationNS - workTimeNS)
	}
}
