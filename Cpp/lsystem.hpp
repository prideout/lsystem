#include <vmath.hpp>
#include <list>

class lsystem
{
public:
    struct CurvePoint {
        vmath::Point3 P;
        vmath::Vector3 N;
    };
    typedef std::list<CurvePoint*> Curve;
    static Curve Evaluate(const char* filename, int seed = 0);
};

