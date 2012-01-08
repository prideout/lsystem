#include <iostream>
#include <lsystem.hpp>
#include <ri.h>

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
    _DrawScene();
    _ReportProgress();
    RiEnd();
}
