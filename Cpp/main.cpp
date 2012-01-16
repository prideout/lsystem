#include <iostream>
#include <stdlib.h>
#include <tinythread.hpp>
#include <lsystem.hpp>
#include <diagnostic.hpp>
#include <ri.h>
#include <Ric.h>
#include <string>

#pragma GCC diagnostic ignored "-Wwrite-strings" 
const bool HaveLicense = true;

using namespace tthread;
using namespace std;
namespace diag = diagnostic;

static void _SetLabel(string label, string groups = "")
{
    const char* pLabel = label.c_str();
    RiAttribute(RI_IDENTIFIER, RI_NAME, &pLabel, RI_NULL);
    if (groups.size() > 0) {
        const char* pGroups = groups.c_str();
        RiAttribute("grouping", RI_GRP_MEMBERSHIP, &pGroups, RI_NULL);
    }
}

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

    _SetLabel("Floor");
    RiTransformBegin();
    RiRotate(90, 1, 0, 0);
    RiDisk(-0.7, 300, 360, RI_NULL);
    RiTransformEnd();
        
    RiWorldEnd();
}

void _ReportProgress()
{
    int previous = 0;
    int progress = 0;
    while (progress < 100) {
        RicProcessCallbacks();
        progress = RicGetProgress();
        if (progress >= 100 || progress < previous)
            break;
        if (progress != previous) {
            printf("\r%3d%%\n", progress);
            previous = progress;
            this_thread::sleep_for(chrono::milliseconds(200));
        }
    }
}

void _CompileShader(const char* basename)
{
    std::string cmd = FormatString("shader %s.sl", basename);
    if (system(cmd.c_str()) > 0) {
        diag::Fatal("Failed to compile RSL file %s\n");
    }
}

int main()
{
    diag::PrintColor(diag::RED, "Evaluating L-System...");

    lsystem ribbon;
    ribbon.Evaluate("Ribbon.xml");
    std::cout << "Success!\n";

    diag::PrintColor(diag::RED, "Compiling shaders...");

    _CompileShader("Vignette");
    _CompileShader("ComputeOcclusion");
    
    diag::PrintColor(diag::RED, "Rendering...");

    RtContextHandle ri = RiGetContext();
    char* launch = "launch:prman? -t -ctrl $ctrlin $ctrlout";
    RiBegin(HaveLicense ? launch : 0);

    {
      RtInt endofframe = RI_TRUE;
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
