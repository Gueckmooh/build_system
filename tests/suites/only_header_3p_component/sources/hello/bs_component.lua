components = require "components"

component = components:NewComponent "hello_exe"

component:Type       "executable"
component:Languages  "CPP"
component:AddSources "src/"

component:CPP():AddBuildOptions "-DBUILD"
profileDebug = component:Profile "Debug"
profileDebug:CPP():AddBuildOptions "-DDEBUG"

component:Requires {
  "stb",
}
