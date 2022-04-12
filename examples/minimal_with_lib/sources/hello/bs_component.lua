components = require "components"

component = components:NewComponent "hello_exe"

component:Type       "executable"
component:Languages  "CPP"
component:AddSources "src/"
component:Requires {
  "greetings",
}
