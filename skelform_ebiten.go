package skelform_ebiten

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/retropaint/skelform_go"
)

type AnimOptions struct {
	Pos_offset   skelform_go.Vec2
	Scale_factor float32
}

func (ao *AnimOptions) Init() {
	ao.Pos_offset = skelform_go.Vec2{X: 0, Y: 0}
	ao.Scale_factor = 0.25
}

func Animate(screen *ebiten.Image, armature skelform_go.Armature, texture image.Image, anim_idx int, frame int, anim_options AnimOptions) {
	animated_bones := skelform_go.Animate(armature, anim_idx, frame)
	tex := ebiten.NewImageFromImage(texture)

	var props []skelform_go.Bone
	for _, bone := range animated_bones {
		props = append(props, bone)

		prop := &bone

		if prop.Tex_idx == -1 {
			continue
		}

		// crop texture to this prop
		tex_offset := skelform_go.Vec2{
			X: armature.Textures[prop.Tex_idx].Offset.X,
			Y: armature.Textures[prop.Tex_idx].Offset.Y,
		}
		tex_size := skelform_go.Vec2{
			X: armature.Textures[prop.Tex_idx].Size.X,
			Y: armature.Textures[prop.Tex_idx].Size.Y,
		}
		sub := tex.SubImage(image.Rectangle{
			Min: image.Point{
				X: int(tex_offset.X),
				Y: int(tex_offset.Y),
			},
			Max: image.Point{
				X: int(tex_offset.X) + int(tex_size.X),
				Y: int(tex_offset.Y) + int(tex_size.Y),
			},
		})

		op := &ebiten.DrawImageOptions{}

		// Ebiten treats positive Y as down
		prop.Pos.Y = -prop.Pos.Y

		// center prop for scale & rot operations
		size := skelform_go.Vec2{
			X: armature.Textures[prop.Tex_idx].Size.X / 2 * prop.Scale.X,
			Y: armature.Textures[prop.Tex_idx].Size.Y / 2 * prop.Scale.Y,
		}
		cos := math.Cos(float64(prop.Rot))
		sin := math.Sin(float64(prop.Rot))
		prop.Pos.X -= size.X*float32(cos) + size.Y*float32(sin)
		prop.Pos.Y += size.X*float32(sin) - size.Y*float32(cos)

		op.GeoM.Scale(float64(prop.Scale.X*anim_options.Scale_factor), float64(prop.Scale.Y*anim_options.Scale_factor))
		op.GeoM.Rotate(-float64(prop.Rot))

		final_pos := skelform_go.Vec2{
			X: prop.Pos.X*anim_options.Scale_factor + anim_options.Pos_offset.X,
			Y: prop.Pos.Y*anim_options.Scale_factor + anim_options.Pos_offset.Y,
		}
		op.GeoM.Translate(float64(final_pos.X), float64(final_pos.Y))

		screen.DrawImage(sub.(*ebiten.Image), op)
	}
}

func Get_frame_by_time(anim skelform_go.Animation, unix_milli int64, reverse bool) int {
	return skelform_go.Get_frame_by_time(anim, unix_milli, reverse)
}

func Load(path string) (skelform_go.Root, image.Image) {
	return skelform_go.Load(path)
}
