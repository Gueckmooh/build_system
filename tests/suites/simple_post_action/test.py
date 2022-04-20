from test_suite import TestSuite


class SimplePostbuildActionSuite(TestSuite):
    def TestBuild(self):
        with self.sandbox() as s:
            self.runBS(["build"]).mustBeOk()
            self.AssertFileExist(".build/bin/new_hello_exe")
            self.AssertFileEqual(
                ".build/bin/hello_exe", ".build/bin/new_hello_exe"
            )
