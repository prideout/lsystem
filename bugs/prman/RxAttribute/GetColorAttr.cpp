#include <stdio.h>
#include <ri.h>
#include <rx.h>
#include <RslPlugin.h>

extern "C" {

RSLEXPORT int getAttr( RslContext *rctx, int /*argc*/, const RslArg* argv[] )
{
    RslStringIter name(argv[1]);  	
    RslColorIter result(argv[0]);
    RtColor value;
    RxInfoType_t rxQueryType;
    int rxQueryCount;
    int status = RxAttribute(*name, 
                             value, sizeof(value),
                             &rxQueryType, &rxQueryCount);
    // Yellow indicates failure.
    if (status) {
        (*result)[0] = 1.0f;
        (*result)[1] = 1.0f;
        (*result)[2] = 0.25f;
    } else {
        (*result)[0] = value[0];
        (*result)[1] = value[1];
        (*result)[2] = value[2];
    }
    return 0;
} 
 
static RslFunction shadeopFunctions[] = {
    {"uniform color GetColorAttr(string)", getAttr, NULL, NULL },
    NULL
};

RSLEXPORT RslFunctionTable RslPublicFunctions(shadeopFunctions);
    
}
