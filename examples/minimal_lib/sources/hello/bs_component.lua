components = require "components"

component = components:NewComponent "hello_lib"

component:Type       "library"
component:Languages  "CPP"
component:AddSources "src/"
