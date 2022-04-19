components = require "components"

component = components:NewComponent "c1"

component:Type       "library"
component:Languages  "CPP"
component:AddSources "src/"

component:Requires "c2"
