# go-scene-renderer
Rendering visual scenes in golang

`go run ./... -video test mp4 out.mp4`

Example (converted from mp4 to gif for illustrative purposes):

![alt text](https://github.com/libeks/go-scene-renderer/blob/main/gallery/cube_sine.gif)

## TODOs:
* Add more intense gradients, like Bezier, etc
* Find Lab color space transformation code, try that out
* Increase anti-aliasing near object edges (adaptive anti-alias)
* Add texture mapping
* Render image in rectanglular windows, recomputing all the triangles that fall within the window, for efficiency
  * This is done, but I need to rethink the Frame abstraction, as I need to have access to CombinedScene with type assertion to make it happen, which isn't great
* Investiate if there are any optimizations using GPU/CUDA


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
  * Set pprof=true in main, then after the run do  `go tool pprof -png . cpu.pprof`