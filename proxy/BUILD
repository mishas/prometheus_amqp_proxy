
load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")


go_library(
    name = "rpc",
    srcs = ["rpc/rpc.go", "rpc/utils.go"],
    deps = ["@com_github_streadway_amqp//:go_default_library"],
    visibility = ["//visibility:public"],
)

go_binary(
    name = "proxy",
    srcs = ["proxy.go"],
    deps = [
        ":rpc",
        "@com_github_streadway_amqp//:go_default_library",
    ],
)

