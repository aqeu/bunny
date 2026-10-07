package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Level struct {
	Width            float64
	GroundY          float64
	Obstacles        []*Obstacle
	PowerUps         []*PowerUp
	FinishX          float64
	CameraX          float64
	BackgroundLayers [3]*ebiten.Image
	LevelNumber      int
	SpawnX           float64
	SpawnY           float64
	Checkpoints      []Checkpoint
	respawnX         float64
	respawnY         float64
}

type Checkpoint struct {
	X, Y   float64
	Active bool
	Sprite *ebiten.Image
}

func NewLevel(levelNum int) *Level {
	l := &Level{
		GroundY:     GroundY,
		LevelNumber: levelNum,
		respawnX:    -1,
		respawnY:    -1,
	}

	switch levelNum {
	case 1:
		l.initLevel1()
	case 2:
		l.initLevel2()
	case 3:
		l.initLevel3()
	default:
		l.initLevel1()
	}

	l.BackgroundLayers = GenerateBackgroundLayers()
	return l
}

func (l *Level) GetRespawnPoint() Checkpoint {
	if l.respawnX >= 0 && l.respawnY >= 0 {
		return Checkpoint{X: l.respawnX, Y: l.respawnY}
	}
	return Checkpoint{X: l.SpawnX, Y: l.SpawnY}
}

func (l *Level) initLevel1() {
	l.Width = 4000
	l.SpawnX = 50
	l.SpawnY = GroundY - BunnyHeight
	l.FinishX = l.Width - 100

	l.Obstacles = append(l.Obstacles, NewObstacle(0, GroundY, l.Width, 60, ObstaclePlatform))

	l.Obstacles = append(l.Obstacles, NewObstacle(400, GroundY-80, 40, 80, ObstaclePipe))
	l.Obstacles = append(l.Obstacles, NewObstacle(700, GroundY-80, 40, 80, ObstaclePipe))
	l.Obstacles = append(l.Obstacles, NewObstacle(1100, GroundY-80, 40, 80, ObstaclePipe))

	l.Obstacles = append(l.Obstacles, NewObstacle(550, GroundY-60, 80, 16, ObstaclePlatform))
	l.Obstacles = append(l.Obstacles, NewObstacle(900, GroundY-100, 80, 16, ObstaclePlatform))
	l.Obstacles = append(l.Obstacles, NewObstacle(1300, GroundY-60, 80, 16, ObstaclePlatform))

	l.Obstacles = append(l.Obstacles, NewObstacle(1600, GroundY-100, 60, 80, ObstacleTunnel))
	l.Obstacles = append(l.Obstacles, NewObstacle(2000, GroundY-100, 60, 80, ObstacleTunnel))

	l.Obstacles = append(l.Obstacles, NewObstacle(2400, GroundY-16, 24, 16, ObstacleSpike))
	l.Obstacles = append(l.Obstacles, NewObstacle(2450, GroundY-16, 24, 16, ObstacleSpike))
	l.Obstacles = append(l.Obstacles, NewObstacle(2500, GroundY-16, 24, 16, ObstacleSpike))

	l.Obstacles = append(l.Obstacles, NewObstacle(2800, GroundY-80, 40, 80, ObstaclePipe))
	l.Obstacles = append(l.Obstacles, NewObstacle(3200, GroundY-80, 40, 80, ObstaclePipe))

	l.Obstacles = append(l.Obstacles, NewMovingObstacle(1800, GroundY-120, 80, 16, ObstacleMovingPlatform, 60, 0.02, "x"))
	l.Obstacles = append(l.Obstacles, NewMovingObstacle(2200, GroundY-150, 80, 16, ObstacleMovingPlatform, 50, 0.025, "y"))

	l.PowerUps = append(l.PowerUps, NewPowerUp(500, GroundY-80, PowerUpCarrot))
	l.PowerUps = append(l.PowerUps, NewPowerUp(1000, GroundY-120, PowerUpCarrot))
	l.PowerUps = append(l.PowerUps, NewPowerUp(1700, GroundY-120, PowerUpSpeedBoost))
	l.PowerUps = append(l.PowerUps, NewPowerUp(2600, GroundY-40, PowerUpCarrot))
	l.PowerUps = append(l.PowerUps, NewPowerUp(3000, GroundY-100, PowerUpExtraLife))
	l.PowerUps = append(l.PowerUps, NewPowerUp(3500, GroundY-80, PowerUpCarrot))
	l.PowerUps = append(l.PowerUps, NewPowerUp(1900, GroundY-160, PowerUpDoubleJump))

	l.Checkpoints = append(l.Checkpoints, Checkpoint{X: 1500, Y: GroundY - BunnyHeight, Active: true})
	l.Checkpoints = append(l.Checkpoints, Checkpoint{X: 3000, Y: GroundY - BunnyHeight, Active: true})
}

func (l *Level) initLevel2() {
	l.Width = 5000
	l.SpawnX = 50
	l.SpawnY = GroundY - BunnyHeight
	l.FinishX = l.Width - 100

	l.Obstacles = append(l.Obstacles, NewObstacle(0, GroundY, l.Width, 60, ObstaclePlatform))

	l.Obstacles = append(l.Obstacles, NewObstacle(350, GroundY-80, 40, 80, ObstaclePipe))
	l.Obstacles = append(l.Obstacles, NewObstacle(600, GroundY-80, 40, 80, ObstaclePipe))
	l.Obstacles = append(l.Obstacles, NewObstacle(950, GroundY-80, 40, 80, ObstaclePipe))
	l.Obstacles = append(l.Obstacles, NewObstacle(1300, GroundY-80, 40, 80, ObstaclePipe))
	l.Obstacles = append(l.Obstacles, NewObstacle(1700, GroundY-80, 40, 80, ObstaclePipe))
	l.Obstacles = append(l.Obstacles, NewObstacle(2100, GroundY-80, 40, 80, ObstaclePipe))
	l.Obstacles = append(l.Obstacles, NewObstacle(2500, GroundY-80, 40, 80, ObstaclePipe))
	l.Obstacles = append(l.Obstacles, NewObstacle(2900, GroundY-80, 40, 80, ObstaclePipe))
	l.Obstacles = append(l.Obstacles, NewObstacle(3300, GroundY-80, 40, 80, ObstaclePipe))
	l.Obstacles = append(l.Obstacles, NewObstacle(3700, GroundY-80, 40, 80, ObstaclePipe))
	l.Obstacles = append(l.Obstacles, NewObstacle(4100, GroundY-80, 40, 80, ObstaclePipe))

	l.Obstacles = append(l.Obstacles, NewObstacle(500, GroundY-70, 60, 16, ObstaclePlatform))
	l.Obstacles = append(l.Obstacles, NewObstacle(800, GroundY-110, 60, 16, ObstaclePlatform))
	l.Obstacles = append(l.Obstacles, NewObstacle(1100, GroundY-70, 60, 16, ObstaclePlatform))
	l.Obstacles = append(l.Obstacles, NewObstacle(1500, GroundY-110, 60, 16, ObstaclePlatform))
	l.Obstacles = append(l.Obstacles, NewObstacle(1900, GroundY-70, 60, 16, ObstaclePlatform))
	l.Obstacles = append(l.Obstacles, NewObstacle(2300, GroundY-110, 60, 16, ObstaclePlatform))
	l.Obstacles = append(l.Obstacles, NewObstacle(2700, GroundY-70, 60, 16, ObstaclePlatform))
	l.Obstacles = append(l.Obstacles, NewObstacle(3100, GroundY-110, 60, 16, ObstaclePlatform))
	l.Obstacles = append(l.Obstacles, NewObstacle(3500, GroundY-70, 60, 16, ObstaclePlatform))
	l.Obstacles = append(l.Obstacles, NewObstacle(3900, GroundY-110, 60, 16, ObstaclePlatform))

	l.Obstacles = append(l.Obstacles, NewObstacle(1400, GroundY-100, 60, 80, ObstacleTunnel))
	l.Obstacles = append(l.Obstacles, NewObstacle(2600, GroundY-100, 60, 80, ObstacleTunnel))
	l.Obstacles = append(l.Obstacles, NewObstacle(3800, GroundY-100, 60, 80, ObstacleTunnel))

	l.Obstacles = append(l.Obstacles, NewObstacle(1000, GroundY-16, 24, 16, ObstacleSpike))
	l.Obstacles = append(l.Obstacles, NewObstacle(1030, GroundY-16, 24, 16, ObstacleSpike))
	l.Obstacles = append(l.Obstacles, NewObstacle(1060, GroundY-16, 24, 16, ObstacleSpike))
	l.Obstacles = append(l.Obstacles, NewObstacle(2200, GroundY-16, 24, 16, ObstacleSpike))
	l.Obstacles = append(l.Obstacles, NewObstacle(2230, GroundY-16, 24, 16, ObstacleSpike))
	l.Obstacles = append(l.Obstacles, NewObstacle(2260, GroundY-16, 24, 16, ObstacleSpike))
	l.Obstacles = append(l.Obstacles, NewObstacle(3400, GroundY-16, 24, 16, ObstacleSpike))
	l.Obstacles = append(l.Obstacles, NewObstacle(3430, GroundY-16, 24, 16, ObstacleSpike))
	l.Obstacles = append(l.Obstacles, NewObstacle(3460, GroundY-16, 24, 16, ObstacleSpike))

	l.Obstacles = append(l.Obstacles, NewMovingObstacle(1600, GroundY-130, 80, 16, ObstacleMovingPlatform, 70, 0.018, "x"))
	l.Obstacles = append(l.Obstacles, NewMovingObstacle(3000, GroundY-140, 80, 16, ObstacleMovingPlatform, 60, 0.022, "y"))

	l.PowerUps = append(l.PowerUps, NewPowerUp(450, GroundY-80, PowerUpCarrot))
	l.PowerUps = append(l.PowerUps, NewPowerUp(750, GroundY-90, PowerUpCarrot))
	l.PowerUps = append(l.PowerUps, NewPowerUp(1200, GroundY-130, PowerUpSpeedBoost))
	l.PowerUps = append(l.PowerUps, NewPowerUp(1600, GroundY-130, PowerUpCarrot))
	l.PowerUps = append(l.PowerUps, NewPowerUp(2000, GroundY-90, PowerUpCarrot))
	l.PowerUps = append(l.PowerUps, NewPowerUp(2400, GroundY-130, PowerUpExtraLife))
	l.PowerUps = append(l.PowerUps, NewPowerUp(2800, GroundY-90, PowerUpCarrot))
	l.PowerUps = append(l.PowerUps, NewPowerUp(3200, GroundY-130, PowerUpSpeedBoost))
	l.PowerUps = append(l.PowerUps, NewPowerUp(3600, GroundY-90, PowerUpCarrot))
	l.PowerUps = append(l.PowerUps, NewPowerUp(4200, GroundY-90, PowerUpCarrot))
	l.PowerUps = append(l.PowerUps, NewPowerUp(1700, GroundY-150, PowerUpShield))

	l.Checkpoints = append(l.Checkpoints, Checkpoint{X: 1200, Y: GroundY - BunnyHeight, Active: true})
	l.Checkpoints = append(l.Checkpoints, Checkpoint{X: 2500, Y: GroundY - BunnyHeight, Active: true})
	l.Checkpoints = append(l.Checkpoints, Checkpoint{X: 3800, Y: GroundY - BunnyHeight, Active: true})
}

func (l *Level) initLevel3() {
	l.Width = 6000
	l.SpawnX = 50
	l.SpawnY = GroundY - BunnyHeight
	l.FinishX = l.Width - 100

	l.Obstacles = append(l.Obstacles, NewObstacle(0, GroundY, l.Width, 60, ObstaclePlatform))

	for i := 0; i < 15; i++ {
		x := 300 + i*400
		l.Obstacles = append(l.Obstacles, NewObstacle(float64(x), GroundY-80, 40, 80, ObstaclePipe))
	}

	platformPositions := []float64{400, 750, 1100, 1500, 1900, 2300, 2700, 3100, 3500, 3900, 4300, 4700, 5100, 5500}
	for _, px := range platformPositions {
		h := 60.0
		if int(px)%2 == 0 {
			h = 110.0
		}
		l.Obstacles = append(l.Obstacles, NewObstacle(px, GroundY-h, 50, 16, ObstaclePlatform))
	}

	l.Obstacles = append(l.Obstacles, NewObstacle(1200, GroundY-100, 60, 80, ObstacleTunnel))
	l.Obstacles = append(l.Obstacles, NewObstacle(2400, GroundY-100, 60, 80, ObstacleTunnel))
	l.Obstacles = append(l.Obstacles, NewObstacle(3600, GroundY-100, 60, 80, ObstacleTunnel))
	l.Obstacles = append(l.Obstacles, NewObstacle(4800, GroundY-100, 60, 80, ObstacleTunnel))

	spikeGroups := []float64{800, 1600, 2000, 2800, 3200, 4000, 4400, 5200, 5600}
	for _, sx := range spikeGroups {
		l.Obstacles = append(l.Obstacles, NewObstacle(sx, GroundY-16, 24, 16, ObstacleSpike))
		l.Obstacles = append(l.Obstacles, NewObstacle(sx+25, GroundY-16, 24, 16, ObstacleSpike))
		l.Obstacles = append(l.Obstacles, NewObstacle(sx+50, GroundY-16, 24, 16, ObstacleSpike))
	}

	l.Obstacles = append(l.Obstacles, NewMovingObstacle(1000, GroundY-130, 80, 16, ObstacleMovingPlatform, 80, 0.02, "x"))
	l.Obstacles = append(l.Obstacles, NewMovingObstacle(2000, GroundY-140, 80, 16, ObstacleMovingPlatform, 70, 0.025, "y"))
	l.Obstacles = append(l.Obstacles, NewMovingObstacle(3000, GroundY-130, 80, 16, ObstacleMovingPlatform, 90, 0.018, "x"))
	l.Obstacles = append(l.Obstacles, NewMovingObstacle(4000, GroundY-150, 80, 16, ObstacleMovingPlatform, 60, 0.03, "y"))
	l.Obstacles = append(l.Obstacles, NewMovingObstacle(5000, GroundY-130, 80, 16, ObstacleMovingPlatform, 80, 0.022, "x"))

	l.PowerUps = append(l.PowerUps, NewPowerUp(450, GroundY-80, PowerUpCarrot))
	l.PowerUps = append(l.PowerUps, NewPowerUp(800, GroundY-90, PowerUpCarrot))
	l.PowerUps = append(l.PowerUps, NewPowerUp(1250, GroundY-130, PowerUpSpeedBoost))
	l.PowerUps = append(l.PowerUps, NewPowerUp(1650, GroundY-130, PowerUpCarrot))
	l.PowerUps = append(l.PowerUps, NewPowerUp(2050, GroundY-90, PowerUpCarrot))
	l.PowerUps = append(l.PowerUps, NewPowerUp(2450, GroundY-130, PowerUpExtraLife))
	l.PowerUps = append(l.PowerUps, NewPowerUp(2850, GroundY-90, PowerUpCarrot))
	l.PowerUps = append(l.PowerUps, NewPowerUp(3250, GroundY-130, PowerUpSpeedBoost))
	l.PowerUps = append(l.PowerUps, NewPowerUp(3650, GroundY-90, PowerUpCarrot))
	l.PowerUps = append(l.PowerUps, NewPowerUp(4050, GroundY-130, PowerUpCarrot))
	l.PowerUps = append(l.PowerUps, NewPowerUp(4450, GroundY-90, PowerUpCarrot))
	l.PowerUps = append(l.PowerUps, NewPowerUp(4850, GroundY-130, PowerUpExtraLife))
	l.PowerUps = append(l.PowerUps, NewPowerUp(5250, GroundY-90, PowerUpCarrot))
	l.PowerUps = append(l.PowerUps, NewPowerUp(5600, GroundY-90, PowerUpCarrot))
	l.PowerUps = append(l.PowerUps, NewPowerUp(1500, GroundY-150, PowerUpShield))
	l.PowerUps = append(l.PowerUps, NewPowerUp(3500, GroundY-150, PowerUpDoubleJump))

	l.Checkpoints = append(l.Checkpoints, Checkpoint{X: 1000, Y: GroundY - BunnyHeight, Active: true})
	l.Checkpoints = append(l.Checkpoints, Checkpoint{X: 2500, Y: GroundY - BunnyHeight, Active: true})
	l.Checkpoints = append(l.Checkpoints, Checkpoint{X: 4000, Y: GroundY - BunnyHeight, Active: true})
	l.Checkpoints = append(l.Checkpoints, Checkpoint{X: 5000, Y: GroundY - BunnyHeight, Active: true})
}

func (l *Level) Update(bunny *Bunny) {
	targetCameraX := bunny.X - ScreenWidth/3
	if targetCameraX < 0 {
		targetCameraX = 0
	}
	maxCameraX := l.Width - ScreenWidth
	if maxCameraX < 0 {
		maxCameraX = 0
	}
	if targetCameraX > maxCameraX {
		targetCameraX = maxCameraX
	}
	l.CameraX += (targetCameraX - l.CameraX) * 0.1
	if l.CameraX > maxCameraX {
		l.CameraX = maxCameraX
	}
	if l.CameraX < 0 {
		l.CameraX = 0
	}

	for _, obs := range l.Obstacles {
		obs.Update()
	}

	for _, pu := range l.PowerUps {
		pu.Update()
	}

	for _, cp := range l.Checkpoints {
		if cp.Active && bunny.X >= cp.X-10 && bunny.X <= cp.X+10 {
			cp.Active = false
			l.respawnX = cp.X
			l.respawnY = cp.Y
		}
	}
}

func (l *Level) Draw(screen *ebiten.Image) {
	DrawBackground(screen, l.BackgroundLayers, l.CameraX)

	for _, obs := range l.Obstacles {
		obs.Draw(screen, l.CameraX)
	}

	for _, pu := range l.PowerUps {
		pu.Draw(screen, l.CameraX)
	}

	for _, cp := range l.Checkpoints {
		if cp.Active {
			cpOp := &ebiten.DrawImageOptions{}
			cpOp.GeoM.Translate(cp.X-l.CameraX, cp.Y)
			cpImg := CreateCheckpointImage()
			screen.DrawImage(cpImg, cpOp)
		}
	}

	finishOp := &ebiten.DrawImageOptions{}
	finishOp.GeoM.Translate(l.FinishX-l.CameraX, GroundY-60)
	finishImg := CreateFinishLineImage()
	screen.DrawImage(finishImg, finishOp)

	ebitenutil.DrawRect(screen, l.FinishX-l.CameraX+10, GroundY-100, 2, 40, ColorRed)
	ebitenutil.DrawRect(screen, l.FinishX-l.CameraX+10, GroundY-100, 12, 12, ColorFinishLine)
}

func (l *Level) IsFinished(bunny *Bunny) bool {
	return bunny.X+bunny.Width >= l.FinishX
}
