load("@pip_deps//:requirements.bzl", "requirement")

exports_files(["requirements.txt"])

py_library(
    name = "prometheus_proxy",
    srcs = ["prometheus_proxy.py"],
    deps = [
        "//external:pika",
        "//external:prometheus_client_py",
    ],
    imports = ["."],
    visibility = ["//visibility:public"],
)
