{
  inputs = {
    flake-parts.url = "github:hercules-ci/flake-parts";
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    git-hooks-nix.url = "github:cachix/git-hooks.nix";
  };

  outputs = inputs @ {flake-parts, ...}:
    flake-parts.lib.mkFlake {inherit inputs;} {
      imports = [
        inputs.git-hooks-nix.flakeModule
      ];

      perSystem = {
        pkgs,
        config,
        ...
      }: {
        pre-commit.settings.hooks = {
        };

        devShells.default = pkgs.mkShell {
          inputsFrom = [config.pre-commit.devShell];

          packages = with pkgs; [
            go
            nodejs_24
            wails
          ];
        };
      };

      systems = ["x86_64-linux" "aarch64-linux"];
    };
}
