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
- `Sampler` is an interface with `GetFrameValue(x,y,t float) float`. It can be converted into a texture using `colors.GetAniTextureFromSampler` along with a texture

Consider a new type:

- `DynamicObject`, which is a collection of `DynamicTriangles` along with some transformations, applied with either `.WithTransform(matrix)` or `.WithDynamicTransform(func(float64) HomogeneousMatrix)`
- A `DynamicObject` can be evaluated at `.Frame(float64)` to get `StaticObject`, which consists of `StaticTriangles`

## A common pattern:

I often run into having three types of a Thing:
* Animated version, evaluated at `(locus Position, t float) retType`
* Dynamic version, evaluated at `(t float) StaticVersion`
* Static version, fixed to a frame, evaluated at `(locus Position) retType`

Here `locus Position` is any positional value. It can be a pixel, it can be a texture space.
The following types implement this approach, with some gaps:
* Textures, implements Dynamic, Static, and Animated version
* Triangles, implements Dynamic and Static version
* Object, implementing Dynamic and Static, no Animated version
* Sampler, implements Animated versions
* Transparency, implements Animated, Dynamic, Static
* Scene, implements Dynamic and Static
* Background, implements Dynamic and Static

Ideally there would be boilerplate that would accomplish this out of the box, without having to rewrite the same code over and over again

## TODOs:

- Add more intense gradients, like Bezier, etc
- Is there some way to create a mapping from pixel space to Window index? If all Windows were uniform, this would be trivial, but they're of various sizes. Maybe storing them as a tree could help? But lookup would be O(log(n))
- Investiate if there are any optimizations using GPU/CUDA
- Add Phong lighting model (this could be another rendering mode)
- Add simple geometry based textures, like the transition between the faces of a cube
- Add the ability to slice arbitrarily into a Perelman slice
- Specify a structure for animation sequences
- Add textures from images
- Add a kaleidoscope, i.e. render a full scene, but fetch pixels from a triangle slice of the full scene
- Develop a way to do visual aberration, where colors move with a slight delay
  - Or color movement is distorted around the edges of the screen
- Explore partial transparencies
  - try Moire patterns
- Add pixel size to rendering, to allow for primitives like point and line, etc
- Move ray calculation to the rendering engine, allow for zooming in/out of the scene
- Create different render methods, generalize the wireframe render. Of course, this would also affect the math to map an image point back onto the pixel value, as well as the code for estimating bounding boxes and wireframes. The current renderer is hardcoded to only handle planar point-based rendering.
  - Planar point-based rendering necessarily renders constant angles onto more pixels closer to the edge, this is why spheres are elongated towards the edge
  - Flat projection would fix this, the rays would be emitted not from one point, but from an image plane, all in the same direction. Equivalent to planar point-based, but with infinite focal length.
  - spherical point-based projection would hold the angle constant everywhere. first, estimate the angle covered by the angle, then divide that (not pixels) into equal parts, then create rays
- Explore using OpenGL bindings (they should work for 1.14, not sure about higher versions)
- Fix homo matrix multiplication to work differently on points vs vectors:
 - points get the full matrix
 - vectors get Mv - M0, i.e. the matrix also needs to be applied to the origin point, for consistency's sake
 - unit vectors get the same handling as vectors, but get normalized after
- Add separate UnitVector type? (would this make math too complex? Do we need to handle the various combinations for adding, etc)

## Ideas from Insta:

- https://www.instagram.com/reel/C5IRkXhgPQp/?utm_source=ig_web_copy_link
- https://www.instagram.com/p/C5G6c3_toa0/?img_index=1
- https://www.instagram.com/p/C5KgPvJovVe/
- https://www.instagram.com/p/C5DQWcuomyF/
- https://www.instagram.com/p/C5MRFySLnfy/
- https://www.instagram.com/p/C46t0QLrsHH/?img_index=1
- https://www.instagram.com/p/C4n6hKVOcIT/?img_index=1
- https://www.instagram.com/p/C4QLY-lNUYs/?img_index=1
- https://www.instagram.com/p/C4nfzZ8yJVz/

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
