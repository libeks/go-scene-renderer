# go-scene-renderer
Rendering visual scenes in golang

`go run ./... -video test out.mp4`

Example (converted from mp4 to gif for illustrative purposes):

![alt text](https://github.com/libeks/go-scene-renderer/blob/main/gallery/cube_sine.gif)

## TODOs:
* Add more intense gradients, like Bezier, etc
* Find Lab color space transformation code, try that out
* Render image in rectanglular windows, recomputing all the triangles that fall within the window, for efficiency
  * This is done, but I need to rethink the Frame abstraction, as I need to have access to CombinedScene with type assertion to make it happen, which isn't great
  * Is there some way to create a mapping from pixel space to Window index? If all Windows were uniform, this would be trivial, but they're of various sizes. Maybe storing them as a tree could help? But lookup would be O(log(n))
* Investiate if there are any optimizations using GPU/CUDA
* Add spring-loaded-mass interactions
  * See https://www.desmos.com/calculator/k01p40v0ct
  * `y=P0*e^{\alpha x}\cos\left(\beta x\right)+C6e^{\alpha x}\sin\left(\beta x\right)\left\{0<x<t\right\}`
    * P0 is initial position, C6 is the velocity component (?), b is friction component, k is spring constant, etc
* Need to refactor how textures and scenes have very similar signatures, but operate on different ranges (0;1) and (-1;1), could one be reused for the other? Do we need to retain the distinction?


## Further Reading:
* https://graphics.stanford.edu/courses/cs348a-09-fall/Handouts/handout15.pdf
* https://en.wikipedia.org/wiki/Homogeneous_coordinates
* https://en.wikipedia.org/wiki/Rotation_matrix
* https://en.wikipedia.org/wiki/Euler_angles#Conventions_by_intrinsic_rotations
* https://www.reddit.com/r/mobilevrstation/comments/xgamsr/typical_ios_shutteringnonplayback_fixes/


## Benchmarks:
Rendering the `SpinningMulticube` scene with the `-video=intermediate` preset takes this long:
* 1m20s - base
* 50s - by caching triangle intermediate results

The same scene with `-video=hidef` (approx 64x the effort, at 4x resolution and 4x anti-aliasing, 4x more frames) takes:
* 1hr - base
* 45m - by caching intermediate triangle results
* 30s - by using windowed triangle culling (very impressive)


# Frequently Asked Questions
* How do I run pprof?
  * Set pprof=true in main, then after the run do  `go tool pprof -png . cpu.pprof` or `go tool pprof -png . mem.pprof`