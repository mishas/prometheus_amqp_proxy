workspace(name = "prometheus_amqp_proxy")

load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository")
load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")


# Go support
http_archive(
    name = "io_bazel_rules_go",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.22.2/rules_go-v0.22.2.tar.gz",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.22.2/rules_go-v0.22.2.tar.gz",
    ],
    sha256 = "142dd33e38b563605f0d20e89d9ef9eda0fc3cb539a14be1bdb1350de2eda659",
)

load("@io_bazel_rules_go//go:deps.bzl", "go_rules_dependencies", "go_register_toolchains")
go_rules_dependencies()
go_register_toolchains()

http_archive(
    name = "bazel_gazelle",
    urls = [
        "https://storage.googleapis.com/bazel-mirror/github.com/bazelbuild/bazel-gazelle/releases/download/v0.20.0/bazel-gazelle-v0.20.0.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.20.0/bazel-gazelle-v0.20.0.tar.gz",
    ],
    sha256 = "d8c45ee70ec39a57e7a05e5027c32b1576cc7f16d9dd37135b0eddde45cf1b10",
)

# Load and call Gazelle dependencies
load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies", "go_repository")

gazelle_dependencies()

go_repository(
    name = "com_github_streadway_amqp",
    commit = "8e4aba63da9fc5571e01c6a45dc809a58cbc5a68",
    importpath = "github.com/streadway/amqp",
)


# Python pip support
http_archive(
    name = "rules_python",
    url = "https://github.com/bazelbuild/rules_python/releases/download/0.0.1/rules_python-0.0.1.tar.gz",
    sha256 = "aa96a691d3a8177f3215b14b0edc9641787abaaa30363a080165d06ab65e1161",
)
load("@rules_python//python:pip.bzl", "pip_repositories")
pip_repositories()

# This rule translates the specified requirements.txt into @pip_deps//:requirements.bzl,
# which itself exposes a pip_install method.
load("@rules_python//python:pip.bzl", "pip_import")
pip_import(
   name = "pip_deps",
   requirements = "//client/python:requirements.txt",
)

# Load the pip_install symbol for my_deps, and create the dependencies' repositories.
load("@pip_deps//:requirements.bzl", "pip_install", "requirement")
pip_install()

# For backward compatibility
bind(
    name = "pika",
    actual = requirement("pika"),
)
bind(
    name = "prometheus_client_py",
    actual = requirement("prometheus_client"),
)
