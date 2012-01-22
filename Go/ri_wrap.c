#include <stdlib.h>
#include <ri.h>

void curve2(float a, float b)
{
    RtToken type = "cubic";
    RtInt ncurves = 1;
    RtInt nvertices[1] = {4};
    RtPoint P[4] = {{0, 0, 0},  {-1, -.5, 1},  {2, .5, 1},  {1, 0, -1}};
    RtToken wrap = "nonperiodic";
    RtFloat nwidths[2] = {a, b};
    RiDisplay("grasshopper", "framebuffer", "rgba", RI_NULL);
    RiFormat (320, 160, 1);
    RiWorldBegin();
    RiDeclare ("width", "varying float");
    RiSurface("constant",RI_NULL);
    RiCurves (type, ncurves, nvertices, wrap, RI_P, (RtPointer) P,  RI_WIDTH, &nwidths, RI_NULL);
    RiWorldEnd();
}

void curve(float a, float b)
{
    char* launch = "launch:prman? -t -ctrl $ctrlin $ctrlout -capture debug.rib";
    RiBegin(launch);
    curve2(a, b);
    RiEnd();
}

