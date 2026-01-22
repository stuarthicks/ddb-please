{
  description = "ddb-please - DynamoDB helper";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";

  outputs = { self, nixpkgs }:
    let
      systems = [ "x86_64-linux" "aarch64-linux" "x86_64-darwin" "aarch64-darwin" ];
      forAll = f: nixpkgs.lib.genAttrs systems (s: f nixpkgs.legacyPackages.${s});
    in {
      packages = forAll (pkgs: rec {
        ddb-please = pkgs.buildGoModule {
          pname = "ddb-please";
          version = self.shortRev or self.dirtyShortRev or "dev";
          src = ./.;
          vendorHash = "sha256-TGbf6DnV4/c2gfoQ7sHLL2IhqxH0REckZScwvfebN6E=";
          ldflags = [ "-s" "-w" ];
          meta.mainProgram = "ddb-please";
        };
        default = ddb-please;
      });

      devShells = forAll (pkgs: {
        default = pkgs.mkShell { packages = [ pkgs.go pkgs.just pkgs.goreleaser ]; };
      });
    };
}
