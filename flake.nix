{
  inputs = {
    nixpkgs.url = github:NixOS/nixpkgs/nixos-23.11;
    nixpkgs_darwin.url = github:NixOS/nixpkgs/nixpkgs-23.11-darwin;
  };

  outputs = { self, nixpkgs, nixpkgs_darwin }: let

    version = self.rev or "dirty";

    packages = nixpkgs: sys: emacs-flavor: let
      pkgs = import nixpkgs { system = sys; };
      epkgs = pkgs.emacsPackagesFor (emacs-flavor pkgs);

      pus = pkgs.buildGoModule {
          name = "pus-${version}";
          version = "${version}";
          src = ./.;
          doCheck = false;
          vendorHash = "sha256-9JkGaiSp+WAzVH6zyzSYi/HEeWMTAtEgLItAx1rre+Y=";
      };

      bootstrap = pkgs.stdenv.mkDerivation {
        name = "pus-bootstrap-${version}";
        version = "${version}";
        builder = "${pkgs.bash}/bin/bash";
        coreutils = pkgs.coreutils;
        args = [ ./bootstrap.sh ./pus.el "${pus}/bin/pus" ];
      };

      pus-el = epkgs.elpaBuild {
        pname = "pus";
        ename = "pus";
        version = version;
        src = [ "${bootstrap}/pus.el" ];
        packageRequires = [];
        meta = {};
      };

    in {
      default = pus;
      bootstrap = bootstrap;
      pus = pus;
      pus-el = pus-el;
    };

  in {
    packages = {
      "x86_64-darwin" = packages nixpkgs_darwin "x86_64-darwin" (pkgs: pkgs.emacs29-macport);
      "aarch64-darwin" = packages nixpkgs_darwin "aarch64-darwin" (pkgs: pkgs.emacs29-macport);
      "x86_64-linux" = packages nixpkgs "x86_64-linux" (pkgs: pkgs.emacs29);
      "aarch64-linux" = packages nixpkgs "aarch64-linux" (pkgs: pkgs.emacs29);
    };
  };
}
