#!/usr/bin/env bash
#
# Build the image cache tarball consumed by packages/image-cache.
#
#   build-image-cache.sh [OUTPUT_TAR]
#
# 1. Builds the "internal" images under images/ with Docker and saves each one
#    as an OCI archive. These stand in for images that are built in-house and
#    never pushed to a registry.
# 2. Starts a throwaway Amazon Linux 2023 container running the same containerd
#    major version as the Bottlerocket variant, imports the archives, pulls a
#    few public images, strips what the host does not need, and tars the
#    resulting containerd root directory. See populate-containerd.sh.
#
# The containerd root must live on a real filesystem, not on the container's
# overlayfs, because the overlayfs snapshotter needs to use it as an upperdir.
# That is why build/ is bind mounted from the host.
#
# Environment:
#   PUBLIC_IMAGES            space separated list of extra images to pull (default: alpine:3.20)
#   KEEP_COMPRESSED_LAYERS   set to 1 to keep compressed layer blobs in the cache
set -euo pipefail

here="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
out="${1:-${here}/../packages/image-cache/image-cache.tar}"
build="${here}/build"

# Build on the same architecture as the AMI. containerd snapshots are
# architecture specific and there is no point caching the wrong ones.
case "$(uname -m)" in
  x86_64)  platform="linux/amd64" ;;
  aarch64) platform="linux/arm64" ;;
  *) echo "unsupported architecture $(uname -m)" >&2; exit 1 ;;
esac

# name:tag as the host will refer to it, and the directory it is built from.
internal_images=(
  "registry.example.internal/python-runtime:3.12 python-runtime"
  "registry.example.internal/log-router:1.0 log-router"
)
public_images="${PUBLIC_IMAGES:-public.ecr.aws/docker/library/alpine:3.20}"

if [ -d "${build}" ] && ! rm -rf "${build}" 2>/dev/null; then
  # A failed sandbox run leaves root-owned files behind.
  docker run --rm -v "${build}:/work" public.ecr.aws/amazonlinux/amazonlinux:2023 \
    sh -c 'rm -rf /work/* /work/.[!.]*'
  rm -rf "${build}"
fi
mkdir -p "${build}/archives"

for entry in "${internal_images[@]}"; do
  read -r ref dir <<<"${entry}"
  echo "==> docker build ${ref}"
  docker build --platform "${platform}" -t "${ref}" "${here}/images/${dir}"
  # ctr can import a docker archive directly; "type=oci" is not required.
  docker save "${ref}" -o "${build}/archives/${dir}.tar"
done

echo "==> populating containerd in a sandbox"
docker run --rm --privileged \
  --platform "${platform}" \
  -v "${build}:/work" \
  -v "${here}/populate-containerd.sh:/populate-containerd.sh:ro" \
  -e "PUBLIC_IMAGES=${public_images}" \
  -e "KEEP_COMPRESSED_LAYERS=${KEEP_COMPRESSED_LAYERS:-0}" \
  -e "HOST_UID=$(id -u)" -e "HOST_GID=$(id -g)" \
  public.ecr.aws/amazonlinux/amazonlinux:2023 \
  bash /populate-containerd.sh

mkdir -p "$(dirname "${out}")"
mv "${build}/image-cache.tar" "${out}"
echo "==> wrote ${out} ($(du -h "${out}" | cut -f1))"
