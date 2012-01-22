package main

import ri "gorman"

func main() {
    
    //Display "grasshopper" "framebuffer" "rgba" 
    //DisplayChannel "float _occlusion" 
    //ShadingRate 4

    launch := "launch:prman? -t -ctrl $ctrlin $ctrlout -capture debug.rib";
    ri.Begin(launch);
    ri.Format(512, 320, 1)
 
    var threadCount int = 2
    ri.Option("limits", "int threads", threadCount)
    ri.Option("statistics", "xmlfilename", "stats.xml")
    
    ri.Temp()
    ri.End();

}