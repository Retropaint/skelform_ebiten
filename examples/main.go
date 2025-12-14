package main

import (
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	// "github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/retropaint/skelform_ebiten"
	skf "github.com/retropaint/skelform_go"
)

type Game struct {
	skellington skf.Armature
	skelTex     []*ebiten.Image
	frame       int
	player_pos  skf.Vec2
	moving      bool
	start_time  time.Time
	dir         int
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
			g.player_pos.X -= speed
			g.moving = true
			g.dir = -1
		case ebiten.KeyD:
			g.player_pos.X += speed
			g.moving = true
			g.dir = 1
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ao := skelform_ebiten.AnimOptions{}
	ao.Init()

	ao.Position = g.player_pos
	ao.Scale = skf.Vec2{X: 0.125, Y: 0.125}
	ao.Scale.X *= float32(g.dir)
	ao.BlendFrames = 10

	skel := &g.skellington

	animIdx := 0

	if g.moving {
		animIdx = 1
	}

	tf0 := skf.TimeFrame(skel.Animations[animIdx], time.Since(g.start_time), false, true)
	//tf1 := skf.TimeFrame(skel.Animations[4], time.Since(g.start_time), false, true)
	skelform_ebiten.Animate(
		skel,
		[]skf.Animation{skel.Animations[animIdx]},
		[]int{tf0},
		//[]skf.Animation{},
		//[]int{0},
		20,
	)
	constructedBones := skelform_ebiten.Construct(
		*skel,
		ao,
	)
	skelform_ebiten.Draw(constructedBones, skel.Styles, g.skelTex, screen)

	g.frame += 1
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1280, 720
}

func main() {
	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle("SkelForm Basic Animation Demo")
	skellington, textures := skelform_ebiten.Load("skellington.skf")
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
		player_pos: skf.Vec2{
			X: float32(size_x) / 2,
			Y: float32(size_y) / 2,
		},
	}); err != nil {
		log.Fatal(err)
	}
}
