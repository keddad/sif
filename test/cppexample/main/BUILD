load("@rules_cc//cc:defs.bzl", "cc_binary", "cc_library")

filegroup(
    name = "hello-greet-fg",
    srcs = [
        "hello-greet.cc",
        "useless.cc",
    ],
)

cc_library(
    name = "hello-greet",
    srcs = [
        "another-useless.cc",
        ":hello-greet-fg",
    ],
    hdrs = [
        "another-useless.h",
        "hello-greet.h",
    ],
)

cc_library(
    name = "useless",
    srcs = ["useless.cc"],
    hdrs = ["useless.h"],
)

cc_library(
    name = "useless_selected",
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
    name = "hello-world_selects",
    srcs = ["hello-world.cc"],
    deps = [
        ":hello-greet",
        ":useless",
        "//lib:hello-time",
    ] + select(
        {"//conditions:default": [
            ":useless_selected",
        ]},
    ),
)

cc_binary(
    name = "hello-world_selectsonly",
    srcs = ["hello-world.cc"],
    deps = select(
        {"//conditions:default": [
            ":hello-greet",
            ":useless",
            "//lib:hello-time",
        ]},
    ),
)

cc_binary(
    name = "hello-world-nouseless",
    srcs = ["hello-world.cc"],
    deps = [
        ":hello-greet",
        "//lib:hello-time",
    ],
)
