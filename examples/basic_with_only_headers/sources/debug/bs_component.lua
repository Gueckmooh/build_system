components = require "components"

component = components:NewComponent "debug"

component:Type       "headers"
component:Languages  "CPP"

component:ExportedHeaders {
  ["src/[DIRS]/*.hpp"] = "debug/[DIRS]/*.hpp",
}
