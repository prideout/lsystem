#include <LSystem.hpp>
#include <pugixml.hpp>
#include <diagnostic.hpp>

namespace diag = diagnostic;

lsystem::XformList lsystem::Evaluate(const char* filename, int seed)
{
    pugi::xml_document doc;
    diag::Print("Reading %s...\n", filename);
    pugi::xml_parse_result result = doc.load_file(filename);
    diag::Check(result, "unable to read XML file %s: %s\n", filename, result.description());
    
    diag::Print("Evaluating L-System...\n");

    XformList retval;
    return retval;
}
