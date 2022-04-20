from test_suite import TestSuite


class MinimalSuite(TestSuite):
    def TestBuild(self):
        with self.sandbox() as s:
            res = self.runBSWithOptions(["build"])
            return res.returncode == 0
