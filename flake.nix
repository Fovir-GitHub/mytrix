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
        gopls
        gotools
        govulncheck
        just
        olm
      ];

      # Shell hooks.
      shellHook = ''
        echo "Entering the development environment!"
        go version
      '';
    };
  };
}
