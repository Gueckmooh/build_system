components = require "components"

component = components:NewComponent "c2"

component:Type       "library"
component:Languages  "CPP"
component:AddSources "src/"

component:Requires "c3"
