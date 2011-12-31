#include <LSystem.hpp>
#include <pugixml.hpp>
#include <diagnostic.hpp>

namespace diag = Diagnostic;

LSystem::XformList LSystem::Evaluate(const char* filename, int seed)
{
    pugi::xml_document doc;
    pugi::xml_parse_result result = doc.load_file(filename);
    diag::Check(result, "unable to read XML file %s: %s\n", filename, result.description());
    
    XformList retval;
    return retval;
}
