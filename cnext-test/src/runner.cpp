#include <cnext/test.h>

int main(int argc, char* argv[]) {
    std::string filter;
    
    for (int i = 1; i < argc; i++) {
        std::string arg = argv[i];
        if (arg == "--filter" || arg == "-f") {
            if (i + 1 < argc) {
                filter = argv[++i];
            }
        } else if (arg == "--help" || arg == "-h") {
            std::cout << "Usage: " << argv[0] << " [options]" << std::endl;
            std::cout << "Options:" << std::endl;
            std::cout << "  -f, --filter <pattern>  Filter tests by name" << std::endl;
            std::cout << "  -h, --help              Show this help message" << std::endl;
            return 0;
        }
    }
    
    cnext::test::TestRunner runner;
    return runner.Run(filter);
}
