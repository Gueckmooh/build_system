from test_suite import TestSuite, assertReturnOk


class MinimalWithDebugSuite(TestSuite):
    def TestBuild(self):
        with self.sandbox() as s:
            self.runBS(["build"]).mustBeOk()

    def TestBuildAndClean(self):
        with self.sandbox() as s:
            self.runBS(["build"]).mustBeOk()
            self.runBS(["clean"]).mustBeOk()

    def TestClean(self):
        with self.sandbox() as s:
            self.runBS(["clean"]).mustBeOk()

    def TestBuildCleanRebuild(self):
        with self.sandbox() as s:
            self.runBS(["build"]).mustBeOk()
            self.runBS(["clean"]).mustBeOk()
            self.runBS(["build"]).mustBeOk()

    def TestBuildDebug(self):
        with self.sandbox() as s:
            self.runBS(["build", "-p", "Debug"]).mustBeOk()

    def TestBuildAndCleanDebug(self):
        with self.sandbox() as s:
            self.runBS(["build", "-p", "Debug"]).mustBeOk()
            self.runBS(["clean"]).mustBeOk()

    def TestBuildCleanRebuildDebug(self):
        with self.sandbox() as s:
            self.runBS(["build", "-p", "Debug"]).mustBeOk()
            self.runBS(["clean"]).mustBeOk()
            self.runBS(["build", "-p", "Debug"]).mustBeOk()
