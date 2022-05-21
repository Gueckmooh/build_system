import tempfile
import shutil
import os.path
import subprocess
from colorama import Fore, Style
from tools.test_assert import AssertError, Asserter
import shellescape
import re


def assertReturnOk(res):
    if res.returncode != 0:
        raise AssertError(
            "Execution failed with error code {}".format(res.returncode)
        )


def getFail():
    return "{}FAIL{}".format(
        Style.BRIGHT + Fore.RED, Style.RESET_ALL + Fore.RESET
    )


def getPass():
    return "{}PASS{}".format(
        Style.BRIGHT + Fore.GREEN, Style.RESET_ALL + Fore.RESET
    )


class CompletedProcessWrapper:
    def __init__(self, res):
        self.__res = res

    def mustBeOk(self):
        print(".", end="", flush=True)
        if self.__res.returncode != 0:
            raise AssertError(
                "Execution failed with error code {} where it should succeed\nstdout:\n{}\nstderr:\n{}".format(
                    self.__res.returncode,
                    self.__res.stdout.decode("utf-8"),
                    self.__res.stderr.decode("utf-8"),
                )
            )
        return self

    def mustBeNOk(self):
        print(".", end="", flush=True)
        if self.__res.returncode == 0:
            raise AssertError(
                "Execution succeeded where it should fail\nstdout:\n{}\nstderr:\n{}".format(
                    self.__res.stdout.decode("utf-8"),
                    self.__res.stderr.decode("utf-8"),
                )
            )
        return self

    def stderrMustContain(self, *cs):
        print(".", end="", flush=True)
        for c in cs:
            if c not in self.__res.stderr.decode("utf-8"):
                raise AssertError(
                    'Could not find "{}" in:\n{}'.format(
                        c, self.__res.stderr.decode("utf-8")
                    )
                )
        return self

    def stdoutMustContain(self, *cs):
        print(".", end="", flush=True)
        for c in cs:
            if c not in self.__res.stdout.decode("utf-8"):
                raise AssertError(
                    'Could not find "{}" in:\n{}'.format(
                        c, self.__res.stdout.decode("utf-8")
                    )
                )
        return self

    def stdoutMustMatch(self, r):
        print(".", end="", flush=True)
        pattern = re.compile(r)
        if pattern.search(self.__res.stdout.decode("utf-8")) is None:
            raise AssertError(
                'Could not find "{}" in:\n{}'.format(
                    r, self.__res.stdout.decode("utf-8")
                )
            )
        return self

    def stdoutMustNotMatch(self, r):
        print(".", end="", flush=True)
        pattern = re.compile(r)
        if pattern.search(self.__res.stdout.decode("utf-8")) is not None:
            raise AssertError(
                'Error "{}" found in:\n{}'.format(
                    r, self.__res.stdout.decode("utf-8")
                )
            )
        return self

    def stdoutMustNotContain(self, *cs):
        print(".", end="", flush=True)
        for c in cs:
            if c in self.__res.stdout.decode("utf-8"):
                raise AssertError(
                    'Error "{}" found in:\n{}'.format(
                        c, self.__res.stdout.decode("utf-8")
                    )
                )
        return self


class TestSuite(Asserter):
    def __init__(self, name: str, d: str, bspath: str):
        self.__name = name
        self.__dir = d
        self.__bspath = bspath
        self.__verbose = False

    def setVerbosity(self, v):
        self.__verbose = v

    def getName(self):
        return self.__name

    def BSPath(self):
        return self.__bspath

    def Run(self, testName: str):
        print(
            "\tRunning test {}{}{}".format(
                Style.BRIGHT, testName, Style.RESET_ALL
            ),
            end="",
            flush=True,
        )
        try:
            getattr(self, testName)()
        except AssertError as e:
            print()
            print(e)
            print("\t" + getFail())
            return False
        print(" " + getPass())
        return True

    def runBS(self, options):
        if self.__verbose:
            print(
                "{}Running command{}: {}".format(
                    Style.BRIGHT,
                    Style.RESET_ALL,
                    " ".join([shellescape.quote(o) for o in options]),
                )
            )
        res = subprocess.run(
            [self.BSPath()] + options,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
        )
        if self.__verbose:
            if len(res.stdout.decode("utf-8")) > 0:
                print("{}stdout{}:".format(Style.BRIGHT, Style.RESET_ALL))
                print(res.stdout.decode("utf-8"))
            if len(res.stderr.decode("utf-8")) > 0:
                print("{}stderr{}:".format(Style.BRIGHT, Style.RESET_ALL))
                print(res.stderr.decode("utf-8"))
        return CompletedProcessWrapper(res)

    def removeFile(self, *filenames):
        for filename in filenames:
            os.remove(filename)

    def runCmd(self, options):
        res = subprocess.run(
            options,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
        )
        return CompletedProcessWrapper(res)

    # def sandbox(self):
    #     return Sandbox()

    class Sandbox:
        def __init__(self, path):
            self.__path = path

        def __enter__(self):
            self.__dir = tempfile.TemporaryDirectory()
            destdir = os.path.join(
                self.__dir.name, os.path.basename(self.__path)
            )
            shutil.copytree(self.__path, destdir)
            self.__cwd = os.getcwd()
            os.chdir(destdir)
            return destdir

        def __exit__(
            self, exception_type, exception_value, exception_traceback
        ):
            os.chdir(self.__cwd)
            self.__dir.cleanup()

    def sandbox(self):
        return self.Sandbox(self.__dir)
