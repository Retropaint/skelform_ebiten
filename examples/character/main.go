package main

import (
	"bytes"
	"fmt"
	"image/color"
	"log"
	"math"
	"path/filepath"
	"runtime"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	// "github.com/hajimehoshi/ebiten/v2/ebitenutil"
	skf_e "github.com/retropaint/skelform_ebiten"
	skf "github.com/retropaint/skelform_go"
)

var (
	mplusFaceSource *text.GoTextFaceSource
)

type Game struct {
	skellington skf.Armature
	skelTex     []*ebiten.Image
	skellina    skf.Armature
	skelaTex    []*ebiten.Image
	playerPos   skf.Vec2
	moving      bool
	skelTime    time.Time
	skelaTime   time.Time
	dir         float32
	style       int
	velocityY   float32
	lastAnimidx int
	groundY     float32
}

// ensures .skf files can be loaded from `go run`,
// not needed in actual projects  
func assetPath(name string) string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), name)
}

func (g *Game) Update() error {
	var speed float32 = 7
	g.moving = false

	var keys []ebiten.Key
	for _, key := range inpututil.AppendJustPressedKeys(keys) {
		switch key {
		case ebiten.Key1:
			g.style = 1
		case ebiten.Key2:
			g.style = 0
		case ebiten.KeySpace:
			g.velocityY = -10
		}
	}
	for _, key := range inpututil.AppendJustReleasedKeys(keys) {
		switch key {
		case ebiten.KeyA:
			g.skelTime = time.Now()
		case ebiten.KeyD:
			g.skelTime = time.Now()
		}
	}
	for _, key := range inpututil.AppendPressedKeys(keys) {
		switch key {
		case ebiten.KeyA:
			g.playerPos.X -= speed
			g.moving = true
			g.dir = -1
		case ebiten.KeyD:
			g.playerPos.X += speed
			g.moving = true
			g.dir = 1
		}
	}

	g.velocityY += 0.2
	g.playerPos.Y += g.velocityY
	if g.playerPos.Y > g.groundY {
		g.playerPos.Y = g.groundY
	}

	return nil
}

func (g *Game) Skellina(screen *ebiten.Image) {
	skela := &g.skellina
	tf0 := skf.TimeFrame(skela.Animations[2], time.Since(g.skelaTime), false, true)
	anims := []skf.Animation{skela.Animations[2]}
	skf_e.Animate(skela, anims, []int{tf0}, []int{20})
	pos := skf.Vec2{X: 50, Y: g.playerPos.Y + 50}
	constructOptions := skf_e.ConstructOptions{Scale: skf.Vec2{X: 0.125, Y: 0.125}, Position: pos}
	finalBones := skf_e.Construct(*skela, constructOptions)
	skf_e.Draw(finalBones, skela.Styles, g.skelaTex, screen)
}

func (g *Game) Skellington(screen *ebiten.Image) {
	constructOptions := skf_e.ConstructOptions{}
	constructOptions.Init()

	skelScale := float32(0.125)
	constructOptions.Position = g.playerPos
	constructOptions.Scale = skf.Vec2{X: skelScale, Y: skelScale}
	constructOptions.Scale.X *= g.dir

	skel := &g.skellington
	animIdx := 0
	if g.moving {
		animIdx = 1
	}
	if g.velocityY < 0 {
		animIdx = 2
	} else if g.playerPos.Y != g.groundY {
		animIdx = 3
	}

	if g.lastAnimidx != animIdx {
		g.skelTime = time.Now()
		g.lastAnimidx = animIdx
	}

	skullScaleY := &bone("Skull", skel.Bones).Scale.Y
	hatRot := &bone("Hat", skel.Bones).Rot

	// skull and hat might be negative from previous frame if looking the other way,
	// so they need to be reset to prevent interpolation from animate()
	*skullScaleY = abs(*skullScaleY)
	*hatRot = abs(*hatRot)

	// animate skellington
	tf0 := skf.TimeFrame(skel.Animations[animIdx], time.Since(g.skelTime), false, true)
	anims := []skf.Animation{skel.Animations[animIdx]}
	skf_e.Animate(skel, anims, []int{tf0}, []int{20})

	// point left shoulder and head to mouse
	mx, my := ebiten.CursorPosition()
	mouseX, mouseY := float32(mx), float32(my)
	mousePos := skf.Vec2{
		X: (-g.playerPos.X/skelScale + mouseX/skelScale) * g.dir,
		Y: g.playerPos.Y/skelScale - mouseY/skelScale,
	}
	bone("Left Shoulder Target", skel.Bones).Pos = mousePos
	bone("Looker", skel.Bones).Pos = mousePos
	if g.style == 0 {
		bone("Hat", skel.Bones).Pos.Y = 520
	} else {
		bone("Hat", skel.Bones).Pos.Y = 600
	}

	// flip skull and hat, and switch shoulder constraint, if looking the other way
	shoulder := bone("LSIK", skel.Bones)
	if (g.dir == 1 && mouseX < g.playerPos.X) || (g.dir == -1 && mouseX > g.playerPos.X) {
		*skullScaleY = -abs(*skullScaleY)
		*hatRot = -abs(*hatRot)
		shoulder.Ik_constraint = "Clockwise"
	} else {
		shoulder.Ik_constraint = "CounterClockwise"
	}

	// construct and draw skellington
	finalBones := skf_e.Construct(*skel, constructOptions)
	skf_e.Draw(finalBones, []skf.Style{skel.Styles[g.style], skel.Styles[1]}, g.skelTex, screen)

	msg := fmt.Sprintf("A - Move Left\nD - Move Right\nSpace - Jump\n1, 2 - Change outfit\nSkellington will look at and reach for cursor")
	op := &text.DrawOptions{}
	op.LayoutOptions.LineSpacing = 27
	op.GeoM.Translate(10, 20)
	op.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, msg, &text.GoTextFace{Source: mplusFaceSource, Size: 20}, op)

}

func (g *Game) Draw(screen *ebiten.Image) {
	g.Skellington(screen)
	//g.Skellina(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1280, 720
}

func main() {
	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle("SkelForm Basic Animation Demo")
	skellington, skelTex := skf_e.Load(assetPath("skellington.skf"))
	skellina, skelaTex := skf_e.Load(assetPath("skellina.skf"))
	size_x, size_y := ebiten.WindowSize()
	groundY := float32(size_y)/2 + 50
	if err := ebiten.RunGame(&Game{
		skellington: skellington,
		skelTex:     skelTex,
		skellina:    skellina,
		skelaTex:    skelaTex,
		dir:         1,
		skelTime:    time.Now(),
		skelaTime:   time.Now(),
		style:       1,
		groundY:     groundY,
		playerPos: skf.Vec2{
			X: float32(size_x) / 2,
			Y: groundY,
		},
	}); err != nil {
		log.Fatal(err)
	}
}

func init() {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		log.Fatal(err)
	}
	mplusFaceSource = s
}

func abs(value float32) float32 {
	return float32(math.Abs(float64(value)))
}

func bone(name string, bones []skf.Bone) *skf.Bone {
	finalBone := &skf.Bone{}
	for i, bone := range bones {
		if bone.Name == name {
			finalBone = &bones[i]
		}
	}
	return finalBone
}
