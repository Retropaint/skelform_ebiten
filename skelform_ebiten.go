package skelform_ebiten

import (
	"image"
	"math"
	"sort"
	"time"

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

func Animate(
	bones []skelform_go.Bone,
	ikFamilies []skelform_go.IkFamily,
	animation skelform_go.Animation,
	screen *ebiten.Image,
	frame int,
	anim_options AnimOptions,
) []skelform_go.Bone {
	var animatedBones []skelform_go.Bone
	for _, bone := range bones {
		animatedBones = append(animatedBones, bone)
	}

	skelform_go.Animate(animatedBones, animation, frame)

	var inheritedBones []skelform_go.Bone
	for _, bone := range animatedBones {
		inheritedBones = append(inheritedBones, bone)
	}

	skelform_go.Inheritance(inheritedBones, make(map[uint]float32))
	var ikRots map[uint]float32
	for i := 0; i < 10; i++ {
		ikRots = skelform_go.InverseKinematics(inheritedBones, ikFamilies)
	}
	skelform_go.Inheritance(animatedBones, ikRots)

	for b := range animatedBones {
		bone := &animatedBones[b]
		bone.Scale = bone.Scale.Mul(skelform_go.Vec2{X: anim_options.Scale_factor, Y: anim_options.Scale_factor})
		bone.Pos.Y = -bone.Pos.Y
		bone.Pos = bone.Pos.Mul(skelform_go.Vec2{X: anim_options.Scale_factor, Y: anim_options.Scale_factor})
		bone.Pos = bone.Pos.Add(anim_options.Pos_offset)
	}

	return animatedBones
}

func Draw(bones []skelform_go.Bone, styles []skelform_go.Style, texture image.Image, screen *ebiten.Image) {
	tex := ebiten.NewImageFromImage(texture)

	sort.Slice(bones, func(i, j int) bool {
		return bones[i].Zindex < bones[j].Zindex
	})

	for _, bone := range bones {
		if len(bone.Style_ids) == 0 {
			continue
		}

		texFields := styles[0].Textures[bone.Tex_idx]

		// crop texture to this bone
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

		// center bone for scale & rot operations
		size := skelform_go.Vec2{
			X: texFields.Size.X / 2 * bone.Scale.X,
			Y: texFields.Size.Y / 2 * bone.Scale.Y,
		}
		cos := math.Cos(float64(bone.Rot))
		sin := math.Sin(float64(bone.Rot))
		bone.Pos.X -= size.X*float32(cos) + size.Y*float32(sin)
		bone.Pos.Y += size.X*float32(sin) - size.Y*float32(cos)

		op.GeoM.Scale(float64(bone.Scale.X), float64(bone.Scale.Y))
		op.GeoM.Rotate(float64(-bone.Rot))

		op.GeoM.Translate(float64(bone.Pos.X), float64(bone.Pos.Y))

		screen.DrawImage(sub.(*ebiten.Image), op)
	}
}

func FormatFrame(anim skelform_go.Animation, frame int, reverse bool, loop bool) int {
	return skelform_go.FormatFrame(anim, frame, reverse, loop)
}

func TimeFrame(anim skelform_go.Animation, time time.Duration, reverse bool, loop bool) int {
	return skelform_go.TimeFrame(anim, time, reverse, loop)
}

func Load(path string) (skelform_go.Root, image.Image) {
	return skelform_go.Load(path)
}
