package main

import ri "gorman"

func main() {

	launch := "launch:prman? -t -ctrl $ctrlin $ctrlout -capture debug.rib"
	ri.Begin(launch)
	ri.Format(512, 320, 1)
	ri.Display("grasshopper", "framebuffer", "rgba")
	ri.ShadingRate(4)

	ri.Option("limits", "int threads", 2)

	ri.Option("statistics",
		"xmlfilename", "stats.xml",
		"endofframe", true)

    ri.PixelSamples(4, 4)
    ri.Projection("perspective", "fov", 30.)
    ri.Translate(0, 0.25, 4)
    ri.Rotate(-20, 1, 0, 0)
    ri.Rotate(180, 1, 0, 0)
    ri.Imager("Vignette")

	ri.Temp()
	ri.End()

}
