#include <diagnostic.hpp>
#include <stdlib.h>
#include <stdio.h>
#include <stdarg.h>
#include <signal.h>
#include <wchar.h>

#define countof(A) (sizeof(A) / sizeof(A[0]))

void diagnostic::Print(const char* pStr, ...)
{
    va_list a;
    va_start(a, pStr);

    char msg[1024] = {0};
    vsnprintf(msg, countof(msg), pStr, a);
    fputs(msg, stderr);
}

static void _FatalError(const char* pStr, va_list a)
{
    char msg[1024] = {0};
    vsnprintf(msg, countof(msg), pStr, a);
    fputs(msg, stderr);
    exit(1);
}

void diagnostic::Fatal(const char* pStr, ...)
{
    va_list a;
    va_start(a, pStr);
    _FatalError(pStr, a);
}

void diagnostic::Check(bool condition, ...)
{
    va_list a;
    const char* pStr;

    if (condition)
        return;

    va_start(a, condition);
    pStr = va_arg(a, const char*);
    _FatalError(pStr, a);
}
