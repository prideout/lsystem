#include <vmath.hpp>
#include <list>

class lsystem
{
public:
    typedef std::list<vmath::Matrix4> XformList;
    static XformList Evaluate(const char* filename, int seed = 0);
};

