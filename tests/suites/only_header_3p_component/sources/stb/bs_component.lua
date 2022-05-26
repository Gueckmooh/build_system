components = require "components"
fs = require "fs"
path = require "path"

component = components:NewComponent "stb"

component:Type       "headers"
component:Languages  "CPP"

function clone_stb(componentPath)
  repoPath = path.Join(componentPath, "stb")
  if not fs.Exists(repoPath) then
    print("Cloning stb repository...")
    g = GitRepository.new {
      url = "https://github.com/nothings/stb",
      revision = "af1a5bc352164740c1cc1354942b1c6b72eacb8a",
      path = repoPath,
    }
    g:CloneAndCheckout()
  end
end

component:AddPrebuildAction(clone_stb)

component:ExportedHeaders {
  ["stb/[DIRS]/*.h"] = "stb/[DIRS]/*.h",
}
