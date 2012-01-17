#include <string>
#include <ri.h>
#include <any.hpp>
#include <map>
#include <diagnostic.hpp>

typedef std::map<std::string, cdiggins::any> RuMap;

struct RuColor {
    RuColor(float r, float g, float b) : r(r), g(g), b(b) {}
    static RuColor FromBytes(unsigned char r, unsigned char g, unsigned char b)
    {
        return RuColor(r/255.0f, g/255.0f, b/255.0f);
    }
    RuColor(float l) : r(l), g(l), b(l) {}
    float r, g, b;
};

inline void RuSurface(const char* shader, RuMap& args)
{
    std::vector<RtToken> tokens(args.size());
    std::vector<RtPointer> params(args.size());
    RuMap::iterator a = args.begin();
    for (int i = 0; a != args.end(); ++a, ++i) {
        tokens[i] = const_cast<char*>(a->first.c_str());
        if (a->second.compatible(1.0f)) {
            params[i] = &(a->second.cast<float>());
        } else if (a->second.compatible("foo")) {
            params[i] = &(a->second.cast<const char*>());
        } else if (a->second.compatible(RuColor(0.0f))) {
            params[i] = &(a->second.cast<RuColor>());
        } else {
            diagnostic::Fatal("Unknown type of %s\n", tokens[i]);
        }
    }
    RiSurfaceV(const_cast<char*>(shader), args.size(), &tokens[0], &params[0]);
}

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
