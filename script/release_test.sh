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

  # Sections other than builds: reference build ids through their "ids:" list
  # (archives, nfpms, ...), so none of those references may point at a build
  # that the platform filter removed.
  referenced_ids="$(awk '
    /^builds:/ { in_builds = 1; next }
    in_builds && /^[^[:space:]]/ { in_builds = 0 }
    in_builds { next }
    match($0, /ids:[[:space:]]*\[[^][]*\]/) {
      list = substr($0, RSTART, RLENGTH)
      sub(/^ids:[[:space:]]*\[/, "", list)
      sub(/\]$/, "", list)
      count = split(list, item, ",")
      for (i = 1; i <= count; i++) {
        gsub(/[[:space:]]/, "", item[i])
        if (item[i] != "") print item[i]
      }
    }
  ' "$capture" | sort -u)"

  for referenced_id in $referenced_ids; do
    if ! printf '%s\n' "$build_ids" | grep -qx "$referenced_id"; then
      printf 'platform %s: ids reference "%s" is not a build of the filtered config\n' \
        "$platform" "$referenced_id" >&2
      exit 1
    fi
  done

  # The selected build is the only one left, so it is the only one that may
  # still be referenced.
  if [ "$referenced_ids" != "$expected_build" ]; then
    printf 'platform %s: expected ids references [%s], got [%s]\n' \
      "$platform" "$expected_build" "$referenced_ids" >&2
    exit 1
  fi

  # A section whose entries were all filtered out must be gone: GoReleaser
  # rejects a config with an empty section.
  if ! awk '
    function is_section_key(s) {
      return s ~ /^[A-Za-z_][A-Za-z0-9_.-]*:[[:space:]]*(#.*)?$/
    }
    is_section_key($0) {
      if (pending) {
        print "section without entries: " pending_line >"/dev/stderr"
        empty = 1
      }
      pending = 1
      pending_line = $0
      next
    }
    $0 !~ /^[[:space:]]*(#|$)/ { pending = 0 }
    END {
      if (pending) {
        print "section without entries: " pending_line >"/dev/stderr"
        empty = 1
      }
      exit empty ? 1 : 0
    }
  ' "$capture"; then
    printf 'platform %s: filtered config has a section without entries\n' "$platform" >&2
    exit 1
  fi
done

echo "platform filtering passed"
