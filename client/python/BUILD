load("@pip_deps//:requirements.bzl", "requirement")

exports_files(["requirements.txt"])


py_library(
    name = "prometheus_proxy",
    srcs = ["prometheus_proxy.py"],
    deps = [
        requirement("pika"),
        requirement("prometheus_client"),
    ],
    imports = ["."],
    visibility = ["//visibility:public"],
)

