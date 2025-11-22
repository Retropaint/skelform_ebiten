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
	Position    skelform_go.Vec2
	Scale       skelform_go.Vec2
	BlendFrames int
}

func (ao *AnimOptions) Init() {
	ao.Position = skelform_go.Vec2{X: 0, Y: 0}
	ao.Scale = skelform_go.Vec2{X: 0.25, Y: 0.25}
	ao.BlendFrames = 0
}

func Animate(
	armature *skelform_go.Armature,
	animations []skelform_go.Animation,
	frames []int,
	blendFrames int,
) {
	for i := range animations {
		skelform_go.Animate(armature.Bones, animations[i], frames[i], blendFrames)
	}

	skelform_go.ResetBones(armature.Bones, animations, frames[0], blendFrames)
}

func Construct(armature skelform_go.Armature, animOptions AnimOptions) []skelform_go.Bone {
	var inheritedBones []skelform_go.Bone
	for _, bone := range armature.Bones {
		inheritedBones = append(inheritedBones, bone)
	}
	skelform_go.Inheritance(inheritedBones, make(map[uint]float32))
	ikRots := skelform_go.InverseKinematics(inheritedBones, armature.Ik_root_ids)

	var finalBones []skelform_go.Bone
	for _, bone := range armature.Bones {
		finalBones = append(finalBones, bone)
		finalBones[len(finalBones)-1].Vertices = nil
		for _, vert := range bone.Vertices {
			finalBones[len(finalBones)-1].Vertices = append(finalBones[len(finalBones)-1].Vertices, vert)
		}
	}
	skelform_go.Inheritance(finalBones, ikRots)
	skelform_go.ConstructVerts(finalBones)

	for b := range finalBones {
		bone := &finalBones[b]
		bone.Scale = bone.Scale.Mul(animOptions.Scale)
		bone.Pos.Y = -bone.Pos.Y
		bone.Pos = bone.Pos.Mul(animOptions.Scale)
		bone.Pos = bone.Pos.Add(animOptions.Position)

		// reverse rot if either scale is negative
		either := animOptions.Scale.X < 0 || animOptions.Scale.Y < 0
		both := animOptions.Scale.X < 0 && animOptions.Scale.Y < 0
		if either && !both {
			bone.Rot = -bone.Rot
		}

		for v := range finalBones[b].Vertices {
			vert := &finalBones[b].Vertices[v]
			vert.Pos.Y = -vert.Pos.Y
			vert.Pos = vert.Pos.Mul(animOptions.Scale)
			vert.Pos = vert.Pos.Add(animOptions.Position)
		}
	}

	return finalBones
}

func Draw(bones []skelform_go.Bone, styles []skelform_go.Style, texture *ebiten.Image, screen *ebiten.Image) {
	sort.Slice(bones, func(i, j int) bool {
		return bones[i].Zindex < bones[j].Zindex
	})

	for b := range bones {
		if len(bones[b].Style_ids) == 0 {
			continue
		}

		texFields := styles[0].Textures[bones[b].Tex_idx]

		if len(bones[b].Vertices) > 0 {
			drawMesh(bones[b], texFields, texture, screen)
			continue
		}

		// crop texture to this bone
		tex_offset := skelform_go.Vec2{
			X: texFields.Offset.X,
			Y: texFields.Offset.Y,
		}
		tex_size := skelform_go.Vec2{
			X: texFields.Size.X,
			Y: texFields.Size.Y,
		}
		sub := texture.SubImage(image.Rectangle{
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
			X: texFields.Size.X / 2 * bones[b].Scale.X,
			Y: texFields.Size.Y / 2 * bones[b].Scale.Y,
		}
		cos := math.Cos(float64(bones[b].Rot))
		sin := math.Sin(float64(bones[b].Rot))
		bones[b].Pos.X -= size.X*float32(cos) + size.Y*float32(sin)
		bones[b].Pos.Y += size.X*float32(sin) - size.Y*float32(cos)

		op.GeoM.Scale(float64(bones[b].Scale.X), float64(bones[b].Scale.Y))
		op.GeoM.Rotate(float64(-bones[b].Rot))

		op.GeoM.Translate(float64(bones[b].Pos.X), float64(bones[b].Pos.Y))

		screen.DrawImage(sub.(*ebiten.Image), op)
	}
}

func drawMesh(bone skelform_go.Bone, tex skelform_go.Texture, fullTex *ebiten.Image, screen *ebiten.Image) {
	var verts []ebiten.Vertex
	var indices []uint16
	for _, vert := range bone.Vertices {
		eb_vert := ebiten.Vertex{
			DstX:   vert.Pos.X,
			DstY:   vert.Pos.Y,
			SrcX:   tex.Offset.X + vert.Uv.X*float32(tex.Size.X),
			SrcY:   tex.Offset.Y + vert.Uv.Y*float32(tex.Size.Y),
			ColorR: 1,
			ColorG: 1,
			ColorB: 1,
			ColorA: 1,
		}
		verts = append(verts, eb_vert)
	}
	for _, idx := range bone.Indices {
		indices = append(indices, uint16(idx))
	}
	screen.DrawTriangles(verts, indices, fullTex, &ebiten.DrawTrianglesOptions{})
}

func FormatFrame(anim skelform_go.Animation, frame int, reverse bool, loop bool) int {
	return skelform_go.FormatFrame(anim, frame, reverse, loop)
}

func TimeFrame(anim skelform_go.Animation, time time.Duration, reverse bool, loop bool) int {
	return skelform_go.TimeFrame(anim, time, reverse, loop)
}

func Load(path string) (skelform_go.Armature, image.Image) {
	return skelform_go.Load(path)
}
