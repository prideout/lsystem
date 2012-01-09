#include <iostream>
#include <lsystem.hpp>
#include <ri.h>

#pragma GCC diagnostic ignored "-Wwrite-strings" 

void _DrawWorld()
{
    RiWorldBegin();
    /*
    RiDeclare("displaychannels", "string");
    RiDeclare("samples", "float");
    RiDeclare("coordsys", "string");
    RiDeclare("hitgroup", "string");
    RiAttribute("cull", "int backfacing", RI_FALSE);
    RiAttribute("cull", "int hidden", RI_FALSE);
    RiAttribute("dice", "int rasterorient", RI_FALSE);
    RiAttribute("visibility", "int diffuse", RI_TRUE);
    RiAttribute("visibility", "int specular", RI_TRUE);
    */
    {
      RtToken tokens[1];
      RtPointer values[1];
      RtColor em = {0.0/255.0,165.0/255.0,211.0/255.0};
      tokens[0] = "color em";
      values[0] = (RtPointer) em;
      RiSurfaceV("ComputeOcclusion", 1, tokens, values);
    }

    RiTransformBegin();
    RiRotate(90, 1, 0, 0);
    RiDisk(-0.7, 300, 360, RI_NULL);
    RiTransformEnd();
    RiWorldEnd();
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
    RiBegin(0); // "launch:prman? -t -ctrl $ctrlin $ctrlout");
    //RiOption("statistics", "endofframe", 1, RI_NULL);
    //    RiOption("statistics", "xmlfilename", "stats.xml");
    RiDisplay("grasshopper", "framebuffer", "rgba", RI_NULL);
    //RiDisplayChannel("float _occlusion");
    RiShadingRate(8);
    RiPixelSamples(4, 4);
    float fov = 30.0f;
    RiProjection("perspective", RI_FOV, &fov, RI_NULL);
    RiTranslate(-1, -1, 20);
    RiRotate(-20, 1, 0, 0);
    RiRotate(180, 1, 0, 0);
    RiImager("Vignette", RI_NULL);
    _DrawWorld();
    _ReportProgress();
    RiEnd();
}
