# Build System
> Yet another C++ build system

## Road to v0.0.1
- [X] Read basic configuration
- [X] Build a basic project, using basic configuration
- [X] Get dependencies to headers
- [X] Build a basic library, using basic configuration
- [X] Compute dependencies between components
- [ ] Check that there is no cycle in dependencies
- [X] Build components and dependencies, then link
- [X] Add build options for components
- [X] Add platforms & profiles
- [X] Add parallel compilation
- [ ] Add pre and post build hooks
- [ ] Add basic integrations tests
- [ ] Do some (massive) cleanup
- [ ] Add documentation

## Road to v0.0.2
- [ ] Generate makefiles to build this basic configuration
- [ ] Look at ninja as a generation back end
- [ ] Generate lua bindings (optional)

## Road to a cleaner v0.0.2
- [ ] Hack gopher lua to have more insights on the parsed lua
- [ ] Find a way to decompile lua functions and convert them to bash
