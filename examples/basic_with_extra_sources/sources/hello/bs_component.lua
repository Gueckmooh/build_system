components = require "components"

component = components:NewComponent "hello_exe"

component:Type       "executable"
component:Languages  "CPP"
component:AddSources "src/"

linuxPlatform = component:Platform "Linux" -- @todo error if Linux platform not declared
linuxPlatform.CPP:AddBuildOptions "-DPLATFORM_LINUX"

component.CPP:AddBuildOptions "-DBUILD"
profileDebug = component:Profile "Debug"
profileDebug.CPP:AddBuildOptions "-DDEBUG"
profileDebug:AddSources "debug/"
