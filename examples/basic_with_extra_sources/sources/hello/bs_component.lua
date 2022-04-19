components = require "components"
path = require "path"
fs = require "fs"

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
component.CPP:AddLinkOptions "-lm"

component:AddPostbuildAction(
  function (targetPath, to)
    base = path.Base(targetPath)
    dir = path.Dir(targetPath)
    newPath = path.Join(dir, "new_" .. base)
    print(string.format("copy %s -> %s", targetPath, newPath))
    fs.CopyFile(targetPath, newPath)
  end
)
