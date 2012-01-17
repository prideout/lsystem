#include <string>
#include <ri.h>

#define RU_PROTOTYPE(Func) \
    \
    inline void \
    Ru##Func##t(const char* token, const char* param, std::string value) \
    { \
        const char* v = value.c_str(); \
        std::string s = "string "; \
        s += param; \
        Ri##Func(const_cast<char*>(token), s.c_str(), &v, RI_NULL); \
    } \
    \
    inline void \
    Ru##Func##i(const char* token, const char* param, RtInt value) \
    { \
        std::string s = std::string("int ") + param; \
        Ri##Func(const_cast<char*>(token), s.c_str(), &value, RI_NULL); \
    } \
    \
    inline void \
    Ru##Func##cf(const char* token, const char* param, RtFloat r, RtFloat g, RtFloat b) \
    { \
        std::string s = std::string("color ") + param; \
        RtColor c = {r, g, b}; \
        Ri##Func(const_cast<char*>(token), s.c_str(), &c, RI_NULL); \
    } \
    \
    inline void \
    Ru##Func##cf(const char* token, const char* param, RtFloat l) \
    { \
        std::string s = std::string("color ") + param; \
        RtColor c = {l, l, l}; \
        Ri##Func(const_cast<char*>(token), s.c_str(), &c, RI_NULL); \
    } \
    \
    inline void \
    Ru##Func##ci(const char* token, const char* param, unsigned char r, unsigned char g, unsigned char b) \
    { \
        std::string s = std::string("color ") + param; \
        RtColor c = {float(r)/255.0f, float(g)/255.0f, float(b)/255.0f}; \
        Ri##Func(const_cast<char*>(token), s.c_str(), &c, RI_NULL); \
    }

RU_PROTOTYPE(Attribute)
RU_PROTOTYPE(Option)
RU_PROTOTYPE(Surface)
