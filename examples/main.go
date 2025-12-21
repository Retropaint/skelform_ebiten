package main

import (
	"log"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	// "github.com/hajimehoshi/ebiten/v2/ebitenutil"
	skf_e "github.com/retropaint/skelform_ebiten"
	skf "github.com/retropaint/skelform_go"
)

type Game struct {
	skellington skf.Armature
	skelTex     []*ebiten.Image
	playerPos   skf.Vec2
	moving      bool
	start_time  time.Time
	dir         float32
}

func (g *Game) Update() error {
	var speed float32 = 7
	g.moving = false

	var keys []ebiten.Key
	for _, key := range inpututil.AppendJustPressedKeys(keys) {
		switch key {
		case ebiten.KeyA:
			g.start_time = time.Now()
		case ebiten.KeyD:
			g.start_time = time.Now()
		}
	}
	for _, key := range inpututil.AppendJustReleasedKeys(keys) {
		switch key {
		case ebiten.KeyA:
			g.start_time = time.Now()
		case ebiten.KeyD:
			g.start_time = time.Now()
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
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	animOptions := skf_e.AnimOptions{}
	animOptions.Init()

	skelScale := float32(0.125)
	animOptions.Position = g.playerPos
	animOptions.Scale = skf.Vec2{X: skelScale, Y: skelScale}
	animOptions.Scale.X *= g.dir
	animOptions.BlendFrames = 10

	skel := &g.skellington
	animIdx := 0
	if g.moving {
		animIdx = 1
	}

	skullScaleY := &bone("Skull", skel.Bones).Scale.Y
	hatRot := &bone("Hat", skel.Bones).Rot

	// skull and hat might be negative from previous frame if looking the other way,
	// so they need to be reset to prevent interpolation from animate()
	*skullScaleY = abs(*skullScaleY)
	*hatRot = abs(*hatRot)

	// animate skellington
	tf0 := skf.TimeFrame(skel.Animations[animIdx], time.Since(g.start_time), false, true)
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

	// flip skull and hat if looking the other way
	if (g.dir == 1 && mouseX < g.playerPos.X) || (g.dir == -1 && mouseX > g.playerPos.X) {
		*skullScaleY = -abs(*skullScaleY)
		*hatRot = -abs(*hatRot)
	}

	// construct and draw skellington
	finalBones := skf_e.Construct(*skel, animOptions)
	skf_e.Draw(finalBones, skel.Styles, g.skelTex, screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1280, 720
}

func main() {
	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle("SkelForm Basic Animation Demo")
	skellington, textures := skf_e.Load("skellington.skf")
	size_x, size_y := ebiten.WindowSize()
	var ebTex []*ebiten.Image
	for _, tex := range textures {
		ebTex = append(ebTex, ebiten.NewImageFromImage(tex))
	}
	if err := ebiten.RunGame(&Game{
		skellington: skellington,
		skelTex:     ebTex,
		dir:         1,
		start_time:  time.Now(),
		playerPos: skf.Vec2{
			X: float32(size_x) / 2,
			Y: float32(size_y) / 2,
		},
	}); err != nil {
		log.Fatal(err)
	}
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
