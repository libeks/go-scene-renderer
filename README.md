# go-scene-renderer
Rendering visual scenes in golang

`go run ./... mp4 /Users/janis.libeks/own/go-scene-renderer/gallery/sin.mp4`

Example (converted from mp4 to gif for illustrative purposes):

![alt text](https://github.com/libeks/go-scene-renderer/blob/main/gallery/cube_sine.gif)

## TODOs:
* Flip coordinates so a scene is rendered in usual x,y plot way
* Add more intense gradients, like Bezier, etc
* Find Lab color space transformation code, try that out
* Add Perlin noise-based scenes

## Further Reading:
* https://graphics.stanford.edu/courses/cs348a-09-fall/Handouts/handout15.pdf
* https://en.wikipedia.org/wiki/Homogeneous_coordinates
* https://en.wikipedia.org/wiki/Rotation_matrix
* https://en.wikipedia.org/wiki/Euler_angles#Conventions_by_intrinsic_rotations
* https://www.reddit.com/r/mobilevrstation/comments/xgamsr/typical_ios_shutteringnonplayback_fixes/
