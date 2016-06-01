
load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")


go_library(
    name = "rpc",
    srcs = ["rpc/rpc.go", "rpc/utils.go"],
    deps = ["//vendor/github.com/streadway/amqp:go_default_library"],
    visibility = ["//visibility:public"],
)

go_binary(
    name = "proxy",
    srcs = ["proxy.go"],
    deps = [
        ":rpc",
        "//vendor/github.com/streadway/amqp:go_default_library",
    ],
)

