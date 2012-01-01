#include <lsystem.hpp>
#include <pugixml.hpp>
#include <diagnostic.hpp>
#include <iostream>

namespace diag = diagnostic;
using std::string;
using namespace vmath;

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
    }
    throw diag::Unreachable();
}

struct StackEntry {
    pugi::xml_node Node;
    int Depth;
    vmath::Matrix4 Transform;
};

typedef std::list<StackEntry> Stack;

lsystem::XformList lsystem::Evaluate(const char* filename, int seed)
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
    int max_depth = doc.attribute("max_depth").as_int();
    while (!stack.empty()) {
        StackEntry entry = stack.back();
        stack.pop_back();
    }
    
    XformList retval;
    return retval;
}
