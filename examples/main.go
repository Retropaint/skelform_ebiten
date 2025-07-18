package main

import (
	"image"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	// "github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/retropaint/skelform_go"
	"github.com/retropaint/skelform_ebiten"
)

type Game struct {
	root       skelform_go.Root
	texture    image.Image
	frame      int
	player_pos skelform_go.Vec2
	moving     bool
	start_time time.Time
}

func (g *Game) Update() error {
	var speed float32 = 5
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
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ao := skelform_ebiten.AnimOptions{}
	ao.Init()

	ao.Pos_offset = g.player_pos

	anim_idx := 0
	if g.moving {
		anim_idx = 1
	}

	elapsed := time.Now().UnixMilli() - g.start_time.UnixMilli()
	time_frame := skelform_ebiten.Get_frame_by_time(g.root.Armature.Animations[anim_idx], elapsed, false)
	skelform_ebiten.Animate(screen, g.root.Armature, g.texture, anim_idx, time_frame, ao)
	g.frame += 1
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1280, 720
}

func main() {
	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle("Hello, World!")
	root, texture := skelform_ebiten.Load("skellington.skf")
	size_x, size_y := ebiten.WindowSize()
	if err := ebiten.RunGame(&Game{
		root:       root,
		texture:    texture,
		start_time: time.Now(),
		player_pos: skelform_go.Vec2{
			X: float32(size_x) / 2,
			Y: float32(size_y) / 2,
		},
	}); err != nil {
		log.Fatal(err)
	}
}
