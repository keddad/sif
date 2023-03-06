import unittest
from pathlib import Path
import os
import subprocess
import tempfile
import shutil

CPPEXAMPLE_WORKSPACE = os.path.join(
    os.path.dirname(__file__), "test", "cppexample")


def copy_testfiles(folder: str):
    temp_dir = tempfile.mkdtemp()
    basename = os.path.basename(folder)
    shutil.copytree(folder, os.path.join(temp_dir, basename))

    return temp_dir, os.path.join(temp_dir, basename)


def get_sif():
    binary = os.path.join(os.path.dirname(__file__), "sif")

    if not os.path.exists(binary):
        raise RuntimeError("Can't find sif binary near test.py")

    return binary


class CppOptimizationTests(unittest.TestCase):
    build_file = "main/BUILD"

    def setUp(self):
        self.temp_dir, self.workspace = copy_testfiles(CPPEXAMPLE_WORKSPACE)
        self.sif = get_sif()

    def tearDown(self):
        shutil.rmtree(self.temp_dir)

    def test_noopt(self):
        res = subprocess.run([self.sif, "--workspace", self.workspace,
                             "--label", "//main:hello-world-nouseless", "--param", "deps"])

        self.assertEqual(res.returncode, 0)

        self.assertEqual(
            (Path(self.workspace) / self.build_file).read_text(),
            (Path(CPPEXAMPLE_WORKSPACE) / self.build_file).read_text()
        )

    def test_simpleopt(self):
        res = subprocess.run([self.sif, "--workspace", self.workspace,
                             "--label", "//main:hello-world", "--param", "deps"])

        self.assertEqual(res.returncode, 0)

        self.assertNotIn(
            ":useless", (Path(self.workspace) / self.build_file).read_text())

    def test_checkers(self):
        res = subprocess.run([self.sif, "--workspace", self.workspace,
                             "--label", "//main:hello-greet-fg", "--param", "srcs", "--check", "//main:hello-greet,//main:hello-world"])
        
        self.assertEqual(res.returncode, 0)

        self.assertNotIn(
            "\"useless.cc\",", (Path(self.workspace) / self.build_file).read_text())