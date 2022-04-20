components = require "components"

component = components:NewComponent "c3"

component:Type       "library"
component:Languages  "CPP"
component:AddSources "src/"

component:Requires "c1"
