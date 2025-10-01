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

func Animate(screen *ebiten.Image, armature skelform_go.Armature, texture image.Image, animIdx int, frame int, anim_options AnimOptions) {
	var animatedBones []skelform_go.Bone
	for _, bone := range armature.Bones {
		animatedBones = append(animatedBones, bone)
	}

	if animIdx < len(armature.Animations)-1 {
		animatedBones = skelform_go.Animate(armature, animIdx, frame)
	}

	var inheritedBones []skelform_go.Bone
	for _, bone := range animatedBones {
		inheritedBones = append(inheritedBones, bone)
	}

	skelform_go.Inheritance(inheritedBones, make(map[uint]float32))
	var ikRots map[uint]float32
	for i := 0; i < 10; i++ {
		ikRots = skelform_go.InverseKinematics(inheritedBones, armature.Ik_families)
	}
	skelform_go.Inheritance(animatedBones, ikRots)

	tex := ebiten.NewImageFromImage(texture)

	for _, bone := range animatedBones {
		if len(bone.Style_idxs) == 0 {
			continue
		}

		texFields := armature.Styles[0].Textures[bone.Tex_idx]

		// crop texture to this prop
		tex_offset := skelform_go.Vec2{
			X: texFields.Offset.X,
			Y: texFields.Offset.Y,
		}
		tex_size := skelform_go.Vec2{
			X: texFields.Size.X,
			Y: texFields.Size.Y,
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
		bone.Pos.Y = -bone.Pos.Y

		// center prop for scale & rot operations
		size := skelform_go.Vec2{
			X: texFields.Size.X / 2 * bone.Scale.X,
			Y: texFields.Size.Y / 2 * bone.Scale.Y,
		}
		cos := math.Cos(float64(bone.Rot))
		sin := math.Sin(float64(bone.Rot))
		bone.Pos.X -= size.X*float32(cos) + size.Y*float32(sin)
		bone.Pos.Y += size.X*float32(sin) - size.Y*float32(cos)

		op.GeoM.Scale(float64(bone.Scale.X*anim_options.Scale_factor), float64(bone.Scale.Y*anim_options.Scale_factor))
		op.GeoM.Rotate(-float64(bone.Rot))

		final_pos := skelform_go.Vec2{
			X: bone.Pos.X*anim_options.Scale_factor + anim_options.Pos_offset.X,
			Y: bone.Pos.Y*anim_options.Scale_factor + anim_options.Pos_offset.Y,
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
