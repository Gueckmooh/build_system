-- The minimal configuration that shouls be handled by the first
-- version of the build system.
version "0.1.0"

project = require "project"

project:Name    "My Pretty Project"
project:Version "0.0.1"

project:Languages     "CPP"     -- Enables C++ compilation

project:AddSources "sources/"

project:DefaultTarget "hello_exe"


project:CPP():Dialect "CPP20"

project:CPP():AddBuildOptions(3)
project:CPP():AddBuildOptions {"-DDEBUG", "-Wall", "-Werror"}

project:CPP():AddLinkOptions "-lm"
project:CPP():AddLinkOptions {"-pthread", "-s"}
