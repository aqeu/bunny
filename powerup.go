package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type PowerUpType int

const (
	PowerUpCarrot PowerUpType = iota
	PowerUpSpeedBoost
	PowerUpExtraLife
	PowerUpDoubleJump
	PowerUpShield
)

type PowerUp struct {
	X, Y      float64
	Width     float64
	Height    float64
	Type      PowerUpType
	Collected bool
	Sprite    *ebiten.Image
	BobOffset float64
	GlowTimer float64
}

func NewPowerUp(x, y float64, puType PowerUpType) *PowerUp {
	p := &PowerUp{
		X:         x,
		Y:         y,
		Width:     20,
		Height:    20,
		Type:      puType,
		Collected: false,
		Sprite:    CreatePowerUpImage(puType),
		BobOffset: 0,
		GlowTimer: 0,
	}
	return p
}

func (p *PowerUp) GetBounds() Rect {
	return Rect{
		X:      p.X,
		Y:      p.Y + p.BobOffset,
		Width:  p.Width,
		Height: p.Height,
	}
}

func (p *PowerUp) Update() {
	if p.Collected {
		return
	}
	p.BobOffset += 0.1
	p.BobOffset = 3 * math.Sin(p.BobOffset)
	p.GlowTimer += 0.05
}

func (p *PowerUp) Draw(screen *ebiten.Image, cameraX float64) {
	if p.Collected {
		return
	}
	alpha := uint8(100 + 155*math.Abs(math.Sin(p.GlowTimer)))
	glowColor := color.RGBA{255, 255, 255, alpha}
	ebitenutil.DrawRect(screen, p.X-cameraX-4, p.Y+p.BobOffset-4, p.Width+8, p.Height+8, glowColor)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(p.X-cameraX, p.Y+p.BobOffset)
	screen.DrawImage(p.Sprite, op)
}
