Library for running [SkelForm](https://skelform.org) animations in
[Ebitengine](https://ebitengine.org/).

## Interactive Example

```
go run github.com/retropaint/skelform_ebiten/examples/character@v0
```

Example includes:

- Simple state-based animations (idling, running, jumping, landing)
- Inverse kinematics
- Mesh deformation
- 2 styles

## Importing

```
import (
  skf_e "github.com/retropaint/skelform_ebiten"
  skf "github.com/retropaint/skelform_go"
)
```

Note: you may need to import `skelform_go` as well for certain structs (`Vec2`,
`Bone`, etc).

## Basic Setup

- `skf_e.Load()` - loads `.skf` file and returns armature & textures, to be used
  later
- `skf_e.Animate()` - transforms the armature's bones based on the animation(s)
- `skf_e.Construct()` - provides the bones from this armature that are ready for
  use
- `skf_e.Draw()` - draws the bones on-screen, with the provided style(s)

### 1. Load:

```
(armature, textures) = skf_e.Load("skelform_file.skf")
```

This should only be called once (eg; before main game loop), and `armature` and
`textures` should be kept for later use.

### 2\. Animate:

```
# use `TimeFrame()` to get the animation frame based on time
time = time.Now()
timeFrame := skf_e.TimeFrame(skel.Animations[0], time.Since(skelTime), false, true)

anim := []skf.Animation{armature.Animations[0]}
skf_e.Animate(armature, anim, []int{time_frame}, []int{0})
```

_Note: not needed if armature is statilc_

### 3\. Construct:

```
constructOptions := skf_e.ConstructOptions{}
constructOptions.Init()

sizeX, sizeY := ebiten.WindowSize()
constructOptions.Position = skf.Vec2{X: sizeX/2, Y: sizeY/2}

finalBones := skf_e.Construct(*armature, constructOptions)
```

Modifications to the armature (eg; aiming at cursor) may be done before or after
construction.

### 4\. Draw:

```
skf_e.Draw(finalBones, armature.Styles, ebTextures, screen)
```
