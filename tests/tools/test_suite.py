import tempfile
import glob
import shutil
import os.path
import subprocess
from colorama import Fore, Style


class AssertError(Exception):
    pass


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


class CompletedProcessToto:
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

    def stderrMustContain(self, c):
        print(".", end="", flush=True)
        if c not in self.__res.stderr.decode("utf-8"):
            raise AssertError(
                'Could not find "{}" in:\n{}'.format(
                    c, self.__res.stderr.decode("utf-8")
                )
            )
        return self

    def stdoutMustContain(self, c):
        print(".", end="", flush=True)
        if c not in self.__res.stdout.decode("utf-8"):
            raise AssertError(
                'Could not find "{}" in:\n{}'.format(
                    c, self.__res.stdout.decode("utf-8")
                )
            )
        return self


class TestSuite:
    def __init__(self, name: str, d: str, bspath: str):
        self.__name = name
        self.__dir = d
        self.__bspath = bspath

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
        res = subprocess.run(
            [self.BSPath()] + options,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
        )
        return CompletedProcessToto(res)

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
