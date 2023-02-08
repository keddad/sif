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
    srcs = [":hello-greet-fg"],
    hdrs = ["hello-greet.h"],
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
