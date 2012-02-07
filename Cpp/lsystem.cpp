#include <vector>
#include <sstream>
#include <lsystem.hpp>
#include <diagnostic.hpp>
#include <stdarg.h>
#include <stdio.h>
#include <stdlib.h>
#include <iostream>
#include <map>
#include <tinythread.hpp>

namespace diag = diagnostic;
using std::string;
using namespace vmath;
using namespace tthread;

static float _radians(float degrees)
{
    return degrees * 3.1415926535f / 180.0f;
}

/// Parse a string in the xform language and generate a 4x4 matrix.
/// Examples:
///   "rx -2 tx 0.1 sa 0.996"
///   "s 0.55 2.0 1.25"
static Matrix4 _ParseTransform(const string& xformString)
{
    typedef std::map<string, Matrix4> Cache;
    static Cache cache;
    static mutex cacheMutex;
    {
        lock_guard<mutex> guard(cacheMutex);
        Cache::iterator i = cache.find(xformString);
        if (i != cache.end()) {
            return i->second;
        }
    }
    
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
            xform *= Matrix4::rotationX(_radians(x));
        } else if (token == "ry") {
            s >> y;
            xform *= Matrix4::rotationY(_radians(y));
        } else if (token == "rz") {
            s >> z;
            xform *= Matrix4::rotationZ(_radians(z));
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
    lock_guard<mutex> guard(cacheMutex);
    cache[xformString] = xform;
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
static pugi::xml_node _PickRule(const pugi::xml_document& doc, const char* name,
                                unsigned int* seed)
{
    string path = FormatString("/rules/rule[@name='%s']", name);
    pugi::xpath_node_set rules = doc.select_nodes(path.c_str());
    pugi::xpath_node_set::const_iterator pRule = rules.begin();
    int sum = 0;
    for (; pRule != rules.end(); ++pRule) {
        int weight = pRule->node().attribute("weight").as_int();
        sum += weight;
    }
    diag::Check(sum > 0, "Unable to find rule %s\n", name);
    
    int n = rand_r(seed) % sum;
    for (pRule = rules.begin(); pRule != rules.end(); ++pRule) {
        int weight = pRule->node().attribute("weight").as_int();
        if (n < weight) {
            return pRule->node();
        }
        n -= weight;
    }
    
    throw diag::Unreachable();
}

void _ProcessRule(void* arg)
{
    lsystem::ThreadInfo* info = reinterpret_cast<lsystem::ThreadInfo*>(arg);
    info->Self->_ProcessRule(info->Entry, &info->Result, info->Seed, info->StackOffset);
}

void lsystem::_ProcessRule(const StackEntry& entry, Curve* result,
                           unsigned int seed, unsigned int stackOffset)
{
    {
        lock_guard<mutex> guard(_mutexCompute);
        while (_threadsCompute >= _maxThreads) {
            _condCompute.wait(_mutexCompute);
        }
        ++_threadsCompute;
    }

    Stack stack;
    stack.push_back(entry);

    // Build the list of transforms    
    size_t progressCount = 0;
    while (!stack.empty()) {
        StackEntry entry = stack.back();
        stack.pop_back();
        
        int local_max = _maxDepth;
        if (entry.Node.attribute("max_depth")) {
            local_max = entry.Node.attribute("max_depth").as_int();
        }
        
        // Mark the end of the ribbon due to global depth constraint
        if (stack.size() + stackOffset >= _maxDepth) {
            result->push_back(0);
            continue;
        }
    
        // Mark the end of the ribbon due to local depth constaint
        Matrix4 matrix = entry.Transform;
        if (entry.Depth >= local_max) {
            // Switch to a different rule is one is specified
            if (entry.Node.attribute("successor")) {
                const char* successor = entry.Node.attribute("successor").value();
                pugi::xml_node rule = _PickRule(_doc, successor, &seed);
                StackEntry newEntry = {rule, 0, matrix};
                stack.push_back(newEntry);
            }
            result->push_back(0);
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
                    pugi::xml_node newRule = _PickRule(_doc, rule, &seed);
                    StackEntry newEntry = { newRule, entry.Depth + 1, matrix };
                    stack.push_back(newEntry);
                } else if (tag == "instance") {
                    CurvePoint* cp = new CurvePoint;
                    Vector4 v = matrix * Vector4(0,0,0,1);
                    cp->P = Point3(v.getXYZ());
                    cp->N = matrix.getUpper3x3() * Vector3(0, 0, 1);
                    result->push_back(cp);
                    
                    // Give users a crude progress indicator by maintaining a count
                    // Note that size() on std::list can have linear complexity (!)
                    if (++_curveLength >= progressCount + 10000) {
                        diag::Print("%d curve segments so far...\n", _curveLength);
                        progressCount = _curveLength;
                    }

                } else {
                    op.print(std::cout);
                    diag::Print("Unrecognized operation: '%s'\n", tag.c_str());
                }
            }
        }
    }
    
    lock_guard<mutex> guard1(_mutexCompute);
    lock_guard<mutex> guard2(_mutexComplete);
    _threadsCompute--;
    _threadsComplete++;
    _condCompute.notify_all();
    _condComplete.notify_all();
}

/// Parse given XML file, evaluate rules, populate member curve
void lsystem::Evaluate(const char* filename, unsigned int seed)
{
    // Parse the XML document
    diag::Print("Reading %s...\n", filename);
    pugi::xml_parse_result result = _doc.load_file(filename);
    _AssignDefaults(_doc);
    diag::Check(result, "Unable to read XML file %s: %s\n", filename, result.description());
    
    diag::Print("Evaluating rules...\n");
    
    pugi::xml_node startRule = _PickRule(_doc, "entry", &seed);
    StackEntry entry = { startRule, 0, Matrix4::identity() };
    
    // Extract a maximum recursion depth from the XML
    _curveLength = 0;
    _maxDepth = _doc.first_child().attribute("max_depth").as_int();
    diag::Print("Max depth is %d\n", _maxDepth);

    // Tip for Mac developers: open "Activity Monitor", then right click its dock icon.
    // Select "Show CPU Usage" and keep it in the doc.
    _threadsComplete = 0;
    _threadsCompute = 0;
    _maxThreads = thread::hardware_concurrency();
    std::vector<thread *> threads;

    // Check for a multithread-amenable lsystem    
    pugi::xml_node op = entry.Node.first_child();
    int count = op.attribute("count").as_int();
    string tag = op.name();

    if (_isThreading && count > 1 && tag == "call" && !op.next_sibling()) {
        string xformString = op.attribute("transforms").value();
        Matrix4 xform = _ParseTransform(xformString);
        const char* rule = op.attribute("rule").value();
        
        std::vector<ThreadInfo> curves(count); // TODO remove this by deriving a new class from tthread::thread
        std::vector<ThreadInfo>::iterator pCurve = curves.begin();
        Matrix4 matrix = entry.Transform;

        while (count--) {
            matrix *= xform;
            ThreadInfo& curve = *pCurve++;
            pugi::xml_node newRule = _PickRule(_doc, rule, &seed);
            StackEntry newEntry = { newRule, entry.Depth + 1, matrix };
            curve.Entry = newEntry;
            curve.Self = this;
            curve.Seed = seed;
            curve.StackOffset = count;
            threads.push_back(new thread(::_ProcessRule, &curve));
        }
        
        lock_guard<mutex> guard(_mutexComplete);
        while (_threadsComplete < threads.size()) {
            _condComplete.wait(_mutexComplete);
        }
        
        for (pCurve = curves.begin(); pCurve != curves.end(); ++pCurve) {
            _curve.splice(_curve.begin(), pCurve->Result);
        }

        std::vector<thread *>::iterator pThread = threads.begin();
        for (; pThread != threads.end(); ++pThread) {
            delete *pThread;
        }

    } else {
        this->_ProcessRule(entry, &_curve, seed, 0);
    }
    
    diag::Print("%d total curve segments.\n", _curveLength);
}

lsystem::~lsystem()
{
    for (Curve::iterator i = _curve.begin(); i != _curve.end(); ++i) {
        delete *i;
    }
}

