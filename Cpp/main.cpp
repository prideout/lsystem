#include <iostream>
#include <lsystem.hpp>

int main()
{
    lsystem ribbon;
    ribbon.Evaluate("Ribbon.xml");
    std::cout << "Success!\n";
}
