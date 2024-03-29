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
    shutil.copytree(folder, os.path.join(temp_dir, basename), symlinks=True)

    return temp_dir, os.path.join(temp_dir, basename)


def get_sif():
    binary = os.path.join(os.path.dirname(__file__), "sif")

    if not os.path.exists(binary):
        raise RuntimeError("Can't find sif binary near test.py")

    return binary


class CppOptimizationTests(unittest.TestCase):
    main_build_file = "main/BUILD"
    lib_build_file = "lib/BUILD"

    expected_hello_world = """
cc_binary(
    name = "hello-world",
    srcs = ["hello-world.cc"],
    deps = [
        ":hello-greet",
        "//lib:hello-time",
    ],
)"""

    def setUp(self):
        self.temp_dir, self.workspace = copy_testfiles(CPPEXAMPLE_WORKSPACE)
        self.sif = get_sif()

    def tearDown(self):
        # Otherwise Bazel servers from previous tests keep running, eventually causing OOM
        res = subprocess.run(["bazel", "shutdown"],
                             shell=True, cwd=self.workspace)

        if res.returncode != 0:
            raise RuntimeError("Can't kill Bazel")

        shutil.rmtree(self.temp_dir)

    def test_noopt(self):
        res = subprocess.run([self.sif, "--workspace", self.workspace,
                             "--label", "//main:hello-world-nouseless", "--params", "deps"])

        self.assertEqual(res.returncode, 0)

        self.assertEqual(
            (Path(self.workspace) / self.main_build_file).read_text(),
            (Path(CPPEXAMPLE_WORKSPACE) / self.main_build_file).read_text()
        )

    def test_simpleopt(self):
        res = subprocess.run([self.sif, "--workspace", self.workspace,
                             "--label", "//main:hello-world", "--params", "deps"])

        self.assertEqual(res.returncode, 0)

        self.assertIn(
            self.expected_hello_world, (Path(self.workspace) / self.main_build_file).read_text())

    def test_checkers(self):
        res = subprocess.run([self.sif, "--workspace", self.workspace,
                             "--label", "//main:hello-greet-fg", "--params", "srcs", "--check", "//main:hello-greet,//main:hello-world"])

        self.assertEqual(res.returncode, 0)

        self.assertNotIn(
            "\"useless.cc\",", (Path(self.workspace) / self.main_build_file).read_text())

    def test_multiparam(self):
        res = subprocess.run([self.sif, "--workspace", self.workspace,
                             "--label", "//main:hello-greet", "--params", "srcs,hdrs", "--check", "//main:hello-greet"])

        self.assertEqual(res.returncode, 0)

        self.assertNotIn(
            "\"another-useless.cc\"", (Path(self.workspace) / self.main_build_file).read_text())

        self.assertNotIn(
            "\"another-useless.h\"", (Path(self.workspace) / self.main_build_file).read_text())

    def test_recopt(self):
        res = subprocess.run([self.sif, "--workspace", self.workspace,
                             "--label", "//main:hello-world", "--params", "deps,hdrs,srcs", "--recparams", "deps"])

        self.assertEqual(res.returncode, 0)

        self.assertEqual((Path(self.workspace) / self.lib_build_file).read_text(), (Path(
            CPPEXAMPLE_WORKSPACE) / self.lib_build_file).read_text(), "Lib was changed!")

        self.assertIn(
            self.expected_hello_world, (Path(self.workspace) / self.main_build_file).read_text())

        self.assertNotIn("\"another-useless.cc\"",
                         (Path(self.workspace) / self.main_build_file).read_text())
        self.assertNotIn("\"another-useless.h\"",
                         (Path(self.workspace) / self.main_build_file).read_text())

    def test_recopt_blacklist(self):
        res = subprocess.run([self.sif, "--workspace", self.workspace,
                             "--label", "//main:hello-world", "--params", "deps,hdrs,srcs", "--recparams", "deps", "--recblacklist", "hello-greet"])

        self.assertEqual(res.returncode, 0)

        self.assertIn(
            self.expected_hello_world, (Path(self.workspace) / self.main_build_file).read_text())

        self.assertIn("\"another-useless.cc\"",
                      (Path(self.workspace) / self.main_build_file).read_text())
        self.assertIn("\"another-useless.h\"",
                      (Path(self.workspace) / self.main_build_file).read_text())

    def test_partial_select(self):
        res = subprocess.run([self.sif, "--workspace", self.workspace,
                             "--label", "//main:hello-world_selects", "--params", "deps"])

        self.assertEqual(res.returncode, 0)

        expected_hello_world_selects = """
cc_binary(
    name = "hello-world_selects",
    srcs = ["hello-world.cc"],
    deps = [
        ":hello-greet",
        "//lib:hello-time",
    ] + select(
        {"//conditions:default": [
            ":useless_selected",
        ]},
    ),
)"""

        self.assertIn(
            expected_hello_world_selects, (Path(self.workspace) / self.main_build_file).read_text())

    def test_partial_selectonly(self):
        res = subprocess.run([self.sif, "--workspace", self.workspace,
                             "--label", "//main:hello-world_selectsonly", "--params", "deps"])

        self.assertEqual(res.returncode, 0)

        self.assertEqual(
            (Path(self.workspace) / self.main_build_file).read_text(),
            (Path(CPPEXAMPLE_WORKSPACE) / self.main_build_file).read_text()
        )
