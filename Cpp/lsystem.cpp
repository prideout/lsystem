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
        pugi::xml_node op = pRule->node().first_child();
        for (; op; op = op.next_sibling()) {
            pugi::xml_attribute xforms = op.attribute("transforms");
            if (!xforms) {
                op.append_attribute("transforms").set_value("");
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
    }
    throw diag::Unreachable();
}

struct StackEntry {
    pugi::xml_node Node;
    int Depth;
    Matrix4 Transform;
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
    int max_depth = doc.first_child().attribute("max_depth").as_int();
    diag::Print("Max depth is %d\n", max_depth);

    // Build the list of transforms    
    XformList retval;
    while (!stack.empty()) {
        StackEntry entry = stack.back();
        stack.pop_back();
        
        int local_max = max_depth;
        if (entry.Node.attribute("max_depth")) {
            local_max = entry.Node.attribute("max_depth").as_int();
        }
        
        // Mark the end of the ribbon due to global depth constraint
        if (stack.size() >= max_depth) {
            retval.push_back(0);
            continue;
        }
    
        // Mark the end of the ribbon due to local depth constaint
        Matrix4 xform = entry.Transform;
        if (entry.Depth >= local_max) {
            // Switch to a different rule is one is specified
            if (entry.Node.attribute("successor")) {
                const char* successor = entry.Node.attribute("successor").value();
                pugi::xml_node rule = _PickRule(doc, successor);
                StackEntry newEntry = {rule, 0, xform};
                stack.push_back(newEntry);
            }
            retval.push_back(0);
            continue;
        }
        
        // Iterate through the operations in the rule
        pugi::xml_node op = entry.Node.first_child();
        for (; op; op = op.next_sibling()) {
            string xformString = op.attribute("transforms").value();
            std::cout << "Transforms string: " << xformString.c_str() << std::endl;
        }
    }
    
/*    

        for statement in rule:
            xform = parse_xform(statement.get("transforms", ""))
            count = int(statement.get("count", 1))
            for n in xrange(count):
                matrix *= xform
                if statement.tag == "call":
                    rule = pick_rule(tree, statement.get("rule"))
                    cloned_matrix = matrix.copy()
                    stack.append((rule, depth + 1, cloned_matrix))
                elif statement.tag == "instance":
                    name = statement.get("shape")
                    if name == "curve":
                        P = Point3(0, 0, 0)
                        N = Vector3(0, 0, 1)
                        P = matrix * P
                        N = matrix.upperLeft() * N
                        shapes.append((P, N))
                    else:
                        shape = (name, matrix)
                        shapes.append(shape)
                else:
                    print "malformed xml"
                    quit()
*/
    return retval;
}
