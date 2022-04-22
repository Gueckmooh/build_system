-- The minimal configuration that shouls be handled by the first
-- version of the build system.
version "0.0.1"

project = require "project"

project:Name    "My Pretty Project"
project:Version "0.0.1"

project:Languages     "CPP"     -- Enables C++ compilation

project:AddSources "sources/"

project:DefaultTarget "hello_exe"


project.CPP:Dialect {"CPP20"}

project.CPP:AddBuildOptions "-O0"
project.CPP:AddBuildOptions {"-DDEBUG", "-Wall", "-Werror", 42}

project.CPP:AddLinkOptions "-lm"
project.CPP:AddLinkOptions {"-pthread", "-s"}
