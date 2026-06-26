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
	Velocity skelform_go.Vec2
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
func Construct(armature *skelform_go.Armature, constOptions ConstructOptions) {
	skelform_go.Construct(armature)

	for b := range armature.Constructed_bones {
		bone := &armature.Constructed_bones[b]

		bone.Scale = bone.Scale.Mul(constOptions.Scale)
		bone.Pos.Y = -bone.Pos.Y
		bone.Pos = bone.Pos.Mul(constOptions.Scale)
		bone.Pos = bone.Pos.Add(constOptions.Position)

		if skelform_go.IsFacingLeft(constOptions.Scale) {
			bone.Rot = -bone.Rot
		}

		if armature.Bones[b].Physics_id != -1 {
			phys := armature.Physics[armature.Bones[b].Physics_id]
			phys.Global_pos = phys.Global_pos.Sub(constOptions.Velocity)
		}

		if armature.Bones[b].Visuals_id != -1 {
			visual := armature.Visuals[armature.Bones[b].Visuals_id]
			for v := range visual.Vertices {
				vert := &visual.Vertices[v]

				vert.Pos.Y = -vert.Pos.Y
				vert.Pos = vert.Pos.Mul(constOptions.Scale)
				vert.Pos = vert.Pos.Add(constOptions.Position)
			}
		}

	}
}

// Draws the bones to the provided screen, using the provided styles and textures.
//
// Recommended: include the whole texture array from the file even if not all will be used,
// as the provided styles will determine the final appearance.
func Draw(bones []skelform_go.Bone, visuals []skelform_go.Visuals, styles []skelform_go.Style, textures []*ebiten.Image, screen *ebiten.Image) {
	// sort bones by Zindex
	sort.Slice(bones, func(i, j int) bool {
		if bones[i].Visuals_id == -1 || bones[j].Visuals_id == -1 {
			return false
		}
		visualsA := visuals[bones[i].Visuals_id]
		visualsB := visuals[bones[j].Visuals_id]
		return visualsA.Zindex <= visualsB.Zindex
	})

	for _, bone := range bones {
		if bone.Visuals_id == -1 {
			continue
		}
		visual := visuals[bone.Visuals_id]

		// get this bone's current texture
		tex, err := skelform_go.GetBoneTexture(visual.Tex, styles)
		if err != nil {
			continue
		}

		// will be used to flip pivot rotations if necessary
		var dir float32
		dir = 1.0
		if skelform_go.IsFacingLeft(bone.Scale) {
			dir = -1.0
		}

		// setup texture pivot
		pivotPos := visual.Pivot_pos.Mul(tex.Size)
		pivotPos = skelform_go.RotateVec2(pivotPos, float64(bone.Rot*dir)).Mul(bone.Scale).Mul(visual.Pivot_scale)
		pivotPos.Y = -pivotPos.Y

		// draw mesh
		if len(visual.Vertices) > 0 {
			drawMesh(visual, tex, textures[tex.AtlasIdx], screen)
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

		scale := skelform_go.Vec2{X: bone.Scale.X * visual.Pivot_scale.X, Y: bone.Scale.Y * visual.Pivot_scale.Y}

		// center bone for scale & rot operations
		size := skelform_go.Vec2{X: tex.Size.X / 2 * scale.X, Y: tex.Size.Y / 2 * scale.Y}
		cos := math.Cos(float64(bone.Rot + visual.Pivot_rot*dir))
		sin := math.Sin(float64(bone.Rot + visual.Pivot_rot*dir))
		bone.Pos.X -= size.X*float32(cos) + size.Y*float32(sin)
		bone.Pos.Y += size.X*float32(sin) - size.Y*float32(cos)

		op.GeoM.Scale(float64(scale.X), float64(scale.Y))
		op.GeoM.Rotate(float64(-bone.Rot - visual.Pivot_rot*dir))

		op.GeoM.Translate(float64(bone.Pos.X+pivotPos.X), float64(bone.Pos.Y+pivotPos.Y))

		screen.DrawImage(sub.(*ebiten.Image), op)
	}
}

func drawMesh(visual skelform_go.Visuals, tex skelform_go.Texture, fullTex *ebiten.Image, screen *ebiten.Image) {
	var verts []ebiten.Vertex
	var indices []uint16
	for _, vert := range visual.Vertices {
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
	for _, idx := range visual.Indices {
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
