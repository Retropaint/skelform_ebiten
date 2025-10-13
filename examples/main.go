package main

import (
	"image"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	// "github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/retropaint/skelform_ebiten"
	skf "github.com/retropaint/skelform_go"
)

type Game struct {
	skellington         skf.Root
	witsy               skf.Root
	skellington_texture image.Image
	witsy_texture       image.Image
	frame               int
	player_pos          skf.Vec2
	moving              bool
	start_time          time.Time
}

func (g *Game) Update() error {
	var speed float32 = 10
	g.moving = false

	var keys []ebiten.Key
	for _, key := range inpututil.AppendPressedKeys(keys) {
		switch key {
		case ebiten.KeyW:
			g.player_pos.Y -= speed
			g.moving = true
		case ebiten.KeyS:
			g.player_pos.Y += speed
			g.moving = true
		case ebiten.KeyA:
			g.player_pos.X -= speed
			g.moving = true
		case ebiten.KeyD:
			g.player_pos.X += speed
			g.moving = true
		case ebiten.KeyUp:
			g.skellington.Armature.Bones[1].Pos.Y += speed
		case ebiten.KeyRight:
			g.skellington.Armature.Bones[1].Pos.X += speed
		case ebiten.KeyLeft:
			g.skellington.Armature.Bones[1].Pos.X -= speed
		case ebiten.KeyDown:
			g.skellington.Armature.Bones[1].Pos.Y -= speed
		case ebiten.KeyF:
			g.moving = true
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ao := skelform_ebiten.AnimOptions{}
	ao.Init()

	ao.Position = g.player_pos
	ao.BlendFrames = 0

	skel := &g.skellington.Armature

	//idx := 0
	if g.moving {
		//idx = 3
		ao.BlendFrames = 20
	}

	tf0 := skf.TimeFrame(skel.Animations[0], time.Since(g.start_time), false, true)
	tf1 := skelform_ebiten.TimeFrame(skel.Animations[1], time.Since(g.start_time), false, true)
	animatedBones := skelform_ebiten.Animate(
		skel,
		[]skf.Animation{skel.Animations[0], skel.Animations[1]},
		[]int{tf0, tf1},
		screen,
		ao,
	)
	skelform_ebiten.Draw(animatedBones, skel.Styles, g.skellington_texture, screen)

	g.frame += 1
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1280, 720
}

func main() {
	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle("Hello, World!")
	skellington, skellington_texture := skelform_ebiten.Load("untitled.skf")
	witsy, witsy_texture := skelform_ebiten.Load("witsy.skf")
	size_x, size_y := ebiten.WindowSize()
	if err := ebiten.RunGame(&Game{
		skellington:         skellington,
		witsy:               witsy,
		skellington_texture: skellington_texture,
		witsy_texture:       witsy_texture,
		start_time:          time.Now(),
		player_pos: skf.Vec2{
			X: float32(size_x) / 2,
			Y: float32(size_y) / 2,
		},
	}); err != nil {
		log.Fatal(err)
	}
}
