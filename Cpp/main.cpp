#include <iostream>
#include <lsystem.hpp>
#include <ri.h>

#pragma GCC diagnostic ignored "-Wwrite-strings" 

void _DrawWorld()
{
    RiWorldBegin();
    RiDeclare("displaychannels", "string");
    RiDeclare("samples", "float");
    RiDeclare("coordsys", "string");
    RiDeclare("hitgroup", "string");
    
    {
        RtToken tokens[2] = { "int backfacing", "int hidden" };
        RtInt bools[2] = { RI_FALSE, RI_FALSE };
        RtPointer values[2] = {&bools[0], &bools[1]};
        RiAttributeV("cull", 2, tokens, values);
    }
    {
        RtToken tokens[2] = { "int diffuse", "int specular" };
        RtInt bools[2] = { RI_TRUE, RI_TRUE };
        RtPointer values[2] = {bools, bools + 1};
        RiAttributeV("visibility", 2, tokens, values);
    }
    {
        RtInt rasterorient = RI_FALSE;
        RiAttribute("dice", "int rasterorient", &rasterorient, RI_NULL);
    }

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

    {
      RtBoolean endofframe = RI_TRUE;
      const char* xmlfilename = "stats.xml";
      RtToken tokens[2];
      RtPointer values[2];
      tokens[0] = "endofframe";
      values[0] = (RtPointer) &endofframe;
      tokens[1] = "xmlfilename";
      values[1] = (RtPointer) &xmlfilename;
      RiOptionV("statistics", 2, tokens, values);
    }

    RiDisplay("grasshopper", "framebuffer", "rgba", RI_NULL);
    RiDisplayChannel("float _occlusion", RI_NULL);
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
