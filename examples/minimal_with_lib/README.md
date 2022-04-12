# Minimal example

This is the minimal example that should be supported by the first
version of the build system.

## Basic project configuration

This config should enable to give the project a name, a version, a
language (only C++ for the moment) and a source directory.

### Configure the project

    project:Name      "project name"
    project:Version   "0.0.1"
    projcet:Languages "CPP"

The two lines above gives the basic configuration of the project, its
name, its version and its languages. Note that we use `CPP` and not
`C++` to denote C++ language to have consistency in the latter namings.

### Add sources directories

    project:AddSources "./sources/"

Adding the `"./sources/"` directory as the sources directory states
that all the subdirectories of `sources` may contain a
component. These components are denoted by the presence of a
`bs_component.lua` file.

Other kind of paths should be matched [not implemented yet] such as:
- ./sources/** which would be equivalent to ./sources/
- ./sources/* which would only match directories under sources
  directory, not recurcively.
- ./sources which would only add sources directory as a possible
  component.
