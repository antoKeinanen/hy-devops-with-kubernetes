{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
  };

  outputs =
    {
      self,
      nixpkgs,
    }:
    let
      pkgs = nixpkgs.legacyPackages."x86_64-linux";

      release-exercise = pkgs.writeShellApplication {
        name = "release-exercise";

        runtimeInputs = with pkgs; [
          git
          gh
        ];

        text = ''
          set -euo pipefail

          if [ "$#" -ne 1 ]; then
            echo "Usage: release-exercise <tag>"
            echo "Example: release-exercise 1.3"
            exit 1
          fi

          tag="$1"
          readme="README.md"

          if [ ! -f "$readme" ]; then
            if [ -f "README" ]; then
              readme="README"
            else
              echo "Could not find README.md or README"
              exit 1
            fi
          fi

          line="- [$tag](https://github.com/antoKeinanen/hy-devops-with-kubernetes/tree/$tag)"

          # Ensure the README ends with a newline before appending.
          if [ -s "$readme" ] && [ "$(tail -c 1 "$readme")" != "" ]; then
            printf '\n' >> "$readme"
          fi

          printf '%s\n' "$line" >> "$readme"

          git add "$readme"
          git commit -m "add $tag"

          current_branch="$(git branch --show-current)"

          if [ -z "$current_branch" ]; then
            echo "You appear to be in detached HEAD state. Cannot push safely."
            exit 1
          fi

          git push -u origin "$current_branch"

          commit_sha="$(git rev-parse HEAD)"

          gh release create "$tag" \
            --target "$commit_sha" \
            --title "$tag" \
            --notes ""
        '';
      };
    in
    {
      devShells."x86_64-linux".default = pkgs.mkShell {
        buildInputs = with pkgs; [
          k3d
          kubectl
          python314
          go
          kubectx
          gh
          release-exercise
        ];

        shellHook = ''
          source <(kubectl completion zsh)
        '';
      };
    };
}
