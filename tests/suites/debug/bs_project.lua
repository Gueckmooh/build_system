-- The minimal configuration that shouls be handled by the first
-- version of the build system.
version "0.1.0"

project = require "project"

project:Name    "My Pretty Project"
project:Version "0.0.1"

project:Languages     "CPP"     -- Enables C++ compilation

project:AddSources "sources/"

project:DefaultTarget "hello_exe"

debugProfile = project:Profile "Debug"

CPP = debugProfile:CPP()

CPP:AddBuildOptions {"-g", "-O0"}
