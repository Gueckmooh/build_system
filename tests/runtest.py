from tools.suite_builder import buildTestSuites, TestSuiteWrapper
import inspect, os.path
from pathlib import Path


def runtests(filenames, bspath):
    for tsw in buildTestSuites(filenames, bspath):
        tsw.RunTests()


def main():
    filename = inspect.getframeinfo(inspect.currentframe()).filename
    basepath = os.path.dirname(os.path.abspath(filename))
    suitepath = os.path.join(basepath, "suites")

    bspath = os.path.join(os.path.dirname(basepath), "bin/bs")
    print(bspath)

    print("Getting tests suites from {}...".format(suitepath))

    testfiles = [str(n) for n in Path(suitepath).rglob("test*.py")]

    runtests(testfiles, bspath)


if __name__ == "__main__":
    main()
