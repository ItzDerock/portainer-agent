{ pkgs, lib, config, inputs, ... }:

{
  packages = with pkgs; [ git wget ];
  languages.go.enable = true;
}
