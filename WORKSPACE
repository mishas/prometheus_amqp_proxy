workspace(name = "prometheus_amqp_proxy")


# Go support
http_archive(
    name = "io_bazel_rules_go",
    url = "https://github.com/bazelbuild/rules_go/releases/download/0.10.3/rules_go-0.10.3.tar.gz",
    sha256 = "feba3278c13cde8d67e341a837f69a029f698d7a27ddbb2a202be7a10b22142a",
)

load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains", "go_repository")
go_rules_dependencies()
go_register_toolchains()

go_repository(
    name = "com_github_streadway_amqp",
    commit = "8e4aba63da9fc5571e01c6a45dc809a58cbc5a68",
    importpath = "github.com/streadway/amqp",
)


# Python pip support
git_repository(
    name = "io_bazel_rules_python",
    remote = "https://github.com/bazelbuild/rules_python.git",
    commit = "b25495c47eb7446729a2ed6b1643f573afa47d99",
)

load("@io_bazel_rules_python//python:pip.bzl", "pip_repositories")
pip_repositories()

# This rule translates the specified requirements.txt into @pip_deps//:requirements.bzl,
# which itself exposes a pip_install method.
load("@io_bazel_rules_python//python:pip.bzl", "pip_import")
pip_import(
   name = "pip_deps",
   requirements = "//client/python:requirements.txt",
)

# Load the pip_install symbol for my_deps, and create the dependencies' repositories.
load("@pip_deps//:requirements.bzl", "pip_install")
pip_install()
