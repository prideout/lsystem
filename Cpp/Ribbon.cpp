#include <iostream>
#include <LSystem.hpp>

int main()
{
    LSystem::XformList ribbon;
    ribbon = LSystem::Evaluate("Ribbon.xml");
    std::cout << "Success!\n";
}
