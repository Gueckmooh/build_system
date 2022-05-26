version "0.0.1"

project = require "project"
components = require "components"

project:Name "My Pretty Project"
project:Version "0.0.1"

project:DefaultTarget "hello"

hello = components:NewComponent "hello"
hello:Type "executable"
hello:Languages "CPP"
hello:AddSources "**.cpp"
