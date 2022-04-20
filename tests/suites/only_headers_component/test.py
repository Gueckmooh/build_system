from test_suite import TestSuite, assertReturnOk


class WithOnlyHeaderSuite(TestSuite):
    def TestBuild(self):
        with self.sandbox() as s:
            self.runBS(["build"]).mustBeNOk()

    def TestBuildUpstream(self):
        with self.sandbox() as s:
            self.runBS(["build", "--build-upstream"]).mustBeOk()

    def TestBuildAndClean(self):
        with self.sandbox() as s:
            self.runBS(["build"]).mustBeNOk()
            self.runBS(["clean"]).mustBeOk()

    def TestClean(self):
        with self.sandbox() as s:
            self.runBS(["clean"]).mustBeOk()

    def TestBuildCleanRebuild(self):
        with self.sandbox() as s:
            self.runBS(["build"]).mustBeNOk()
            self.runBS(["clean"]).mustBeOk()
            self.runBS(["build"]).mustBeNOk()

    def TestBuildUpstreamAndClean(self):
        with self.sandbox() as s:
            self.runBS(["build", "--build-upstream"]).mustBeOk()
            self.runBS(["clean"]).mustBeOk()

    def TestBuildUpstreamCleanRebuild(self):
        with self.sandbox() as s:
            self.runBS(["build", "--build-upstream"]).mustBeOk()
            self.runBS(["clean"]).mustBeOk()
            self.runBS(["build", "--build-upstream"]).mustBeOk()
