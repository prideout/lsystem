#include <exception>
#include <string>

std::string FormatString(const char* pStr, ...);

namespace diagnostic
{
    void Print(const char* pStr, ...);
    void Fatal(const char* pStr, ...);
    void Check(bool condition, ...);
    std::exception& Unreachable();
}

