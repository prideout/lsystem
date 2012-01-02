#include <vector>
#include <sstream>
#include <lsystem.hpp>
#include <pugixml.hpp>
#include <diagnostic.hpp>
#include <iostream>
#include <tinythread.hpp>

namespace diag = diagnostic;
using std::string;
using namespace vmath;
using namespace tthread;

/// Safe utility function for sprintf-style formatting
/// (light alternative to std::ostringstream and boost::format)
static string _Format(const char* pStr, ...)
{
    va_list a;
    va_start(a, pStr);
    int bytes = 1 + vsnprintf(0, 0, pStr, a);
    char* msg = (char*) malloc(bytes);
    va_start(a, pStr);
    vsnprintf(msg, bytes, pStr, a);
    string retval(msg);
    free(msg);
    return retval;
}

/// Parse a string in the xform language and generate a 4x4 matrix.
/// Examples:
///   "rx -2 tx 0.1 sa 0.996"
///   "s 0.55 2.0 1.25"
static Matrix4 _ParseTransform(const std::string& xformString)
{
    std::istringstream s(xformString);
    Matrix4 xform = Matrix4::identity();
    while (s) {
        std::string token;
        s >> std::skipws >> token;
        float x, y, z;
        if (token == "s") {
            s >> x >> y >> z;
            xform *= Matrix4::scale(Vector3(x, y, z));
        } else if (token == "sa") {
            s >> x;
            xform *= Matrix4::scale(Vector3(x, x, x));
        } else if (token == "t") {
            s >> x >> y >> z;
            xform *= Matrix4::translation(Vector3(x, y, z));
        } else if (token == "tx") {
            s >> x;
            xform *= Matrix4::translation(Vector3(x, 0, 0));
        } else if (token == "ty") {
            s >> y;
            xform *= Matrix4::translation(Vector3(0, y, 0));
        } else if (token == "tz") {
            s >> z;
            xform *= Matrix4::translation(Vector3(0, 0, z));
        } else if (token == "rx") {
            s >> x;
            xform *= Matrix4::rotationX(x);
        } else if (token == "ry") {
            s >> y;
            xform *= Matrix4::rotationY(y);
        } else if (token == "rz") {
            s >> z;
            xform *= Matrix4::rotationZ(z);
        } else if (token == "sx") {
            s >> x;
            xform *= Matrix4::scale(Vector3(x, 1, 1));
        } else if (token == "sy") {
            s >> y;
            xform *= Matrix4::scale(Vector3(1, y, 1));
        } else if (token == "sz") {
            s >> z;
            xform *= Matrix4::scale(Vector3(1, 1, z));
        } else if (token == "") {
            // harmless.
        } else {
            diag::Fatal("Unrecognized xform token: %s\n", token.c_str());
        }
    }
    return xform;
}

/// Add default values to any missing attributes
static void _AssignDefaults(pugi::xml_document& doc)
{
    pugi::xpath_node_set rules = doc.select_nodes("/rules/rule");
    pugi::xpath_node_set::const_iterator pRule = rules.begin();
    for (; pRule != rules.end(); ++pRule) {
        pugi::xml_attribute weight = pRule->node().attribute("weight");
        if (!weight) {
            pRule->node().append_attribute("weight").set_value("1");
        }
        pugi::xml_node op = pRule->node().first_child();
        for (; op; op = op.next_sibling()) {
            pugi::xml_attribute xforms = op.attribute("transforms");
            if (!xforms) {
                op.append_attribute("transforms").set_value("");
            }
            pugi::xml_attribute count = op.attribute("count");
            if (!count) {
                op.append_attribute("count").set_value("1");
            }
        }
    }
}

/// Pick a rule with the given name at random, respecting weights
static pugi::xml_node _PickRule(const pugi::xml_document& doc, const char* name)
{
    string path = _Format("/rules/rule[@name='%s']", name);
    pugi::xpath_node_set rules = doc.select_nodes(path.c_str());
    pugi::xpath_node_set::const_iterator pRule = rules.begin();
    int sum = 0;
    for (; pRule != rules.end(); ++pRule) {
        int weight = pRule->node().attribute("weight").as_int();
        sum += weight;
    }
    diag::Check(sum > 0, "Unable to find rule %s\n", name);
    
    int n = rand() % sum;
    for (pRule = rules.begin(); pRule != rules.end(); ++pRule) {
        int weight = pRule->node().attribute("weight").as_int();
        if (n < weight) {
            return pRule->node();
        }
        n -= weight;
    }
    
    throw diag::Unreachable();
}

struct StackEntry {
    pugi::xml_node Node;
    int Depth;
    Matrix4 Transform;
};

typedef std::list<StackEntry> Stack;

/// Parse given XML file, evaluate rules, populate member curve
void lsystem::Evaluate(const char* filename, int seed)
{
    // Seed the random number generator.
    srand(seed);
    
    // Parse the XML document
    pugi::xml_document doc;
    diag::Print("Reading %s...\n", filename);
    pugi::xml_parse_result result = doc.load_file(filename);
    _AssignDefaults(doc);
    diag::Check(result, "Unable to read XML file %s: %s\n", filename, result.description());
    
    diag::Print("Evaluating L-System...\n");
    
    // Create a small stack for tracking L-System transforms
    Stack stack; {
        pugi::xml_node entry = _PickRule(doc, "entry");
        StackEntry start = { entry, 0, Matrix4::identity() };
        stack.push_back(start);
    }
    
    // Extract a maximum recursion depth from the XML
    int max_depth = doc.first_child().attribute("max_depth").as_int();
    diag::Print("Max depth is %d\n", max_depth);

    // TODO
    // Tip for Mac developers: open "Activity Monitor", then right click its dock icon.
    // Select "Show CPU Usage" and keep it in the doc.
    int maxThreads = thread::hardware_concurrency();
    std::vector<thread *> threads;
    threads.reserve(maxThreads);

    // Build the list of transforms    
    size_t progressCount = 0;
    size_t curveLength = 0;
    while (!stack.empty()) {
        StackEntry entry = stack.back();
        stack.pop_back();
        
        int local_max = max_depth;
        if (entry.Node.attribute("max_depth")) {
            local_max = entry.Node.attribute("max_depth").as_int();
        }
        
        // Mark the end of the ribbon due to global depth constraint
        if (stack.size() >= max_depth) {
            _curve.push_back(0);
            continue;
        }
    
        // Mark the end of the ribbon due to local depth constaint
        Matrix4 matrix = entry.Transform;
        if (entry.Depth >= local_max) {
            // Switch to a different rule is one is specified
            if (entry.Node.attribute("successor")) {
                const char* successor = entry.Node.attribute("successor").value();
                pugi::xml_node rule = _PickRule(doc, successor);
                StackEntry newEntry = {rule, 0, matrix};
                stack.push_back(newEntry);
            }
            _curve.push_back(0);
            continue;
        }
        
        // Iterate through the operations in the rule
        pugi::xml_node op = entry.Node.first_child();
        for (; op; op = op.next_sibling()) {
            string xformString = op.attribute("transforms").value();
            Matrix4 xform = _ParseTransform(xformString);
            int count = op.attribute("count").as_int();
            while (count--) {
                matrix *= xform;
                string tag = op.name();
                if (tag == "call") {
                    const char* rule = op.attribute("rule").value();
                    pugi::xml_node newRule = _PickRule(doc, rule);
                    StackEntry newEntry = { newRule, entry.Depth + 1, matrix };
                    stack.push_back(newEntry);
                } else if (tag == "instance") {
                    CurvePoint* cp = new CurvePoint;
                    cp->P = Point3((matrix * Vector4(0,0,0,1)).getXYZ());
                    cp->N = matrix.getUpper3x3() * Vector3(0, 0, 1);
                    _curve.push_back(cp);
                    
                    // Give users a crude progress indicator by maintaining a count
                    // Note that size() on std::list can have linear complexity (!)
                    if (++curveLength >= progressCount + 10000) {
                        diag::Print("%d curve segments so far...\n", curveLength);
                        progressCount = curveLength;
                    }

                } else {
                    op.print(std::cout);
                    diag::Print("Unrecognized operation: '%s'\n", tag.c_str());
                }
            }
        }
    }
    
    diag::Print("%d total curve segments.\n", curveLength);
}

lsystem::~lsystem()
{
    for (Curve::iterator i = _curve.begin(); i != _curve.end(); ++i) {
        delete *i;
    }
}

