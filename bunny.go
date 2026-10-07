package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Bunny struct {
	X, Y               float64
	VX, VY             float64
	Width              float64
	Height             float64
	Grounded           bool
	Lives              int
	SpeedBoostTimer    float64
	SpeedBoostDuration float64
	DoubleJump         bool
	DoubleJumpUsed     bool
	Shield             bool
	InvincibleTimer    float64
	CoyoteTimer        float64
	JumpBufferTimer    float64
	animFrame          int
	animTimer          float64
	facingRight        bool
	sprites            [4]*ebiten.Image
	dead               bool
	landingTimer       float64
	deathFlashTimer    float64
}

func NewBunny(x, y float64) *Bunny {
	b := &Bunny{
		X:           x,
		Y:           y,
		VX:          0,
		VY:          0,
		Width:       BunnyWidth,
		Height:      BunnyHeight,
		Lives:       MaxLives,
		Grounded:    true,
		facingRight: true,
	}
	for i := 0; i < 4; i++ {
		b.sprites[i] = CreateBunnyImageFrame(i)
	}
	return b
}

func (b *Bunny) GetBounds() Rect {
	return Rect{
		X:      b.X,
		Y:      b.Y,
		Width:  b.Width,
		Height: b.Height,
	}
}

func (b *Bunny) Update(keys map[ebiten.Key]bool, levelWidth float64, obstacles []*Obstacle) {
	if b.dead {
		b.landingTimer += DeltaTime
		b.deathFlashTimer += DeltaTime
		return
	}

	if b.InvincibleTimer > 0 {
		b.InvincibleTimer -= DeltaTime
		if b.InvincibleTimer < 0 {
			b.InvincibleTimer = 0
		}
	}

	speed := BunnySpeed
	if b.SpeedBoostTimer > 0 {
		speed = BunnySpeed * SpeedBoostMultiplier
		b.SpeedBoostTimer -= DeltaTime
		if b.SpeedBoostTimer <= 0 {
			b.SpeedBoostTimer = 0
		}
	}

	b.Grounded = false

	desiredVX := 0.0
	if keys[ebiten.KeyLeft] || keys[ebiten.KeyA] {
		desiredVX = -speed
		b.facingRight = false
	} else if keys[ebiten.KeyRight] || keys[ebiten.KeyD] {
		desiredVX = speed
		b.facingRight = true
	}
	b.X += desiredVX
	b.resolveHorizontalCollisions(obstacles)
	if b.X < 0 {
		b.X = 0
	}
	if b.X > levelWidth-b.Width {
		b.X = levelWidth - b.Width
	}

	wantJump := keys[ebiten.KeySpace] || keys[ebiten.KeyUp] || keys[ebiten.KeyW]
	if wantJump {
		b.JumpBufferTimer = JumpBufferTime
	} else if b.JumpBufferTimer > 0 {
		b.JumpBufferTimer -= DeltaTime
	}

	b.VY += Gravity
	if b.VY > MaxFallSpeed {
		b.VY = MaxFallSpeed
	}
	b.Y += b.VY
	if b.Y >= GroundY-b.Height {
		b.Y = GroundY - b.Height
		b.VY = 0
		b.Grounded = true
	}
	b.resolvePlatformCollisions(obstacles)

	if b.Grounded {
		b.CoyoteTimer = CoyoteTime
		b.DoubleJumpUsed = false
	} else {
		b.CoyoteTimer -= DeltaTime
	}

	if b.JumpBufferTimer > 0 && b.CoyoteTimer > 0 {
		b.VY = JumpStrength
		b.Grounded = false
		b.CoyoteTimer = 0
		b.JumpBufferTimer = 0
		PlayJumpSound()
	} else if b.JumpBufferTimer > 0 && b.DoubleJump && !b.DoubleJumpUsed && !b.Grounded {
		b.VY = JumpStrength * 0.9
		b.DoubleJumpUsed = true
		b.JumpBufferTimer = 0
		PlayDoubleJumpSound()
	}

	b.animTimer += DeltaTime
	if b.animTimer > 0.15 {
		b.animTimer = 0
		b.animFrame = (b.animFrame + 1) % 4
	}
}

func (b *Bunny) resolveHorizontalCollisions(obstacles []*Obstacle) {
	bunnyBounds := b.GetBounds()
	for _, obs := range obstacles {
		if obs.Type == ObstaclePlatform || obs.Type == ObstacleMovingPlatform {
			continue
		}
		obsBounds := obs.GetBounds()
		if bunnyBounds.Intersects(obsBounds) {
			overlapLeft := (b.X + b.Width) - obs.X
			overlapRight := (obs.X + obs.Width) - b.X
			if overlapLeft < overlapRight {
				b.X = obs.X - b.Width
			} else {
				b.X = obs.X + obs.Width
			}
		}
	}
}

func (b *Bunny) resolvePlatformCollisions(obstacles []*Obstacle) {
	if b.VY < 0 {
		return
	}
	bunnyBounds := b.GetBounds()
	for _, obs := range obstacles {
		if obs.Type != ObstaclePlatform && obs.Type != ObstacleMovingPlatform {
			continue
		}
		if obs.Y >= GroundY-1 && obs.Type == ObstaclePlatform {
			continue
		}
		obsBounds := obs.GetBounds()
		if bunnyBounds.Intersects(obsBounds) {
			prevBottom := (b.Y - b.VY) + b.Height
			if prevBottom <= obs.Y+4 {
				b.Y = obs.Y - b.Height
				b.VY = 0
				b.Grounded = true
				return
			}
		}
	}
}

func (b *Bunny) Reset(x, y float64) {
	b.X = x
	b.Y = y
	b.VX = 0
	b.VY = 0
	b.Grounded = true
	b.SpeedBoostTimer = 0
	b.DoubleJump = false
	b.DoubleJumpUsed = false
	b.Shield = false
	b.InvincibleTimer = InvincibleTime
	b.CoyoteTimer = 0
	b.JumpBufferTimer = 0
	b.dead = false
	b.landingTimer = 0
	b.deathFlashTimer = 0
	b.animFrame = 0
	b.animTimer = 0
	b.facingRight = true
}

func (b *Bunny) Draw(screen *ebiten.Image, cameraX float64) {
	if b.dead {
		flashOn := int(b.deathFlashTimer*10)%2 == 0
		if flashOn {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(b.X-cameraX, b.Y)
			screen.DrawImage(b.sprites[0], op)
		}
		return
	}
	if b.InvincibleTimer > 0 && int(b.InvincibleTimer*20)%2 == 0 {
		return
	}
	op := &ebiten.DrawImageOptions{}
	if !b.facingRight {
		op.GeoM.Scale(-1, 1)
		op.GeoM.Translate(-b.Width, 0)
	}
	op.GeoM.Translate(b.X-cameraX, b.Y)
	screen.DrawImage(b.sprites[b.animFrame], op)

	if b.Shield {
		shieldOp := &ebiten.DrawImageOptions{}
		shieldOp.GeoM.Translate(b.X-cameraX-2, b.Y-2)
		ebitenutil.DrawRect(screen, b.X-cameraX-2, b.Y-2, b.Width+4, b.Height+4, color.RGBA{100, 100, 255, 100})
	}
}

func (b *Bunny) Die() {
	b.Lives--
	b.dead = true
	b.landingTimer = 0
	b.deathFlashTimer = 0
}

func (b *Bunny) AddLife() {
	if b.Lives < MaxLives {
		b.Lives++
	}
}

func (b *Bunny) ActivateSpeedBoost() {
	b.SpeedBoostTimer = SpeedBoostDuration
}

func (b *Bunny) ActivateDoubleJump() {
	b.DoubleJump = true
	b.DoubleJumpUsed = false
}

func (b *Bunny) ActivateShield() {
	b.Shield = true
}

func (b *Bunny) UseShield() {
	b.Shield = false
}
