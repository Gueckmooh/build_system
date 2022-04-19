# Build System
> Yet another C++ build system heavily inspired by
> [dune](https://dune.build/) and
> [premmake](https://premake.github.io/). A project is defined by a
> `bs_project.lua` file. And each software component is defined by a
> `bs_component.lua` file.

## Basic configuration
A project is defined by a `bs_project.lua` file which is placed at the
root directory of the project. A software component should contain a
`bs_component.lua` file in its root directory.

Here is the basic configuration for this example project:

```
my_project
├── bs_project.lua
└── sources
    └── hello
        ├── bs_component.lua
        └── src
            └── main.cpp
```

### Project configuration

The bare minimum configuration required a project **name**, a project
**version**, the **languages** in which the project is written, and
the path to the **path** to directory containing software components.

Here is the syntax of the `bs_project.lua` file:
```lua
project = require "project"

project:Name       "My Pretty Project"
project:Version    "0.0.1"

project:Languages  "CPP"     -- Enables C++ compilation

project:AddSources "sources/" -- All the directories and subdirectories
                              -- of "sources/" could contain a component
```

### Component configuration

The bare minimum configuration for a component requires a component
**name**, a component **type**, the **languages** in which the
component is written and the **location** of its sources.

```lua
components = require "components"

component = components:NewComponent "hello"

component:Type       "executable" -- Tells the build system to build 
                                  -- an executable
component:Languages  "CPP"        -- Enables C++ compilation
component:AddSources "src/"       -- The source files will be searched 
                                  -- in src and its subdirectories
```
