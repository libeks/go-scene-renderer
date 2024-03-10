# go-scene-renderer
Rendering gifs in golang


`go run ./... gif /Users/janis.libeks/own/go-scene-renderer/gallery/sin.gif`

Example:
![alt text](https://github.com/libeks/go-scene-renderer/blob/main/gallery/sin.gif)

## TODOs:
* Flip coordinates so a scene is rendered in usual x,y plot way
* Add more intense gradients, like Bezier, etc
* Find Lab color space transformation code, try that out
* Add Perlin noise-based scenes
* Be smarter about palette building based on how many pixels will need to change
	* Consider the distributions of gradient values in frame