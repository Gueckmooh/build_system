from test_suite import TestSuite, assertReturnOk


class CPPProfileIncorrect1(TestSuite):
    def TestBuild(self):
        with self.sandbox() as s:
            self.runBS(["build"]).mustBeNOk()

    def TestBuildAndClean(self):
        with self.sandbox() as s:
            self.runBS(["build"]).mustBeNOk()
            self.runBS(["clean"]).mustBeNOk()

    def TestClean(self):
        with self.sandbox() as s:
            self.runBS(["clean"]).mustBeNOk()

    def TestBuildCleanRebuild(self):
        with self.sandbox() as s:
            self.runBS(["build"]).mustBeNOk()
            self.runBS(["clean"]).mustBeNOk()
            self.runBS(["build"]).mustBeNOk()
