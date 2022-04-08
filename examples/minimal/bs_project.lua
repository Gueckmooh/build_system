-- The minimal configuration that shouls be handled by the first
-- version of the build system.
project = require "project"

project:Name    "My Pretty Project"
project:Version "0.0.1"

project:Languages     "CPP"     -- Enables C++ compilation

project:AddSources "./sources/"
