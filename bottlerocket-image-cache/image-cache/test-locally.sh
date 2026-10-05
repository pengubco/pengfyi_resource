#!/usr/bin/env bash
#
# Exercise the image cache the way the Bottlerocket host will, without building
# an AMI: unpack the cache into one directory (the "root volume"), point
# image-cache-link at an empty second directory (the "workload volume"), start
# containerd on the second directory, and run the cached images.
#
#   test-locally.sh [IMAGE_CACHE_TAR]
set -euo pipefail

here="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
tar_path="${1:-${here}/../packages/image-cache/image-cache.tar}"
work="${here}/build/test"

if [ -d "${work}" ] && ! rm -rf "${work}" 2>/dev/null; then
  # A failed run leaves root-owned files behind.
  docker run --rm -v "${work}:/work" public.ecr.aws/amazonlinux/amazonlinux:2023 \
    sh -c 'rm -rf /work/* /work/.[!.]*'
  rm -rf "${work}"
fi
mkdir -p "${work}"

echo "==> building image-cache-link"
(cd "${here}/../sources/image-cache-tools" && \
  CGO_ENABLED=0 go build -o "${work}/image-cache-link" ./cmd/image-cache-link)

docker run --rm --privileged \
  -v "${work}:/work" \
  -v "${tar_path}:/image-cache.tar:ro" \
  -e "HOST_UID=$(id -u)" -e "HOST_GID=$(id -g)" \
  public.ecr.aws/amazonlinux/amazonlinux:2023 \
  bash -euxo pipefail -c '
dnf -y -q install containerd tar util-linux findutils

# cgroup v2 forbids processes in a cgroup that also has child cgroups. The
# container starts with everything in the root cgroup, so move it into a leaf
# before runc tries to create siblings for the demo containers.
mkdir -p /sys/fs/cgroup/init
while read -r pid; do echo "${pid}" > /sys/fs/cgroup/init/cgroup.procs 2>/dev/null || true; done < /sys/fs/cgroup/cgroup.procs
echo "+cpu +memory +pids" > /sys/fs/cgroup/cgroup.subtree_control

# The cache: read-only, as it will be on the Bottlerocket root volume.
mkdir -p /work/cache && tar -C /work/cache -xpf /image-cache.tar
mount --bind -o ro /work/cache /work/cache

# The workload volume: empty.
mkdir -p /work/containerd
/work/image-cache-link --root /work/containerd --cache /work/cache
find /work/containerd -maxdepth 4 | head -20

mkdir -p /run/containerd
containerd --root /work/containerd --state /run/containerd >/work/containerd.log 2>&1 &
for _ in $(seq 1 30); do ctr version >/dev/null 2>&1 && break; sleep 1; done

ctr image ls
ctr run --rm registry.example.internal/python-runtime:3.12 demo-runtime python /opt/runtime/bootstrap.py serve --port 8080
ctr run --rm registry.example.internal/log-router:1.0 demo-router
ctr run --rm public.ecr.aws/docker/library/alpine:3.20 demo-alpine cat /etc/alpine-release

# Nothing was pulled: the cache alone served all three images.
ctr content ls | wc -l

umount /work/cache
chown -R "${HOST_UID}:${HOST_GID}" /work
'
echo "==> all three cached images ran from an empty containerd root"
