#include <iostream>
#include <lsystem.hpp>

int main()
{
    lsystem::XformList ribbon;
    ribbon = lsystem::Evaluate("Ribbon.xml");
    std::cout << "Success!\n";
}
