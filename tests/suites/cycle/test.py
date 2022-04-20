from test_suite import TestSuite, assertReturnOk


class CycleSuite(TestSuite):
    def TestBuild(self):
        with self.sandbox() as s:
            self.runBS(["build"]).mustBeNOk().stdoutMustContain(
                "Forbidden cyclic component dependencies"
            )

    def TestBuildUpstream(self):
        with self.sandbox() as s:
            self.runBS(
                ["build", "--build-upstream"]
            ).mustBeNOk().stdoutMustContain(
                "Forbidden cyclic component dependencies"
            )
