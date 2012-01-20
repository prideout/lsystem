#include <iostream>
#include <stdlib.h>
#include <stdio.h>
#include <Ric.h>
#include <string>
#include <vector>

#include "tinythread.hpp"
#include "lsystem.hpp"
#include "diagnostic.hpp"
#include "ru.h"

#pragma GCC diagnostic ignored "-Wwrite-strings" 
const bool HaveLicense = true;

using namespace tthread;
using namespace std;
using namespace vmath;
namespace diag = diagnostic;

static void
_SetLabel(string label, string groups = "")
{
    RuAttributet(RI_IDENTIFIER, RI_NAME, label);
    if (groups.size() > 0) {
        RuAttributet("grouping", RI_GRP_MEMBERSHIP, groups);
    }
}

static void
_DrawCurves(const lsystem& ribbon)
{
    std::vector<float> points;
    std::vector<float> normals;
    std::vector<RtInt> vertsPerCurve;
    
    lsystem::Curve::const_iterator c = ribbon.GetCurve().begin();
    for (RtInt count = 0;; ++c) {
        if (c == ribbon.GetCurve().end() || !*c) {
            if (count == 1) {
                points.pop_back(); normals.pop_back();
                points.pop_back(); normals.pop_back();
                points.pop_back(); normals.pop_back();
            } else if (count > 1) {
                vertsPerCurve.push_back(count);
            }
            count = 0;
            if (c == ribbon.GetCurve().end()) {
                break;
            }
            continue;
        }
        Point3 P = (*c)->P;
        Vector3 N = (*c)->N;
        points.push_back(P.getX()); points.push_back(P.getY());
        points.push_back(P.getZ()); normals.push_back(N.getX());
        normals.push_back(N.getY()); normals.push_back(N.getZ());
        ++count;
    }

    RtInt total = 0;
    std::vector<RtInt>::const_iterator i = vertsPerCurve.begin();
    for (; i != vertsPerCurve.end(); ++i) {
        total += *i;
    }
    diag::Check(total == points.size() / 3, "Incorrect curve topology: %d total verts, %d points\n",
        total, points.size() / 3);

    RtInt ncurves = RtInt(vertsPerCurve.size());
    float curveWidth = 0.025f;
    RiCurves("linear", ncurves, &vertsPerCurve[0],
        "nonperiodic",
        RI_P, &points[0],
        RI_N, &normals[0],
        "constantwidth", &curveWidth,
        RI_NULL);
}

static void
_DrawWorld(const lsystem& ribbon)
{
    RiWorldBegin();
    RiDeclare("displaychannels", "string");
    RiDeclare("samples", "float");
    RiDeclare("coordsys", "string");
    RiDeclare("em", "color");

    RuAttributei("cull", "backfacing", RI_FALSE);
    RuAttributei("cull", "hidden", RI_FALSE);
    RuAttributei("visibility", "diffuse", RI_TRUE);
    RuAttributei("visibility", "specular", RI_TRUE);
    RuAttributei("dice", "rasterorient", RI_FALSE);

    // Create a list of common shader parameters:
    RuMap args;
    args["displaychannels"] = "_occlusion";
    args["samples"] = 128.0f;
    args["filename"] = "";

    // Insert the floor, which fills the background:
    _SetLabel("Floor");
    args["em"] = RuColor::FromBytes(0, 165, 211);
    RuSurface("ComputeOcclusion", args);
    RiTransformBegin();
    RiRotate(90, 1, 0, 0);
    RiDisk(-0.7, 300, 360, RI_NULL);
    RiTransformEnd();
    
    // Insert the L-System art:
    _SetLabel("Sculpture");
    args["em"] = RuColor(1.0f);
    RuSurface("ComputeOcclusion", args);
    RiTransformBegin();
    RiRotate(90, 1, 0, 0);
    RiTranslate(0, 0, -0.55);
    _DrawCurves(ribbon);
    RiTransformEnd();
            
    RiWorldEnd();
}

static void
_ReportProgress()
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

static void
_CompileShader(const char* basename)
{
    std::string cmd = FormatString("shader %s.sl", basename);
    if (system(cmd.c_str()) > 0) {
        diag::Fatal("Failed to compile RSL file %s\n");
    }
}

static void
_InitCamera()
{
    float fov = 30.0f;
    RiProjection("perspective", RI_FOV, &fov, RI_NULL);
    RiTranslate(0, -0.25, 10);
    RiRotate(-20, 1, 0, 0);
    RiRotate(180, 1, 0, 0);
    RiImager("Vignette", RI_NULL);
}

int
main(int argc, char** argv)
{
    bool isThreading = (argc == 1);
    diag::PrintColor(diag::RED, "Evaluating %s L-System...",
                     isThreading ? "multithreaded" : "single-threaded");
    lsystem ribbon(isThreading);
    ribbon.Evaluate("Ribbon.xml");
    std::cout << "Success!\n";

    diag::PrintColor(diag::RED, "Compiling shaders...");
    _CompileShader("Vignette");
    _CompileShader("ComputeOcclusion");
    
    diag::PrintColor(diag::RED, "Rendering...");
    RtContextHandle ri = RiGetContext();
    char* launch = "launch:prman? -t -ctrl $ctrlin $ctrlout -capture debug.rib";
    RiBegin(HaveLicense ? launch : 0);
    RuOptioni("statistics", "endofframe", RI_TRUE);
    RuOptiont("statistics", "xmlfilename", "stats.xml");
    RiFormat(1024, 640, 1.0f);
    RiDisplay("grasshopper", "framebuffer", "rgba", RI_NULL);
    RiDisplayChannel("float _occlusion", RI_NULL);
    RiShadingRate(8);
    RiPixelSamples(4, 4);
    _InitCamera();
    _DrawWorld(ribbon);
    _ReportProgress();
    RiEnd();
}
