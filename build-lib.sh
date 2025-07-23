#!/usr/bin/env bash
set -euo pipefail

# assume this script lives in your project root
ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
RUST_LIB_DIR="${ROOT_DIR}/rust-lib"

cd "${RUST_LIB_DIR}"
cargo build --release

# pick the right file extension by OS
case "$(uname -s)" in
  Linux)
    LIB_NAME="liborderffi.so"
    ;;
  Darwin)
    LIB_NAME="liborderffi.dylib"
    ;;
  *)
    echo "Unsupported OS: $(uname -s)" >&2
    exit 1
    ;;
esac

# copy it back to the project root
cp "target/release/${LIB_NAME}" "${ROOT_DIR}/"
echo "${LIB_NAME} has been copied to ${ROOT_DIR}/"
