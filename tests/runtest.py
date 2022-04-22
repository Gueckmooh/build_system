from tools.suite_builder import buildTestSuites, TestSuiteWrapper
import inspect, os.path
from pathlib import Path
import sys
import argparse


def runtests(filenames, bspath, suite=None, verbose=False):
    for tsw in buildTestSuites(filenames, bspath):
        if verbose:
            tsw.setVerbosity(True)
        if suite is None or tsw.getSuiteName() == suite:
            if not tsw.RunTests():
                return False
    return True


def listsuites(filenames, bspath):
    for tsw in buildTestSuites(filenames, bspath):
        print(tsw.getSuiteName())


def main():
    parser = argparse.ArgumentParser(description="Run test suites")
    parser.add_argument(
        "--suite",
        dest="suite",
        action="store",
        type=str,
        default=None,
        help="Select a test suite to run.",
    )
    parser.add_argument(
        "--list-suites",
        dest="do_list_suites",
        action="store_true",
        default=False,
        help="List all test suites.",
    )
    parser.add_argument(
        "--verbose",
        dest="verbose",
        action="store_true",
        default=False,
        help="Print more informations.",
    )
    args = parser.parse_args()

    filename = inspect.getframeinfo(inspect.currentframe()).filename
    basepath = os.path.dirname(os.path.abspath(filename))
    suitepath = os.path.join(basepath, "suites")

    bspath = os.path.join(os.path.dirname(basepath), "bin/bs")
    if not os.path.exists(bspath):
        print("Error while setup tests: could not find 'bs' bin.")
        print("You must run 'make build' before running the tests.")
        sys.exit(1)

    testfiles = [str(n) for n in Path(suitepath).rglob("test*.py")]

    if args.do_list_suites:
        listsuites(testfiles, bspath)
        sys.exit(0)

    if not runtests(testfiles, bspath, args.suite, args.verbose):
        sys.exit(1)


if __name__ == "__main__":
    main()
