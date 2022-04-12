components = require "components"

component = components:NewComponent "greetings"

component:Type       "library"
component:Languages  "CPP"
component:AddSources "src/"

component:ExportedHeaders {
  ["export/[DIRS]/*.hpp"] = "greetings/[DIRS]/*.hpp",
}
