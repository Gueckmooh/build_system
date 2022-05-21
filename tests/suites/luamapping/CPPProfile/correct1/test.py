from test_suite import TestSuite, assertReturnOk


class CPPProfileCorrect1(TestSuite):
    def TestBuild(self):
        with self.sandbox() as s:
            self.runBS(["build"]).mustBeOk()

    def TestBuildVerbose(self):
        with self.sandbox() as s:
            self.runBS(["build", "--verbose"]).mustBeOk().stdoutMustContain(
                "-DDEBUG", "-O0", "-Wall", "-Werror"
            ).stdoutMustNotMatch(r"hello_exe.*-Wall").stdoutMustContain(
                "-s", "-lm", "-pthread"
            )

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
