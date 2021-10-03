{ pkgs ? import <nixpkgs> {} }:

pkgs.buildGoModule {
  pname = "jtk";
  version = "0.1";

  vendorSha256 = null;

  src = ./.;
}