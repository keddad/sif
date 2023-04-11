load("@rules_cc//cc:defs.bzl", "cc_binary", "cc_library")

filegroup(
    name = "hello-greet-fg",
    srcs = [
        "hello-greet.cc",
        "useless.cc",
    ]
)

cc_library(
    name = "hello-greet",
    srcs = [":hello-greet-fg", "another-useless.cc"],
    hdrs = ["hello-greet.h", "another-useless.h"],
)

cc_library(
    name = "useless",
    srcs = ["useless.cc"],
    hdrs = ["useless.h"],
)

cc_binary(
    name = "hello-world",
    srcs = ["hello-world.cc"],
    deps = [
        ":hello-greet",
        ":useless",
        "//lib:hello-time",
    ],
)

cc_binary(
    name = "hello-world-nouseless",
    srcs = ["hello-world.cc"],
    deps = [
        ":hello-greet",
        "//lib:hello-time",
    ],
)
