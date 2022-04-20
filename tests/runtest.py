from tools.suite_builder import buildTestSuites, TestSuiteWrapper
import inspect, os.path
from pathlib import Path
import sys


def runtests(filenames, bspath):
    for tsw in buildTestSuites(filenames, bspath):
        if not tsw.RunTests():
            return False
    return True


def main():
    filename = inspect.getframeinfo(inspect.currentframe()).filename
    basepath = os.path.dirname(os.path.abspath(filename))
    suitepath = os.path.join(basepath, "suites")

    bspath = os.path.join(os.path.dirname(basepath), "bin/bs")
    if not os.path.exists(bspath):
        print("Error while setup tests: could not find 'bs' bin.")
        print("You must run 'make build' before running the tests.")
        sys.exit(1)

    testfiles = [str(n) for n in Path(suitepath).rglob("test*.py")]

    if not runtests(testfiles, bspath):
        sys.exit(1)


if __name__ == "__main__":
    main()
