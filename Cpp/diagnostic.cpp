#include <diagnostic.hpp>
#include <stdlib.h>
#include <stdio.h>
#include <stdarg.h>

#define countof(A) (sizeof(A) / sizeof(A[0]))
    
/// Safe utility function for sprintf-style formatting
/// (light alternative to std::ostringstream and boost::format)
std::string FormatString(const char* pStr, ...)
{
    va_list a;
    va_start(a, pStr);
    int bytes = 1 + vsnprintf(0, 0, pStr, a);
    char* msg = (char*) malloc(bytes);
    va_start(a, pStr);
    vsnprintf(msg, bytes, pStr, a);
    std::string retval(msg);
    free(msg);
    return retval;
}

void diagnostic::Print(const char* pStr, ...)
{
    va_list a;
    va_start(a, pStr);
    int bytes = 1 + vsnprintf(0, 0, pStr, a);
    char* msg = (char*) malloc(bytes);
    va_start(a, pStr);
    vsnprintf(msg, bytes, pStr, a);
    fputs(msg, stdout);
    free(msg);
}

void diagnostic::PrintColor(const char* prefix, const char* pStr, ...)
{
    va_list a;
    va_start(a, pStr);
    int bytes = 1 + vsnprintf(0, 0, pStr, a);
    char* msg = (char*) malloc(bytes);
    va_start(a, pStr);
    vsnprintf(msg, bytes, pStr, a);
    fputs(prefix, stdout);
    fputs(msg, stdout);
    puts(ANSI_RESET);
    free(msg);
}

void diagnostic::Fatal(const char* pStr, ...)
{
    va_list a;
    va_start(a, pStr);
    int bytes = 1 + vsnprintf(0, 0, pStr, a);
    char* msg = (char*) malloc(bytes);
    va_start(a, pStr);
    vsnprintf(msg, bytes, pStr, a);
    fputs(msg, stderr);
    free(msg);
    exit(1);
}

void diagnostic::Check(bool condition, ...)
{
    if (condition)
        return;

    va_list a;
    va_start(a, condition);
    const char* pStr = va_arg(a, const char*);
    int bytes = 1 + vsnprintf(0, 0, pStr, a);
    char* msg = (char*) malloc(bytes);
    va_start(a, condition);
    pStr = va_arg(a, const char*);
    vsnprintf(msg, bytes, pStr, a);
    fputs(msg, stderr);
    free(msg);
    exit(1);
}

class _UnreachableException : public std::exception
{   
    virtual const char* what() const throw() {
        return "Internal error: unreachable code";
    }
};

std::exception& diagnostic::Unreachable()
{
    static _UnreachableException e;
    Print("Internal error: unreachable code.\n");
    return e;
}
