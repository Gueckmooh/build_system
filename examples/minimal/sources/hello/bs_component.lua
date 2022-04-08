components = require "components"

component = components:NewComponent "hello_exe"

component:Languages "CPP"
component:AddSources "./src/"
