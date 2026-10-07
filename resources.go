package main

import (
	"bytes"
	"fmt"
	"image/color"
	"math"
	"os"
	"path/filepath"
	"runtime"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	ScreenWidth  = 800
	ScreenHeight = 600
)

const (
	AudioSampleRate  = 44100
	BeepFreqJump     = 440.0
	BeepFreqCollect  = 880.0
	BeepFreqHit      = 150.0
	BeepFreqComplete = 660.0
	BeepDuration     = 0.1
)

var (
	ColorSkyTop      = color.RGBA{135, 206, 235, 255}
	ColorSkyBottom   = color.RGBA{200, 230, 255, 255}
	ColorGrass       = color.RGBA{34, 139, 34, 255}
	ColorGrassDark   = color.RGBA{20, 100, 20, 255}
	ColorGround      = color.RGBA{139, 90, 43, 255}
	ColorPipeGreen   = color.RGBA{0, 180, 0, 255}
	ColorPipeDark    = color.RGBA{0, 120, 0, 255}
	ColorPipeCap     = color.RGBA{0, 220, 0, 255}
	ColorTunnelBrown = color.RGBA{139, 90, 43, 255}
	ColorTunnelDark  = color.RGBA{80, 50, 20, 255}
	ColorPlatform    = color.RGBA{160, 120, 60, 255}
	ColorPlatformTop = color.RGBA{34, 139, 34, 255}
	ColorCarrot      = color.RGBA{255, 140, 0, 255}
	ColorCarrotTop   = color.RGBA{0, 180, 0, 255}
	ColorSpeedBoost  = color.RGBA{255, 255, 0, 255}
	ColorExtraLife   = color.RGBA{255, 50, 50, 255}
	ColorBunnyBody   = color.RGBA{240, 240, 240, 255}
	ColorBunnyEar    = color.RGBA{255, 180, 180, 255}
	ColorBunnyEye    = color.RGBA{0, 0, 0, 255}
	ColorBunnyNose   = color.RGBA{255, 150, 150, 255}
	ColorSpike       = color.RGBA{100, 100, 100, 255}
	ColorFinishLine  = color.RGBA{255, 215, 0, 255}
	ColorWhite       = color.White
	ColorBlack       = color.Black
	ColorRed         = color.RGBA{255, 0, 0, 255}
	ColorBlue        = color.RGBA{0, 0, 255, 255}
	ColorYellow      = color.RGBA{255, 255, 0, 255}
	ColorMountain1   = color.RGBA{80, 140, 80, 255}
	ColorMountain2   = color.RGBA{60, 110, 60, 255}
	ColorCloud       = color.RGBA{255, 255, 255, 200}
	ColorGroundDark  = color.RGBA{80, 50, 20, 255}
	ColorDoubleJump  = color.RGBA{0, 200, 255, 255}
	ColorShield      = color.RGBA{100, 100, 255, 255}
	ColorParticle    = color.RGBA{255, 200, 100, 255}
)

var audioCtx *audio.Context
var soundEnabled = true

func InitAudio() {
	audioCtx = audio.NewContext(AudioSampleRate)
}

func generateBeepSamples(freq float64, duration float64) []byte {
	numSamples := int(duration * float64(AudioSampleRate))
	buf := make([]byte, numSamples*2)
	for i := 0; i < numSamples; i++ {
		t := float64(i) / float64(AudioSampleRate)
		env := 1.0
		fadeLen := int(0.01 * float64(AudioSampleRate))
		if i < fadeLen {
			env = float64(i) / float64(fadeLen)
		} else if i > numSamples-fadeLen {
			env = float64(numSamples-i) / float64(fadeLen)
		}
		sample := math.Sin(2*math.Pi*freq*t) * 0.3 * env
		sampleInt := int16(sample * 32767)
		buf[i*2] = byte(sampleInt)
		buf[i*2+1] = byte(sampleInt >> 8)
	}
	return buf
}

func playBeep(freq float64, duration float64) {
	if audioCtx == nil || !soundEnabled {
		return
	}
	samples := generateBeepSamples(freq, duration)
	player, err := audioCtx.NewPlayer(bytes.NewReader(samples))
	if err != nil {
		return
	}
	player.Rewind()
	player.SetVolume(0.3)
	player.Play()
}

func SetSoundEnabled(enabled bool) {
	soundEnabled = enabled
}

func IsSoundEnabled() bool {
	return soundEnabled
}

func PlayJumpSound() {
	playBeep(BeepFreqJump, BeepDuration)
}

func PlayCollectSound() {
	playBeep(BeepFreqCollect, BeepDuration*0.5)
}

func PlayHitSound() {
	playBeep(BeepFreqHit, BeepDuration*2)
}

func PlayLevelCompleteSound() {
	playBeep(BeepFreqComplete, BeepDuration)
	go func() {
		playBeep(BeepFreqComplete*1.5, BeepDuration)
	}()
}

func PlayGameOverSound() {
	playBeep(BeepFreqHit, BeepDuration*3)
}

func PlayDoubleJumpSound() {
	playBeep(BeepFreqJump*1.5, BeepDuration*0.8)
}

func PlayShieldSound() {
	playBeep(BeepFreqCollect*2, BeepDuration*0.3)
	go func() {
		playBeep(BeepFreqCollect*2.5, BeepDuration*0.3)
	}()
}

func LoadImage(filename string) (*ebiten.Image, error) {
	img, _, err := ebitenutil.NewImageFromFile(filename)
	if err == nil {
		return img, nil
	}
	return nil, err
}

func CreateBunnyImage() *ebiten.Image {
	return CreateBunnyImageFrame(0)
}

func CreateBunnyImageFrame(frame int) *ebiten.Image {
	img := ebiten.NewImage(32, 32)
	ebitenutil.DrawRect(img, 4, 8, 24, 18, ColorBunnyBody)
	ebitenutil.DrawRect(img, 6, 4, 20, 12, ColorBunnyBody)
	ebitenutil.DrawRect(img, 8, 0, 16, 10, ColorBunnyBody)
	ebitenutil.DrawRect(img, 10, 0, 5, 10, ColorBunnyEar)
	ebitenutil.DrawRect(img, 17, 0, 5, 10, ColorBunnyEar)
	ebitenutil.DrawRect(img, 11, 0, 3, 6, color.RGBA{255, 200, 200, 255})
	ebitenutil.DrawRect(img, 18, 0, 3, 6, color.RGBA{255, 200, 200, 255})
	ebitenutil.DrawRect(img, 10, 4, 3, 3, ColorBunnyEye)
	ebitenutil.DrawRect(img, 19, 4, 3, 3, ColorBunnyEye)
	ebitenutil.DrawRect(img, 14, 8, 4, 2, ColorBunnyNose)
	ebitenutil.DrawRect(img, 8, 7, 2, 2, color.RGBA{255, 180, 180, 255})
	ebitenutil.DrawRect(img, 22, 7, 2, 2, color.RGBA{255, 180, 180, 255})
	switch frame {
	case 0:
		ebitenutil.DrawRect(img, 6, 26, 6, 6, ColorBunnyBody)
		ebitenutil.DrawRect(img, 20, 26, 6, 6, ColorBunnyBody)
	case 1:
		ebitenutil.DrawRect(img, 4, 26, 6, 6, ColorBunnyBody)
		ebitenutil.DrawRect(img, 22, 26, 6, 6, ColorBunnyBody)
	case 2:
		ebitenutil.DrawRect(img, 6, 26, 6, 6, ColorBunnyBody)
		ebitenutil.DrawRect(img, 20, 26, 6, 6, ColorBunnyBody)
	case 3:
		ebitenutil.DrawRect(img, 8, 26, 6, 6, ColorBunnyBody)
		ebitenutil.DrawRect(img, 18, 26, 6, 6, ColorBunnyBody)
	}
	ebitenutil.DrawRect(img, 2, 16, 4, 4, ColorWhite)
	return img
}

func CreateObstacleImage(obsType ObstacleType) *ebiten.Image {
	var img *ebiten.Image
	switch obsType {
	case ObstaclePipe:
		img = ebiten.NewImage(40, 80)
		ebitenutil.DrawRect(img, 4, 16, 32, 64, ColorPipeGreen)
		ebitenutil.DrawRect(img, 0, 8, 40, 12, ColorPipeCap)
		ebitenutil.DrawRect(img, 4, 16, 32, 4, ColorPipeDark)
	case ObstacleTunnel:
		img = ebiten.NewImage(60, 80)
		ebitenutil.DrawRect(img, 0, 0, 60, 20, ColorTunnelBrown)
		ebitenutil.DrawRect(img, 0, 60, 60, 20, ColorTunnelBrown)
		ebitenutil.DrawRect(img, 0, 20, 8, 40, ColorTunnelBrown)
		ebitenutil.DrawRect(img, 52, 20, 8, 40, ColorTunnelBrown)
		ebitenutil.DrawRect(img, 0, 0, 60, 1, ColorTunnelDark)
		ebitenutil.DrawRect(img, 0, 79, 60, 1, ColorTunnelDark)
	case ObstaclePlatform:
		img = ebiten.NewImage(80, 16)
		ebitenutil.DrawRect(img, 0, 4, 80, 12, ColorPlatform)
		ebitenutil.DrawRect(img, 0, 0, 80, 4, ColorPlatformTop)
		ebitenutil.DrawRect(img, 0, 14, 80, 2, ColorGroundDark)
	case ObstacleSpike:
		img = ebiten.NewImage(24, 20)
		ebitenutil.DrawRect(img, 0, 16, 24, 4, ColorSpike)
		ebitenutil.DrawRect(img, 2, 10, 6, 6, ColorSpike)
		ebitenutil.DrawRect(img, 9, 4, 6, 12, ColorSpike)
		ebitenutil.DrawRect(img, 16, 10, 6, 6, ColorSpike)
	}
	if img == nil {
		img = ebiten.NewImage(1, 1)
	}
	return img
}

func CreatePowerUpImage(puType PowerUpType) *ebiten.Image {
	img := ebiten.NewImage(20, 20)
	switch puType {
	case PowerUpCarrot:
		ebitenutil.DrawRect(img, 8, 4, 4, 12, ColorCarrot)
		ebitenutil.DrawRect(img, 6, 6, 8, 4, ColorCarrot)
		ebitenutil.DrawRect(img, 4, 8, 12, 2, ColorCarrot)
		ebitenutil.DrawRect(img, 8, 0, 4, 6, ColorCarrotTop)
		ebitenutil.DrawRect(img, 9, 4, 2, 1, ColorWhite)
		ebitenutil.DrawRect(img, 9, 8, 2, 1, ColorWhite)
	case PowerUpSpeedBoost:
		ebitenutil.DrawRect(img, 8, 0, 4, 10, ColorSpeedBoost)
		ebitenutil.DrawRect(img, 4, 6, 4, 6, ColorSpeedBoost)
		ebitenutil.DrawRect(img, 10, 6, 4, 6, ColorSpeedBoost)
		ebitenutil.DrawRect(img, 6, 10, 8, 4, ColorSpeedBoost)
		ebitenutil.DrawRect(img, 7, 1, 6, 2, ColorWhite)
	case PowerUpExtraLife:
		ebitenutil.DrawRect(img, 4, 6, 4, 4, ColorExtraLife)
		ebitenutil.DrawRect(img, 12, 6, 4, 4, ColorExtraLife)
		ebitenutil.DrawRect(img, 6, 4, 8, 4, ColorExtraLife)
		ebitenutil.DrawRect(img, 8, 8, 4, 8, ColorExtraLife)
		ebitenutil.DrawRect(img, 6, 12, 8, 4, ColorExtraLife)
	case PowerUpDoubleJump:
		ebitenutil.DrawRect(img, 6, 2, 8, 16, ColorDoubleJump)
		ebitenutil.DrawRect(img, 4, 6, 12, 8, ColorDoubleJump)
		ebitenutil.DrawRect(img, 8, 0, 4, 4, ColorWhite)
		ebitenutil.DrawRect(img, 2, 10, 4, 4, ColorWhite)
		ebitenutil.DrawRect(img, 14, 10, 4, 4, ColorWhite)
	case PowerUpShield:
		ebitenutil.DrawRect(img, 4, 4, 12, 12, ColorShield)
		ebitenutil.DrawRect(img, 6, 2, 8, 16, ColorShield)
		ebitenutil.DrawRect(img, 2, 6, 16, 8, ColorShield)
		ebitenutil.DrawRect(img, 8, 6, 4, 8, ColorWhite)
	}
	return img
}

func CreateSkyImage() *ebiten.Image {
	img := ebiten.NewImage(ScreenWidth, ScreenHeight)
	for y := 0; y < ScreenHeight; y++ {
		t := float64(y) / float64(ScreenHeight)
		r := uint8(float64(ColorSkyTop.R) + t*(float64(ColorSkyBottom.R)-float64(ColorSkyTop.R)))
		g := uint8(float64(ColorSkyTop.G) + t*(float64(ColorSkyBottom.G)-float64(ColorSkyTop.G)))
		b := uint8(float64(ColorSkyTop.B) + t*(float64(ColorSkyBottom.B)-float64(ColorSkyTop.B)))
		cy := float64(y)
		ebitenutil.DrawRect(img, 0, cy, ScreenWidth, 1, color.RGBA{r, g, b, 255})
	}
	return img
}

func CreateMountainImage() *ebiten.Image {
	img := ebiten.NewImage(ScreenWidth*2, ScreenHeight)
	mountains := []struct{ x, w, h int }{
		{0, 200, 120},
		{150, 180, 90},
		{300, 220, 140},
		{500, 160, 80},
		{650, 200, 110},
		{800, 180, 95},
		{950, 220, 130},
	}
	for _, m := range mountains {
		peakX := m.x + m.w/2
		baseY := ScreenHeight
		for sx := m.x; sx < peakX; sx++ {
			ratio := float64(sx-m.x) / float64(peakX-m.x)
			sy := int(float64(baseY) - ratio*float64(m.h))
			ebitenutil.DrawRect(img, float64(sx), float64(sy), 1, float64(baseY-sy+1), ColorMountain1)
		}
		for sx := peakX; sx < m.x+m.w; sx++ {
			ratio := float64(sx-peakX) / float64(m.x+m.w-peakX)
			sy := int(float64(baseY) - ratio*float64(m.h))
			ebitenutil.DrawRect(img, float64(sx), float64(sy), 1, float64(baseY-sy+1), ColorMountain1)
		}
	}
	return img
}

func CreateGroundImage() *ebiten.Image {
	img := ebiten.NewImage(ScreenWidth*2, 60)
	ebitenutil.DrawRect(img, 0, 0, float64(ScreenWidth*2), 8, ColorGrass)
	ebitenutil.DrawRect(img, 0, 8, float64(ScreenWidth*2), 52, ColorGround)
	for i := 0; i < ScreenWidth*2; i += 12 {
		ebitenutil.DrawRect(img, float64(i), 0, 2, 6, ColorGrassDark)
	}
	return img
}

func CreateCloudImage() *ebiten.Image {
	img := ebiten.NewImage(80, 30)
	ebitenutil.DrawRect(img, 10, 8, 60, 14, ColorCloud)
	ebitenutil.DrawRect(img, 20, 4, 30, 10, ColorCloud)
	ebitenutil.DrawRect(img, 40, 6, 20, 8, ColorCloud)
	return img
}

func CreateFinishLineImage() *ebiten.Image {
	img := ebiten.NewImage(20, 60)
	for y := 0; y < 60; y += 10 {
		for x := 0; x < 20; x += 10 {
			var c color.Color
			if (x/10+y/10)%2 == 0 {
				c = ColorFinishLine
			} else {
				c = ColorBlack
			}
			ebitenutil.DrawRect(img, float64(x), float64(y), 10, 10, c)
		}
	}
	return img
}

func CreateCheckpointImage() *ebiten.Image {
	img := ebiten.NewImage(20, 40)
	ebitenutil.DrawRect(img, 8, 0, 4, 40, ColorFinishLine)
	ebitenutil.DrawRect(img, 4, 0, 12, 8, ColorFinishLine)
	ebitenutil.DrawRect(img, 2, 8, 16, 4, ColorFinishLine)
	ebitenutil.DrawRect(img, 0, 12, 20, 4, ColorFinishLine)
	return img
}

func GenerateBackgroundLayers() [3]*ebiten.Image {
	return [3]*ebiten.Image{
		CreateSkyImage(),
		CreateMountainImage(),
		CreateGroundImage(),
	}
}

func DrawBackground(screen *ebiten.Image, layers [3]*ebiten.Image, cameraX float64) {
	skyOp := &ebiten.DrawImageOptions{}
	screen.DrawImage(layers[0], skyOp)

	mountainW := layers[1].Bounds().Dx()
	mountainOffset := int(-cameraX*0.3) % mountainW
	if mountainOffset > 0 {
		mountainOffset -= mountainW
	}
	for offset := mountainOffset; offset > -ScreenWidth; offset -= mountainW {
		mountainOp := &ebiten.DrawImageOptions{}
		mountainOp.GeoM.Translate(float64(offset), 0)
		screen.DrawImage(layers[1], mountainOp)
	}
	for offset := mountainOffset + mountainW; offset < ScreenWidth; offset += mountainW {
		mountainOp := &ebiten.DrawImageOptions{}
		mountainOp.GeoM.Translate(float64(offset), 0)
		screen.DrawImage(layers[1], mountainOp)
	}

	groundW := layers[2].Bounds().Dx()
	groundOffset := int(-cameraX) % groundW
	if groundOffset > 0 {
		groundOffset -= groundW
	}
	groundY := float64(ScreenHeight - 60)
	for offset := groundOffset; offset > -ScreenWidth; offset -= groundW {
		groundOp := &ebiten.DrawImageOptions{}
		groundOp.GeoM.Translate(float64(offset), groundY)
		screen.DrawImage(layers[2], groundOp)
	}
	for offset := groundOffset + groundW; offset < ScreenWidth; offset += groundW {
		groundOp := &ebiten.DrawImageOptions{}
		groundOp.GeoM.Translate(float64(offset), groundY)
		screen.DrawImage(layers[2], groundOp)
	}
}

func DrawText(screen *ebiten.Image, text string, x, y int) {
	ebitenutil.DebugPrintAt(screen, text, x, y)
}

func DrawTextCenter(screen *ebiten.Image, text string, y int) {
	w := len(text) * 6
	x := (ScreenWidth - w) / 2
	ebitenutil.DebugPrintAt(screen, text, x, y)
}

func getHighScorePath() string {
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		return filepath.Join(filepath.Dir(filename), "highscore.dat")
	}
	return "highscore.dat"
}

func LoadHighScore() int {
	data, err := os.ReadFile(getHighScorePath())
	if err != nil {
		return 0
	}
	var score int
	_, err = fmt.Sscanf(string(data), "%d", &score)
	if err != nil {
		return 0
	}
	return score
}

func SaveHighScore(score int) {
	_ = os.WriteFile(getHighScorePath(), []byte(fmt.Sprintf("%d", score)), 0644)
}
