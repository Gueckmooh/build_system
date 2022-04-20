components = require "components"
path = require "path"
fs = require "fs"

component = components:NewComponent "hello_exe"

component:Type       "executable"
component:Languages  "CPP"
component:AddSources "src/"

component:AddPostbuildAction(
  function (targetPath)
    base = path.Base(targetPath)
    dir = path.Dir(targetPath)
    newPath = path.Join(dir, "new_" .. base)
    fs.CopyFile(targetPath, newPath)
  end
)
