# go-scene-renderer
Rendering visual scenes in golang

`go run ./... -video test mp4 out.mp4`

Example (converted from mp4 to gif for illustrative purposes):

![alt text](https://github.com/libeks/go-scene-renderer/blob/main/gallery/cube_sine.gif)

## TODOs:
* Add more intense gradients, like Bezier, etc
* Find Lab color space transformation code, try that out
* Add Perlin noise-based scenes
* Add wireframe rendering
* Increase anti-aliasing near object edges (adaptive anti-alias)
* Add texture mapping
* Generalize color application to unit cube
* Add render options into command parameters instead of being in code


## Further Reading:
* https://graphics.stanford.edu/courses/cs348a-09-fall/Handouts/handout15.pdf
* https://en.wikipedia.org/wiki/Homogeneous_coordinates
* https://en.wikipedia.org/wiki/Rotation_matrix
* https://en.wikipedia.org/wiki/Euler_angles#Conventions_by_intrinsic_rotations
* https://www.reddit.com/r/mobilevrstation/comments/xgamsr/typical_ios_shutteringnonplayback_fixes/
