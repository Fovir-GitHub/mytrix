{
  description = "Devshell for Golang.";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";

  outputs = {
    self,
    nixpkgs,
  }: let
    system = "x86_64-linux";
    pkgs = import nixpkgs {
      inherit system;
      config.permittedInsecurePackages = [
        "olm-3.2.16"
      ];
    };
  in {
    devShells.${system}.default = pkgs.mkShell {
      # Add packages here.
      buildInputs = with pkgs; [
        go
        golangci-lint
        gopls
        gotools
        govulncheck
        just
        litecli
        markdown-toc
        olm
        pre-commit
        prettier
        sqlite
      ];

      # Shell hooks.
      shellHook = ''
        echo "Entering the development environment!"
        go version
      '';
    };
  };
}
