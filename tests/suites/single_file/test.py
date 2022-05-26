from test_suite import TestSuite, assertReturnOk


class SingleFileSuite(TestSuite):
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
