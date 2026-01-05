package skelform_ebiten

import (
	"image"
	"math"
	"sort"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/retropaint/skelform_go"
)

// Options for `Construct()`.
//
// Position: adds each bone's position by this much
// Scale: multiplies each bone's scale by this much
type ConstructOptions struct {
	Position skelform_go.Vec2
	Scale    skelform_go.Vec2
}

func (co *ConstructOptions) Init() {
	co.Position = skelform_go.Vec2{X: 0, Y: 0}
	co.Scale = skelform_go.Vec2{X: 0.25, Y: 0.25}
}

// Transforms an armature's bones based on the provided animation(s) and their frame(s).
//
// `smoothFrames` is used to smoothly interpolate transforms. Mainly used for smooth animation transitions. Higher frames are smoother.
//
// Note: smoothFrames should ideally be set to 0 (or empty) when reversing animations.
func Animate(armature *skelform_go.Armature, animations []skelform_go.Animation, frames []int, smoothFrames []int) {
	skelform_go.Animate(armature, animations, frames, smoothFrames)
}

// Returns the constructed array of bones from this armature.
//
// While constructing, several options (positional offset, scale) may be set.
func Construct(armature skelform_go.Armature, constOptions ConstructOptions) []skelform_go.Bone {
	finalBones := skelform_go.Construct(&armature)

	for b := range finalBones {
		bone := &finalBones[b]
		bone.Scale = bone.Scale.Mul(constOptions.Scale)
		bone.Pos.Y = -bone.Pos.Y
		bone.Pos = bone.Pos.Mul(constOptions.Scale)
		bone.Pos = bone.Pos.Add(constOptions.Position)

		skelform_go.CheckBoneFlip(bone, constOptions.Scale)

		for v := range finalBones[b].Vertices {
			vert := &finalBones[b].Vertices[v]
			vert.Pos.Y = -vert.Pos.Y
			vert.Pos = vert.Pos.Mul(constOptions.Scale)
			vert.Pos = vert.Pos.Add(constOptions.Position)
		}
	}

	return finalBones
}

// Draws the bones to the provided screen, using the provided styles and textures.
//
// Recommended: include the whole texture array from the file even if not all will be used,
// as the provided styles will determine the final appearance.
func Draw(bones []skelform_go.Bone, styles []skelform_go.Style, textures []*ebiten.Image, screen *ebiten.Image) {
	sort.Slice(bones, func(i, j int) bool {
		return bones[i].Zindex < bones[j].Zindex
	})

	finalTextures := skelform_go.SetupBoneTextures(bones, styles)

	for b := range bones {
		tex, ok := finalTextures[uint(bones[b].Id)]
		if !ok {
			continue
		}

		if len(bones[b].Vertices) > 0 {
			drawMesh(bones[b], tex, textures[tex.AtlasIdx], screen)
			continue
		}

		// crop texture to this bone
		tex_offset := skelform_go.Vec2{
			X: tex.Offset.X,
			Y: tex.Offset.Y,
		}
		tex_size := skelform_go.Vec2{
			X: tex.Size.X,
			Y: tex.Size.Y,
		}
		sub := textures[tex.AtlasIdx].SubImage(image.Rectangle{
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
			X: tex.Size.X / 2 * bones[b].Scale.X,
			Y: tex.Size.Y / 2 * bones[b].Scale.Y,
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

// Returns the properly bound animation frame based on the provided animation.
func FormatFrame(anim skelform_go.Animation, frame int, reverse bool, loop bool) int {
	return skelform_go.FormatFrame(anim, frame, reverse, loop)
}

// Returns the animation frame based on the provided time.
func TimeFrame(anim skelform_go.Animation, time time.Duration, reverse bool, loop bool) int {
	return skelform_go.TimeFrame(anim, time, reverse, loop)
}

// Loads an `.skf` file.
func Load(path string) (skelform_go.Armature, []*ebiten.Image) {
	armature, textures := skelform_go.Load(path)
	var ebTextures []*ebiten.Image
	for _, tex := range textures {
		ebTextures = append(ebTextures, ebiten.NewImageFromImage(tex))
	}
	return armature, ebTextures
}
