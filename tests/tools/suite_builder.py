import ast
import os.path


class TestSuiteWrapper:
    def __init__(self, testSuite, testsToRun):
        self.__testSuite = testSuite
        self.__testsToRun = testsToRun

    def RunTests(self):
        for test in self.__testsToRun:
            if not self.__testSuite.Run(test):
                return False
        return True


def reconcileImports(imports):
    for i in imports:
        if isinstance(i, ast.ImportFrom) and i.module == "test_suite":
            i.module = "tools.test_suite"


def buildTestSuites(filenames, bspath):
    for filename in filenames:
        with open(filename) as file:
            node = ast.parse(file.read())
        classes = [n for n in node.body if isinstance(n, ast.ClassDef)]
        imports = [
            n
            for n in node.body
            if (isinstance(n, ast.Import) or isinstance(n, ast.ImportFrom))
        ]
        reconcileImports(imports)
        ts = []
        for c in classes:
            for b in c.bases:
                if b.id == "TestSuite":
                    ts += [c]
                    break
        for t in ts:
            a = imports + [t]
            code = ast.unparse(a)
            loc = {}
            exec(compile(code, filename, "exec"), {}, loc)
            methods = [n for n in t.body if isinstance(n, ast.FunctionDef)]
            testsToRun = []
            for m in methods:
                if m.name.startswith("Test"):
                    testsToRun += [m.name]
            yield TestSuiteWrapper(
                loc[t.name](t.name, os.path.dirname(filename), bspath),
                testsToRun,
            )


# for tsw in buildTestSuites(["toto.py"]):
#     tsw.RunTests()
