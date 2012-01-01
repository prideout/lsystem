#include <exception>

namespace diagnostic
{
    void Print(const char* pStr, ...);
    void Fatal(const char* pStr, ...);
    void Check(bool condition, ...);
    std::exception& Unreachable();
}

