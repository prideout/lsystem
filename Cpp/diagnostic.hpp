#include <exception>
#include <string>

std::string FormatString(const char* pStr, ...);

namespace diagnostic
{
    void Print(const char* pStr, ...);
    void Fatal(const char* pStr, ...);
    void Check(bool condition, ...);
    void PrintColor(const char* prefix, const char* pStr, ...);
    std::exception& Unreachable();

    // https://wiki.archlinux.org/index.php/Color_Bash_Prompt#List_of_colors_for_prompt_and_Bash
    static const char* BLACK   = "\e[1;30m";
    static const char* RED     = "\e[1;31m";
    static const char* GREEN   = "\e[1;32m";
    static const char* YELLOW  = "\e[1;33m";
    static const char* BLUE    = "\e[1;34m";
    static const char* MAGENTA = "\e[1;35m";
    static const char* CYAN    = "\e[1;36m";
    static const char* WHITE   = "\e[1;37m";
    static const char* ANSI_RESET = "\e[0m";
}

