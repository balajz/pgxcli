#!/bin/bash
set -euo pipefail

repo_root="$(cd "$(dirname "$0")/.." && pwd)"
tmpdir="$(mktemp -d)"
trap 'rm -rf "$tmpdir"; rm -f "$repo_root/.goreleaser.generated.yaml"' EXIT

cat >"$tmpdir/goreleaser" <<'EOF'
#!/bin/bash
set -euo pipefail

config=""
while [ "$#" -gt 0 ]; do
  case "$1" in
    -f)
      config="$2"
      shift 2
      ;;
    *)
      shift
      ;;
  esac
done

cp "$config" "${CAPTURE_FILE:?}"
EOF
chmod +x "$tmpdir/goreleaser"

for platform_and_build in "linux linux" "macos darwin" "windows windows"; do
  platform="${platform_and_build% *}"
  expected_build="${platform_and_build#* }"
  capture="$tmpdir/${platform}.yaml"

  (
    cd "$repo_root"
    PATH="$tmpdir:$PATH" CAPTURE_FILE="$capture" ./script/release --local --platform "$platform"
  )

  build_ids="$(awk '
    /^builds:/ { in_builds = 1; next }
    in_builds && /^[^[:space:]]/ { in_builds = 0 }
    in_builds && /^  - id: / { print $3 }
  ' "$capture")"

  test "$build_ids" = "$expected_build"
done

echo "platform filtering passed"
