# go-scene-renderer
Rendering visual scenes in golang

`go run ./... gif /Users/janis.libeks/own/go-scene-renderer/gallery/sin.gif`

Example:
![alt text](https://github.com/libeks/go-scene-renderer/blob/main/gallery/sin.gif)

## TODOs:
* Flip coordinates so a scene is rendered in usual x,y plot way
* Add more intense gradients, like Bezier, etc
* Find Lab color space transformation code, try that out
* Add Perlin noise-based scenes
* Use progress bar for rendering progress
	* https://github.com/schollz/progressbar
* Figure out why video files cannot be transfered to iPhone. Is it due to a missing audio track?
	* no, empty audio is still rejected

GIFs:
* Be smarter about palette building based on how many pixels will need to change
	* Consider the distributions of gradient values in frame



## Further Reading:
* https://graphics.stanford.edu/courses/cs348a-09-fall/Handouts/handout15.pdf
* https://en.wikipedia.org/wiki/Homogeneous_coordinates
* https://en.wikipedia.org/wiki/Rotation_matrix
* https://en.wikipedia.org/wiki/Euler_angles#Conventions_by_intrinsic_rotations
* https://www.reddit.com/r/mobilevrstation/comments/xgamsr/typical_ios_shutteringnonplayback_fixes/
