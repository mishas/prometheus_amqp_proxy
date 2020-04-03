load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "rpc.go",
        "utils.go",
    ],
    importpath = "github.com/mishas/prometheus_amqp_proxy/proxy/rpc",
    visibility = ["//visibility:public"],
    deps = ["@com_github_streadway_amqp//:go_default_library"],
)
