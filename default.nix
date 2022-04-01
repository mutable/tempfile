{ platform, ... }:

platform.buildGo.package {
  name = "github.com/mutable/tempfile";

  srcs = [
    ./tempfile.go
  ];

  deps = with platform.third_party; [
    gopkgs."golang.org".x.sys.unix
  ];
}
