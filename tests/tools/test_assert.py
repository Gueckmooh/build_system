import os.path
import filecmp


class AssertError(Exception):
    pass


class Asserter:
    def AssertTrue(self, b):
        print(".", end="", flush=True)
        if not b:
            raise AssertError("Expected True got False")

    def AssertFileExist(self, file):
        print(".", end="", flush=True)
        if not os.path.exists(file):
            raise AssertError('File "{}" does not exit'.format(file))

    def AssertFileEqual(self, f1, f2):
        print(".", end="", flush=True)
        if not filecmp.cmp(f1, f2):
            raise AssertError(
                'Files "{}" and "{}" are not equal'.format(f1, f2)
            )
