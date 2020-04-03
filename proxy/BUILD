load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["proxy.go"],
    importpath = "github.com/mishas/prometheus_amqp_proxy/proxy",
    visibility = ["//visibility:private"],
    deps = [
        "//proxy/rpc:go_default_library",
        "@com_github_streadway_amqp//:go_default_library",
    ],
)

go_binary(
    name = "proxy",
    embed = [":go_default_library"],
    visibility = ["//visibility:public"],
)
