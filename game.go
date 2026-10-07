package main

import (
	"fmt"
	"math"
	"math/rand/v2"
	"strconv"

	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type GameState int

const (
	StateMenu GameState = iota
	StatePlaying
	StatePaused
	StateLevelComplete
	StateGameOver
	StateSettings
)

type Game struct {
	State             GameState
	CurrentLevel      int
	Score             int
	Lives             int
	LevelTimer        float64
	Bunny             *Bunny
	Level             *Level
	BackgroundLayers  [3]*ebiten.Image
	StateTimer        float64
	SoundEnabled      bool
	PreviousState     GameState
	ShakeTimer        float64
	ShakeMag          float64
	HighScore         int
	Particles         []*Particle
	CheckpointX       float64
	CheckpointY       float64
	ComboCount        int
	ComboTimer        float64
	SettingsSelection int
	prevPKey          bool
	prevMKey          bool
}

type Particle struct {
	X, Y    float64
	VX, VY  float64
	Life    float64
	MaxLife float64
	Color   color.Color
	Size    float64
}

func NewGame() *Game {
	g := &Game{
		State:        StateMenu,
		Lives:        MaxLives,
		Score:        0,
		SoundEnabled: true,
		HighScore:    LoadHighScore(),
	}
	g.BackgroundLayers = GenerateBackgroundLayers()
	return g
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return fmt.Errorf("quit")
	}

	pKey := ebiten.IsKeyPressed(ebiten.KeyP)
	if pKey && !g.prevPKey && g.State == StatePlaying {
		g.State = StatePaused
		g.PreviousState = StatePlaying
	} else if pKey && !g.prevPKey && g.State == StatePaused {
		g.State = g.PreviousState
	}
	g.prevPKey = pKey

	mKey := ebiten.IsKeyPressed(ebiten.KeyM)
	if mKey && !g.prevMKey {
		g.SoundEnabled = !g.SoundEnabled
		SetSoundEnabled(g.SoundEnabled)
	}
	g.prevMKey = mKey

	if g.ShakeTimer > 0 {
		g.ShakeTimer -= DeltaTime
		if g.ShakeTimer < 0 {
			g.ShakeTimer = 0
		}
	}

	if g.ComboTimer > 0 {
		g.ComboTimer -= DeltaTime
		if g.ComboTimer < 0 {
			g.ComboTimer = 0
			g.ComboCount = 0
		}
	}

	for i := len(g.Particles) - 1; i >= 0; i-- {
		p := g.Particles[i]
		p.X += p.VX
		p.Y += p.VY
		p.VY += 0.2
		p.Life -= DeltaTime
		if p.Life <= 0 {
			g.Particles = append(g.Particles[:i], g.Particles[i+1:]...)
		}
	}

	switch g.State {
	case StateMenu:
		g.updateMenu()
	case StatePlaying:
		g.updatePlaying()
	case StatePaused:
		g.updatePaused()
	case StateLevelComplete:
		g.updateLevelComplete()
	case StateGameOver:
		g.updateGameOver()
	case StateSettings:
		g.updateSettings()
	}
	return nil
}

func (g *Game) updateMenu() {
	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		g.startGame()
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.State = StateSettings
		g.SettingsSelection = 0
		g.PreviousState = StateMenu
	}
	if ebiten.IsGamepadButtonPressed(0, ebiten.GamepadButton0) || ebiten.IsGamepadButtonPressed(0, ebiten.GamepadButton1) {
		g.startGame()
	}
	if ebiten.IsGamepadButtonPressed(0, ebiten.GamepadButton6) {
		g.State = StateSettings
		g.SettingsSelection = 0
		g.PreviousState = StateMenu
	}
}

func (g *Game) updatePlaying() {
	keys := map[ebiten.Key]bool{
		ebiten.KeyLeft:  ebiten.IsKeyPressed(ebiten.KeyLeft),
		ebiten.KeyRight: ebiten.IsKeyPressed(ebiten.KeyRight),
		ebiten.KeyA:     ebiten.IsKeyPressed(ebiten.KeyA),
		ebiten.KeyD:     ebiten.IsKeyPressed(ebiten.KeyD),
		ebiten.KeySpace: ebiten.IsKeyPressed(ebiten.KeySpace),
		ebiten.KeyUp:    ebiten.IsKeyPressed(ebiten.KeyUp),
		ebiten.KeyW:     ebiten.IsKeyPressed(ebiten.KeyW),
	}

	if ebiten.IsGamepadButtonPressed(0, ebiten.GamepadButton16) {
		keys[ebiten.KeyLeft] = true
	}
	if ebiten.IsGamepadButtonPressed(0, ebiten.GamepadButton17) {
		keys[ebiten.KeyRight] = true
	}
	if ebiten.IsGamepadButtonPressed(0, ebiten.GamepadButton0) || ebiten.IsGamepadButtonPressed(0, ebiten.GamepadButton1) {
		keys[ebiten.KeySpace] = true
	}
	if ebiten.IsGamepadButtonPressed(0, ebiten.GamepadButton9) {
		g.State = StatePaused
		g.PreviousState = StatePlaying
	}

	g.Bunny.Update(keys, g.Level.Width, g.Level.Obstacles)
	g.Level.Update(g.Bunny)
	g.LevelTimer += DeltaTime

	cp := g.Level.GetRespawnPoint()
	if cp.X > 0 || cp.Y > 0 {
		g.CheckpointX = cp.X
		g.CheckpointY = cp.Y
	}

	if g.Bunny.Y > ScreenHeight+50 {
		g.handleBunnyDeath()
		return
	}

	if g.checkObstacleCollision() {
		g.handleBunnyDeath()
		return
	}

	g.checkPowerUpCollision()
	if g.Level.IsFinished(g.Bunny) {
		g.completeLevel()
	}
}

func (g *Game) updatePaused() {
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.State = StateSettings
		g.SettingsSelection = 0
		g.PreviousState = StatePaused
	}
	if ebiten.IsGamepadButtonPressed(0, ebiten.GamepadButton6) {
		g.State = StateSettings
		g.SettingsSelection = 0
		g.PreviousState = StatePaused
	}
}

func (g *Game) updateSettings() {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) || ebiten.IsKeyPressed(ebiten.KeyBackspace) {
		g.State = g.PreviousState
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		g.SettingsSelection--
		if g.SettingsSelection < 0 {
			g.SettingsSelection = 1
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		g.SettingsSelection++
		if g.SettingsSelection > 1 {
			g.SettingsSelection = 0
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyEnter) || ebiten.IsKeyPressed(ebiten.KeySpace) {
		if g.SettingsSelection == 0 {
			g.SoundEnabled = !g.SoundEnabled
			SetSoundEnabled(g.SoundEnabled)
		} else if g.SettingsSelection == 1 {
			g.State = g.PreviousState
		}
	}
	if ebiten.IsGamepadButtonPressed(0, ebiten.GamepadButton6) {
		g.State = g.PreviousState
	}
	if ebiten.IsGamepadButtonPressed(0, ebiten.GamepadButton0) || ebiten.IsGamepadButtonPressed(0, ebiten.GamepadButton1) {
		if g.SettingsSelection == 0 {
			g.SoundEnabled = !g.SoundEnabled
			SetSoundEnabled(g.SoundEnabled)
		} else if g.SettingsSelection == 1 {
			g.State = g.PreviousState
		}
	}
}

func (g *Game) updateLevelComplete() {
	g.StateTimer += DeltaTime
	if ebiten.IsKeyPressed(ebiten.KeyEnter) && g.StateTimer > 1.0 {
		g.nextLevel()
	}
	if ebiten.IsGamepadButtonPressed(0, ebiten.GamepadButton0) && g.StateTimer > 1.0 {
		g.nextLevel()
	}
}

func (g *Game) updateGameOver() {
	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		g.startGame()
	}
	if ebiten.IsGamepadButtonPressed(0, ebiten.GamepadButton0) || ebiten.IsGamepadButtonPressed(0, ebiten.GamepadButton1) {
		g.startGame()
	}
}

func (g *Game) startGame() {
	g.State = StatePlaying
	g.CurrentLevel = 1
	g.Score = 0
	g.Lives = MaxLives
	g.LevelTimer = 0
	g.StateTimer = 0
	g.ShakeTimer = 0
	g.ComboCount = 0
	g.ComboTimer = 0
	g.Particles = nil
	g.loadLevel(1)
}

func (g *Game) loadLevel(levelNum int) {
	g.Level = NewLevel(levelNum)
	g.Bunny = NewBunny(g.Level.SpawnX, g.Level.SpawnY)
	g.LevelTimer = 0
	g.StateTimer = 0
	g.CheckpointX = g.Level.SpawnX
	g.CheckpointY = g.Level.SpawnY
	g.Particles = nil
}

func (g *Game) handleBunnyDeath() {
	if g.Bunny.InvincibleTimer > 0 {
		return
	}
	if g.Bunny.Shield {
		g.Bunny.UseShield()
		g.ShakeTimer = ScreenShakeTime
		g.ShakeMag = ScreenShakeMag * 0.5
		PlayShieldSound()
		g.spawnParticles(g.Bunny.X+g.Bunny.Width/2, g.Bunny.Y+g.Bunny.Height/2, 20, ColorShield)
		return
	}
	PlayHitSound()
	g.ShakeTimer = ScreenShakeTime
	g.ShakeMag = ScreenShakeMag
	g.spawnParticles(g.Bunny.X+g.Bunny.Width/2, g.Bunny.Y+g.Bunny.Height/2, 30, ColorRed)
	g.Bunny.Die()
	g.Lives = g.Bunny.Lives
	if g.Lives <= 0 {
		g.State = StateGameOver
		g.StateTimer = 0
		PlayGameOverSound()
		if g.Score > g.HighScore {
			g.HighScore = g.Score
			SaveHighScore(g.HighScore)
		}
	} else {
		g.Bunny.Reset(g.CheckpointX, g.CheckpointY)
	}
}

func (g *Game) checkObstacleCollision() bool {
	bunnyBounds := g.Bunny.GetBounds()
	for _, obs := range g.Level.Obstacles {
		if obs.Type == ObstaclePlatform || obs.Type == ObstacleTunnel || obs.Type == ObstacleMovingPlatform {
			continue
		}
		obsBounds := obs.GetBounds()
		if bunnyBounds.Intersects(obsBounds) {
			return true
		}
	}
	return false
}

func (g *Game) checkPowerUpCollision() {
	bunnyBounds := g.Bunny.GetBounds()
	for _, pu := range g.Level.PowerUps {
		if pu.Collected {
			continue
		}
		puBounds := pu.GetBounds()
		if bunnyBounds.Intersects(puBounds) {
			pu.Collected = true
			g.collectPowerUp(pu)
		}
	}
}

func (g *Game) collectPowerUp(pu *PowerUp) {
	PlayCollectSound()
	g.ComboCount++
	g.ComboTimer = 2.0
	comboMultiplier := 1
	if g.ComboCount >= 5 {
		comboMultiplier = 3
	} else if g.ComboCount >= 3 {
		comboMultiplier = 2
	}
	g.spawnParticles(pu.X+pu.Width/2, pu.Y+pu.Height/2, 10, ColorYellow)
	switch pu.Type {
	case PowerUpCarrot:
		g.Score += CarrotPoints * comboMultiplier
	case PowerUpSpeedBoost:
		g.Bunny.ActivateSpeedBoost()
	case PowerUpExtraLife:
		g.Bunny.AddLife()
	case PowerUpDoubleJump:
		g.Bunny.ActivateDoubleJump()
	case PowerUpShield:
		g.Bunny.ActivateShield()
	}
}

func (g *Game) spawnParticles(x, y float64, count int, c color.Color) {
	for i := 0; i < count; i++ {
		angle := math.Pi*2*float64(i)/float64(count) + rand.Float64()*0.5
		speed := 2 + rand.Float64()*3
		p := &Particle{
			X:       x,
			Y:       y,
			VX:      math.Cos(angle) * speed,
			VY:      math.Sin(angle)*speed - 2,
			Life:    0.5 + rand.Float64()*0.5,
			MaxLife: 1.0,
			Color:   c,
			Size:    3 + rand.Float64()*3,
		}
		g.Particles = append(g.Particles, p)
	}
}

func (g *Game) completeLevel() {
	PlayLevelCompleteSound()
	g.State = StateLevelComplete
	g.StateTimer = 0
	timeTaken := g.LevelTimer
	basePoints := g.CurrentLevel * BasePointsPerLevel
	timeBonus := int(1000.0 / math.Max(timeTaken, 1.0))
	comboBonus := g.ComboCount * 5
	totalScore := basePoints + g.Score + timeBonus + comboBonus
	g.Score = totalScore
	if g.Score > g.HighScore {
		g.HighScore = g.Score
		SaveHighScore(g.HighScore)
	}
}

func (g *Game) nextLevel() {
	g.CurrentLevel++
	if g.CurrentLevel > 3 {
		g.CurrentLevel = 1
	}
	g.loadLevel(g.CurrentLevel)
	g.State = StatePlaying
}

func (g *Game) Draw(screen *ebiten.Image) {
	shakeX := 0.0
	shakeY := 0.0
	if g.ShakeTimer > 0 {
		shakeX = (rand.Float64() - 0.5) * g.ShakeMag * 2
		shakeY = (rand.Float64() - 0.5) * g.ShakeMag * 2
	}

	switch g.State {
	case StateMenu:
		g.drawMenu(screen)
	case StatePlaying:
		g.drawPlaying(screen, shakeX, shakeY)
	case StatePaused:
		g.drawPlaying(screen, shakeX, shakeY)
		g.drawPausedOverlay(screen)
	case StateLevelComplete:
		g.drawLevelComplete(screen, shakeX, shakeY)
	case StateGameOver:
		g.drawGameOver(screen)
	case StateSettings:
		g.drawSettings(screen)
	}

	for _, p := range g.Particles {
		ebitenutil.DrawRect(screen, p.X, p.Y, p.Size, p.Size, p.Color)
	}
}

func (g *Game) drawPausedOverlay(screen *ebiten.Image) {
	ebitenutil.DebugPrintAt(screen, "PAUSED", ScreenWidth/2-30, ScreenHeight/2-30)
	ebitenutil.DebugPrintAt(screen, "Press P to resume", ScreenWidth/2-60, ScreenHeight/2)
}

func (g *Game) drawMenu(screen *ebiten.Image) {
	DrawBackground(screen, g.BackgroundLayers, 0)
	DrawTextCenter(screen, "BUNNY OBBY", 150)
	DrawTextCenter(screen, "================", 180)
	DrawTextCenter(screen, "Arrow Keys / A-D: Move", 250)
	DrawTextCenter(screen, "Space / W / Up: Jump", 270)
	DrawTextCenter(screen, "Double Jump: Unlockable!", 290)
	DrawTextCenter(screen, "Escape: Quit", 310)
	DrawTextCenter(screen, "Press ENTER to Start", 370)
	DrawTextCenter(screen, "P: Pause  S: Settings", 400)
	DrawTextCenter(screen, "M: Toggle Sound", 420)
	DrawTextCenter(screen, "High Score: "+strconv.Itoa(g.HighScore), 450)
}

func (g *Game) drawPlaying(screen *ebiten.Image, shakeX, shakeY float64) {
	g.Level.Draw(screen)
	g.Bunny.Draw(screen, g.Level.CameraX)
	g.drawHUD(screen)
}

func (g *Game) drawHUD(screen *ebiten.Image) {
	ebitenutil.DebugPrintAt(screen, "Score: "+strconv.Itoa(g.Score), 10, 10)
	livesText := "Lives: " + strconv.Itoa(g.Lives)
	ebitenutil.DebugPrintAt(screen, livesText, 10, 30)
	ebitenutil.DebugPrintAt(screen, "Level: "+strconv.Itoa(g.CurrentLevel), 10, 50)
	if g.Bunny.SpeedBoostTimer > 0 {
		ebitenutil.DebugPrintAt(screen, "SPEED BOOST!", 10, 70)
	}
	if g.Bunny.DoubleJump && !g.Bunny.DoubleJumpUsed {
		ebitenutil.DebugPrintAt(screen, "DOUBLE JUMP READY", 10, 90)
	} else if g.Bunny.DoubleJump {
		ebitenutil.DebugPrintAt(screen, "Double Jump: Used", 10, 90)
	}
	if g.Bunny.Shield {
		ebitenutil.DebugPrintAt(screen, "SHIELD ACTIVE", 10, 110)
	}
	if g.ComboCount >= 3 {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("COMBO x%d!", g.ComboCount), 10, 130)
	}

	soundText := "Sound: ON"
	if !g.SoundEnabled {
		soundText = "Sound: OFF"
	}
	ebitenutil.DebugPrintAt(screen, soundText, 10, 150)
}

func (g *Game) drawSettings(screen *ebiten.Image) {
	DrawBackground(screen, g.BackgroundLayers, 0)
	DrawTextCenter(screen, "SETTINGS", 150)
	DrawTextCenter(screen, "================", 180)

	soundText := "Sound: ON"
	if !g.SoundEnabled {
		soundText = "Sound: OFF"
	}
	if g.SettingsSelection == 0 {
		ebitenutil.DebugPrintAt(screen, "> "+soundText, ScreenWidth/2-60, 250)
	} else {
		ebitenutil.DebugPrintAt(screen, "  "+soundText, ScreenWidth/2-60, 250)
	}

	if g.SettingsSelection == 1 {
		ebitenutil.DebugPrintAt(screen, "> Back", ScreenWidth/2-30, 290)
	} else {
		ebitenutil.DebugPrintAt(screen, "  Back", ScreenWidth/2-30, 290)
	}

	DrawTextCenter(screen, "Up/Down: Navigate", 350)
	DrawTextCenter(screen, "Enter: Select", 370)
	DrawTextCenter(screen, "Esc: Back", 390)
}

func (g *Game) drawLevelComplete(screen *ebiten.Image, shakeX, shakeY float64) {
	g.Level.Draw(screen)
	g.Bunny.Draw(screen, g.Level.CameraX)
	ebitenutil.DebugPrintAt(screen, "LEVEL COMPLETE!", ScreenWidth/2-80, ScreenHeight/2-60)
	ebitenutil.DebugPrintAt(screen, "Score: "+strconv.Itoa(g.Score), ScreenWidth/2-60, ScreenHeight/2-20)
	ebitenutil.DebugPrintAt(screen, "Time: "+fmt.Sprintf("%.1fs", g.LevelTimer), ScreenWidth/2-50, ScreenHeight/2+10)
	ebitenutil.DebugPrintAt(screen, "Press ENTER for next level", ScreenWidth/2-100, ScreenHeight/2+40)
}

func (g *Game) drawGameOver(screen *ebiten.Image) {
	DrawBackground(screen, g.BackgroundLayers, 0)
	ebitenutil.DebugPrintAt(screen, "GAME OVER", ScreenWidth/2-60, ScreenHeight/2-60)
	ebitenutil.DebugPrintAt(screen, "Final Score: "+strconv.Itoa(g.Score), ScreenWidth/2-70, ScreenHeight/2-20)
	ebitenutil.DebugPrintAt(screen, "High Score: "+strconv.Itoa(g.HighScore), ScreenWidth/2-70, ScreenHeight/2+10)
	ebitenutil.DebugPrintAt(screen, "Press ENTER to Restart", ScreenWidth/2-80, ScreenHeight/2+50)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}
