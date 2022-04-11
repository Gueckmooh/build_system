components = require "components"

component = components:NewComponent "hello_lib"

component:Type       "library"
component:Languages  "CPP"
component:AddSources "src/"

component:ExportedHeaders {
  ["export/[DIRS]/*.hpp"] = "hello/[DIRS]/*.hpp",
}
