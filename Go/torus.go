package main

import ri "gorman"

func main() {

	launch := "launch:prman? -t -ctrl $ctrlin $ctrlout -capture debug.rib"
	ri.Begin(launch)
	ri.Format(512, 320, 1)
	ri.Display("grasshopper", "framebuffer", "rgba")

	//DisplayChannel "float _occlusion" 
	//ShadingRate 4

	ri.Option("limits", "int threads", 2)

	ri.Option("statistics",
		"xmlfilename", "stats.xml",
		"endofframe", 1) // try true

	ri.Temp()
	ri.End()

}
