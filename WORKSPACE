
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
git_repository(
    name = "io_bazel_rules_go",
    remote = "https://github.com/bazelbuild/rules_go.git",
    tag = "0.0.2",
)
load("@io_bazel_rules_go//go:def.bzl", "go_repositories")
go_repositories()

