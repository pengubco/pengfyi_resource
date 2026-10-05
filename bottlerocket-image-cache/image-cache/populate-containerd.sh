#!/usr/bin/env bash
#
# Runs inside the Amazon Linux 2023 sandbox started by build-image-cache.sh.
# Populates a fresh containerd root with the internal and public images, then
# tars it up as the image cache.
set -euxo pipefail

root=/work/containerd
archives=/work/archives

dnf -y -q install containerd tar findutils

# Point containerd at the bind mounted work directory. The default config is
# fine otherwise: overlayfs snapshotter, no CRI needed.
mkdir -p "${root}" /run/containerd
containerd --root "${root}" --state /run/containerd \
  --address /run/containerd/containerd.sock \
  >/work/containerd.log 2>&1 &
containerd_pid=$!

for _ in $(seq 1 30); do
  ctr version >/dev/null 2>&1 && break
  sleep 1
done
ctr version

# Internal images: built in-house, never pushed anywhere. Importing the archive
# unpacks the layers into overlayfs snapshots as well as storing the blobs.
for archive in "${archives}"/*.tar; do
  ctr image import "${archive}"
done

# Public images that most workloads start from. Pulling them here means the
# host never has to.
for image in ${PUBLIC_IMAGES}; do
  ctr image pull "${image}"
done

ctr image ls

# containerd keeps the compressed layer blobs so it can push images later. The
# host only runs images, from the unpacked snapshots, so the compressed copies
# are dead weight: often as large as the snapshots themselves. Only the blobs
# carrying the "uncompressed" label are layers; manifests and configs stay.
if [ "${KEEP_COMPRESSED_LAYERS}" != "1" ]; then
  ctr content ls \
    | awk '/containerd\.io\/uncompressed=sha256/ { print $1 }' \
    | xargs -r ctr content rm
fi

# Stopping containerd too soon after the deletes loses some of them.
sleep 3
kill "${containerd_pid}"
wait "${containerd_pid}" || true

# Drop state that belongs to the sandbox, not to the images.
find "${root}" -type c -delete
rm -rf \
  "${root}/io.containerd.snapshotter.v1.native" \
  "${root}/io.containerd.snapshotter.v1.blockfile" \
  "${root}/io.containerd.snapshotter.v1.btrfs" \
  "${root}/io.containerd.runtime.v2.task" \
  "${root}/io.containerd.grpc.v1.introspection" \
  "${root}/io.containerd.sandbox.controller.v1.shim" \
  "${root}/io.containerd.snapshotter.v1.erofs" \
  "${root}/tmpmounts"

du -sh "${root}"/*

tar -C "${root}" -cf /work/image-cache.tar .

# Everything under /work was created as root. Hand the result back to the
# host user and drop the containerd root so the host can clean up build/.
rm -rf "${root}"
chown -R "${HOST_UID}:${HOST_GID}" /work
