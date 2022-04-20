import tempfile
import glob
import shutil
import os.path
import subprocess


class TestSuite:
    def __init__(self, name: str, d: str, bspath: str):
        self.__name = name
        self.__dir = d
        self.__bspath = bspath

    def BSPath(self):
        return self.__bspath

    def Run(self, suiteName: str):
        print(self.__name, suiteName)
        getattr(self, suiteName)()

    def runBSWithOptions(self, options):
        return subprocess.run([self.BSPath()] + options)

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
