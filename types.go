package main

const (
	BunnyWidth           = 32.0
	BunnyHeight          = 32.0
	MaxLives             = 3
	BunnySpeed           = 4.0
	SpeedBoostMultiplier = 2.0
	JumpStrength         = -12.0
	Gravity              = 0.5
	MaxFallSpeed         = 12.0
	SpeedBoostDuration   = 3.0
	GroundY              = 500.0
	CarrotPoints         = 10
	BasePointsPerLevel   = 100

	CoyoteTime      = 0.12
	JumpBufferTime  = 0.12
	InvincibleTime  = 1.5
	ScreenShakeTime = 0.3
	ScreenShakeMag  = 6.0
	TargetFPS       = 60.0
	DeltaTime       = 1.0 / TargetFPS
)

type Rect struct {
	X      float64
	Y      float64
	Width  float64
	Height float64
}

func (r Rect) Intersects(other Rect) bool {
	return r.X < other.X+other.Width &&
		r.X+r.Width > other.X &&
		r.Y < other.Y+other.Height &&
		r.Y+r.Height > other.Y
}

func (r Rect) ContainsPoint(px, py float64) bool {
	return px >= r.X && px < r.X+r.Width && py >= r.Y && py < r.Y+r.Height
}
