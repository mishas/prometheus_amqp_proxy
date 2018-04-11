workspace(name = "prometheus_amqp_proxy")

new_git_repository(
    name = "pika_git",
    remote = "https://github.com/pika/pika.git",
    tag = "0.10.0",
    build_file_content = """
py_library(
    name = "pika",
    srcs = glob(["pika/*.py", "pika/**/*.py"]),
    visibility = ["//visibility:public"],
)        
    """,
)
             
bind( 
    name = "pika",
    actual = "@pika_git//:pika"
)   
    
new_git_repository(
    name = "prometheus_client_py_git",
    remote = "https://github.com/prometheus/client_python.git",
    tag = "0.0.13",
    build_file_content = """
py_library(
    name = "prometheus_client",
    srcs = glob(["prometheus_client/*.py"]),
    visibility = ["//visibility:public"],
)
    """,
)

bind(
    name = "prometheus_client_py",
    actual = "@prometheus_client_py_git//:prometheus_client",
)

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
