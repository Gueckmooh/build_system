#include <iostream>
#include <math.h>

#include <debug/debug.hpp>

int main(void) {
    std::cout << "-> " << fmod(42, 42) << std::endl;
    std::cout << "Hello, World!" << std::endl;
    debug("toto");
    return 0;
}
