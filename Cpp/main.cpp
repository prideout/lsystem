#include <iostream>
#include <lsystem.hpp>

int main()
{
    lsystem::Curve ribbon;
    ribbon = lsystem::Evaluate("Ribbon.xml");
    std::cout << "Success!\n";
}
