#include <iostream>
#include <lsystem.hpp>
#include <ri.h>

#pragma GCC diagnostic ignored "-Wwrite-strings" 

void _DrawScene()
{
}

void _ReportProgress()
{
}

int main()
{
    lsystem ribbon;
    ribbon.Evaluate("Ribbon.xml");
    std::cout << "Success!\n";

    RtContextHandle ri = RiGetContext();
    RiBegin("launch:prman? -t -ctrl $ctrlin $ctrlout");
    RiOption("statistics", "endofframe", 1);
    RiOption("statistics", "xmlfilename", "stats.xml");
    RiDisplay("grasshopper", "framebuffer", "rgba");
    RiDisplayChannel("float _occlusion");
    RiShadingRate(8);
    RiPixelSamples(4, 4);

      /*
    RiProjection(ri.PERSPECTIVE, {ri.FOV: 30})
    RiTranslate(-1, -1, 20)
    RiRotate(-20, 1, 0, 0)
    RiRotate(180, 1, 0, 0)
    curve = [1, 1, 0.8, 0.1, 0.9, 0.2, 1, 1, 1, 1]
    RiCamera("world", {"float[10] shutteropening": curve})
    RiImager("Vignette")
    RiWorldBegin()
    RiDeclare("displaychannels", "string")
    RiDeclare("samples", "float")
    RiDeclare("coordsys", "string")
    RiDeclare("hitgroup", "string")
    RiAttribute("cull", dict(hidden=0,backfacing=0))
    RiAttribute("dice", dict(rasterorient=0))
    RiAttribute("visibility", dict(diffuse=1,specular=1))
*/
    _DrawScene();

    RiWorldEnd();

    _ReportProgress();
    RiEnd();
}
