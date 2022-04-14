-- The minimal configuration that shouls be handled by the first
-- version of the build system.
project = require "project"

project:Name    "My Pretty Project"
project:Version "0.0.1"

project:Languages     "CPP"     -- Enables C++ compilation

project:AddSources "sources/"
project:DefaultTarget "hello_exe"

project.CPP:Dialect "CPP20"
project.CPP:AddBuildOptions {"-Wall", "-Werror"}
