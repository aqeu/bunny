package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type ObstacleType int

const (
	ObstaclePipe ObstacleType = iota
	ObstacleTunnel
	ObstaclePlatform
	ObstacleSpike
	ObstacleMovingPlatform
)

type Obstacle struct {
	X, Y      float64
	Width     float64
	Height    float64
	Type      ObstacleType
	Sprite    *ebiten.Image
	StartX    float64
	StartY    float64
	MoveRange float64
	MoveSpeed float64
	MoveAxis  string
	Phase     float64
}

func NewObstacle(x, y, width, height float64, obsType ObstacleType) *Obstacle {
	o := &Obstacle{
		X:      x,
		Y:      y,
		Width:  width,
		Height: height,
		Type:   obsType,
		Sprite: CreateObstacleImage(obsType),
		StartX: x,
		StartY: y,
	}
	return o
}

func NewMovingObstacle(x, y, width, height float64, obsType ObstacleType, moveRange, moveSpeed float64, moveAxis string) *Obstacle {
	o := NewObstacle(x, y, width, height, obsType)
	o.MoveRange = moveRange
	o.MoveSpeed = moveSpeed
	o.MoveAxis = moveAxis
	o.Phase = 0
	return o
}

func (o *Obstacle) GetBounds() Rect {
	return Rect{
		X:      o.X,
		Y:      o.Y,
		Width:  o.Width,
		Height: o.Height,
	}
}

func (o *Obstacle) Update() {
	if o.Type == ObstacleMovingPlatform {
		o.Phase += o.MoveSpeed
		if o.MoveAxis == "x" {
			o.X = o.StartX + math.Sin(o.Phase)*o.MoveRange
		} else if o.MoveAxis == "y" {
			o.Y = o.StartY + math.Sin(o.Phase)*o.MoveRange
		}
	}
}

func (o *Obstacle) Draw(screen *ebiten.Image, cameraX float64) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(o.X-cameraX, o.Y)
	screen.DrawImage(o.Sprite, op)
}
