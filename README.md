# go-scene-renderer

Rendering visual scenes in golang

`go run ./... -video test out.mp4`

Example (converted from mp4 to gif for illustrative purposes):

![alt text](https://github.com/libeks/go-scene-renderer/blob/main/gallery/cube_sine.gif)

See more examples on [Insta](https://www.instagram.com/cube.gif/).

## Organization

Rendering accepts these types/interfaces:

- `Background`, which specifies the color at each x,y coordinate, each in the range of (-1,1). A background is not intended to contain objects, as no optimizations will not be performed in their rendering.
- `Scene`, a set of objects with a `Background`, intended to be displayed in a single frame. When rendering,
  the `Flatten` method will be called, which returns a set of `Triangle`s and a `Frame` for the background.
  - When rendering, these triangles are organized into a set of render `Window`s, each describing the sub-Scene of the square part of the image. This is done to make scene rendering much more efficient.
- `DynamicBackground`, a entity that returns a `Background` for each timestamp in (0,1), not intended to contain any objects
- `DynamicScene`, an entity that returns a Scene for each timestamp in (0,1)

As with other places, `Dynamic` refers to being renderable on a range of frames, whereas `Static` describes a snapshot in a single frame.

For help, here are other types of interfaces/objects which may come in handy:

- `Texture` and `DynamicTexture`, each specifying the color at each pixel from (0,1) in two dimensions. The domain of a Texture (0,1) is different from Frame (-1,1), but a Texture can be used as a Frame using the TextureToFrame helper.
  - A `DynamicTexture` can be converted into `DynamicBackground` using `BackgroundFromTexture()`
  - A static `Texture` can be converted into `DynamicTexture` using `StaticTexture()`
- `Gradient`, specifying a color from a gradient, in the range (0,1)
- `DynamicObject` is an object in a scene, which has a `Frame(float64)` method, returning a `StaticObject` (a collection of `StaticTriangles`), and a `GetWireframe` method, allowing for wireframe rendering.
- `Triangle` is the basic entity of object rendering. Triangles are bidirectional, with `DynamicTriangle` and `StaticTriangle` versions, skinned with the respective types of `Texture`.
- `Parallelogram` is a helper that contains two adjoining triangles in a plane, it contains a helper for mapping textures correctly onto the two contained triangles.
- `HomogeneousMatrix` contains the logic for doing three types of homogeneous transformations, which are:
  _ Translation by an arbitrary 3D vector (`TranslationMatrix`),
  _ Rotation by a radian angle around one of the three axes (`RotateMatrixX`, `RotateMatrixY`,`RotateMatrixZ`), and \* Scaling of all axes (`ScaleMatrix`).
- These matrices can be combined using `MatrixProduct`, applied right to left.

Consider a new type:

- `DynamicObject`, which is a collection of `DynamicTriangles` along with some transformations, applied with either `.WithTransform(matrix)` or `.WithDynamicTransform(func(float64) HomogeneousMatrix)`
- A `DynamicObject` can be evaluated at `.Frame(float64)` to get `StaticObject`, which consists of `StaticTriangles`

## TODOs:

- Add more intense gradients, like Bezier, etc
- Find Lab color space transformation code, try that out
- Is there some way to create a mapping from pixel space to Window index? If all Windows were uniform, this would be trivial, but they're of various sizes. Maybe storing them as a tree could help? But lookup would be O(log(n))
- Investiate if there are any optimizations using GPU/CUDA
- Add spring-loaded-mass interactions
  - See https://www.desmos.com/calculator/k01p40v0ct
  - `y=P0*e^{\alpha x}\cos\left(\beta x\right)+C6e^{\alpha x}\sin\left(\beta x\right)\left\{0<x<t\right\}`
    - P0 is initial position, C6 is the velocity component (?), b is friction component, k is spring constant, etc
- Add Phong lighting model
- Improve bounding box computation when one vertex of triangles goes behind the camera
- Add textures from images

## Further Reading:

- https://graphics.stanford.edu/courses/cs348a-09-fall/Handouts/handout15.pdf
- https://en.wikipedia.org/wiki/Homogeneous_coordinates
- https://en.wikipedia.org/wiki/Rotation_matrix
- https://en.wikipedia.org/wiki/Euler_angles#Conventions_by_intrinsic_rotations
- https://www.reddit.com/r/mobilevrstation/comments/xgamsr/typical_ios_shutteringnonplayback_fixes/
- https://users.csc.calpoly.edu/~zwood/teaching/csc471/2017F/barycentric.pdf
- https://en.wikipedia.org/wiki/Barycentric_coordinate_system

## Benchmarks:

Rendering the `SpinningMulticube` scene with the `-video=intermediate` preset takes this long:

- 1m20s - base
- 50s - by caching triangle intermediate results

The same scene with `-video=hidef` (approx 64x the effort, at 4x resolution and 4x anti-aliasing, 4x more frames) takes:

- 1hr - base
- 45m - by caching intermediate triangle results
- 30s - by using windowed triangle culling (very impressive)

# Frequently Asked Questions

- How do I run pprof?
  - Set pprof=true in main, then after the run do `go tool pprof -png . cpu.pprof` or `go tool pprof -png . mem.pprof`
- How do I add frame numbers to a video?
  - After the video is rendered, run something like this:
    `ffmpeg -i gallery/out_hd.mp4 \ -vf "drawtext=fontfile=Arial.ttf: text=%{n}: x=(w-tw)/2: y=h-(2*lh): fontcolor=white: box=1: boxcolor=0x00000099" \ gallery/out_with_frames.mp4`
