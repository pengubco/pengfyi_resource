# bottlerocket-image-cache

A [Bottlerocket](https://bottlerocket.dev/) variant that ships a set of container images on the
read-only root volume and serves them to containerd from an empty, separately encryptable
workload volume. Boot an instance from the AMI and the images are already there; containerd never
pulls them, and the volume that holds everything containerd writes starts out blank.

The companion blog post explains why and walks through the mechanics:
[Shipping container images on the Bottlerocket root volume](https://peng.fyi/post/pre-load-container-images-on-the-bottlerocket-read-only-root-fs-volume/).

## How it works

```
build time   image-cache/build-image-cache.sh
               docker build the internal images, docker save them
               run containerd 2.2 in an AL2023 sandbox, ctr import / ctr pull
               prune, tar /var/lib/containerd        -> packages/image-cache/image-cache.tar

AMI build    packages/image-cache          untar to /usr/share/image-cache (root volume, os_t)
                                           var-lib-image\x2dcache.mount: overlay, context=local_t
             packages/image-cache-tools    /usr/bin/image-cache-link (Go)
             packages/image-cache-release  the boot units below
             variants/aws-image-cache      the variant; amispec.toml adds a blank /dev/xvdb

first boot   check-workload-volume.service    blkid /dev/xvdb, skip the next two on reboot
             format-workload-volume.service   mkfs.ext4 -L WORKLOAD
             mnt-workload.mount               LABEL=WORKLOAD -> /mnt/workload
             var-lib-image\x2dcache.mount     /usr/share/image-cache -> /var/lib/image-cache
             prepare-workload-volume.service  image-cache-link --root /mnt/workload/containerd --cache /var/lib/image-cache
             var-lib-containerd.mount         /mnt/workload/containerd --rbind-> /var/lib/containerd
             containerd.service               starts with the cached images already in its root
```

image-cache-link copies containerd's two bolt databases into the empty root and symlinks every
content blob and every overlayfs snapshot back into the cache. containerd follows the symlinks.

## Layout

| Path | What |
|---|---|
| `image-cache/images/*` | Dockerfiles for the two stand-in "internal" images: a Python runtime and a log router sidecar |
| `image-cache/build-image-cache.sh` | Builds the tarball. Needs Docker. |
| `image-cache/populate-containerd.sh` | Runs inside the sandbox: ctr import, ctr pull, prune, tar |
| `image-cache/test-locally.sh` | Runs the whole mechanism in a container, no AMI needed |
| `sources/image-cache-tools` | Go module for image-cache-link, with unit tests |
| `packages/image-cache` | RPM: the tarball on the root volume plus the relabeling overlay mount |
| `packages/image-cache-tools` | RPM: image-cache-link |
| `packages/image-cache-release` | RPM: the first-boot units for the workload volume |
| `packages/settings-*`, `sources/settings-*` | The settings model every Bottlerocket variant must provide; copied from upstream aws-dev |
| `variants/aws-image-cache` | The variant manifest and AMI block device spec |

Kit and SDK versions in [Twoliter.toml](Twoliter.toml) track upstream Bottlerocket v1.66.0.

## Getting the source

This project lives in a subfolder of the pengfyi_resource repository. Every command below runs
from that subfolder:

```sh
git clone https://github.com/pengubco/pengfyi_resource.git
cd pengfyi_resource/bottlerocket-image-cache
```

## Prerequisites

- Linux x86_64 or aarch64 host with Docker. The build runs containerd inside a privileged
  container, so Docker Desktop style setups that cannot run privileged containers will not work.
- Rust toolchain with `cargo make` (`cargo install cargo-make`).
- Go 1.24 or later, only for `image-cache/test-locally.sh` and `go test`.
- AWS CLI with credentials that can register AMIs and launch EC2 instances, for the last steps.

About 40 GiB of free disk and 30 minutes for the first build. Subsequent builds reuse the
downloaded SDK and kits and finish in a minute or two.

## Step 1: build the image cache

```sh
cargo make build-image-cache
```

This writes `packages/image-cache/image-cache.tar` (about 77 MiB) containing three images:

```
registry.example.internal/python-runtime:3.12   built from image-cache/images/python-runtime
registry.example.internal/log-router:1.0        built from image-cache/images/log-router
public.ecr.aws/docker/library/alpine:3.20       pulled
```

Add your own by editing `internal_images` in build-image-cache.sh or setting `PUBLIC_IMAGES`.

To see the whole mechanism work without building an AMI:

```sh
image-cache/test-locally.sh
```

It unpacks the tarball into a read-only bind mount, runs image-cache-link against an empty
directory, starts containerd there, and runs all three images.

## Step 2: build the image

```sh
cargo make build-variant
```

Output lands in `build/images/x86_64-aws-image-cache/latest/`. The first run installs Twoliter
into `tools/twoliter`, pulls the SDK and kits from public ECR, and generates local Secure Boot
keys under `sbkeys/local`.

To build only one package while iterating:

```sh
cargo make -e PACKAGE=image-cache-release build-package
```

## Step 3: register the AMI

```sh
cp Infra.toml.example Infra.toml   # or: export PUBLISH_REGIONS=us-west-2
cargo make ami
```

The AMI gets two block device mappings: `/dev/xvda` (the root volume, from the built image, with
the read-only root partition and Bottlerocket's writable data partition on it) and `/dev/xvdb`
(the workload volume, 20 GiB blank gp3, from
[variants/aws-image-cache/amispec.toml](variants/aws-image-cache/amispec.toml)). The AMI ID is
printed and written to `build/images/x86_64-aws-image-cache/latest/bottlerocket-aws-image-cache-x86_64-amis.json`.

## Step 4: launch and look around

Launch with an instance profile that allows SSM Session Manager
(`AmazonSSMManagedInstanceCore`), and enable the admin container through user data so you can
get a root shell:

```sh
cat > userdata.toml <<'EOF'
[settings.host-containers.admin]
enabled = true
EOF

aws ec2 run-instances \
  --image-id "${AMI_ID}" \
  --instance-type t3.small \
  --iam-instance-profile Name="${INSTANCE_PROFILE}" \
  --subnet-id "${SUBNET_ID}" \
  --user-data file://userdata.toml \
  --tag-specifications 'ResourceType=instance,Tags=[{Key=Name,Value=image-cache-demo}]'
```

To encrypt the workload volume with your own KMS key, add a `--block-device-mappings` override
for `/dev/xvdb` with `Encrypted=true,KmsKeyId=...`. Nothing on the host changes; the volume still
arrives blank and gets formatted on first boot.

Connect through SSM to the control container, then enter the admin container and drop to the host:

```sh
aws ssm start-session --target "${INSTANCE_ID}"
```

```
$ enter-admin-container
$ sudo sheltie
```

Now on the host:

```
# findmnt /var/lib/image-cache /mnt/workload /var/lib/containerd
TARGET              SOURCE                                   FSTYPE  OPTIONS
/var/lib/image-cache overlay                                 overlay ro,...,lowerdir=/usr/share/image-cache:/var/lib/image-cache,context=...local_t...
/mnt/workload       /dev/nvme1n1                             ext4    rw,...
/var/lib/containerd /dev/nvme1n1[/containerd]                ext4    rw,...

# ls -l /var/lib/containerd/io.containerd.snapshotter.v1.overlayfs/snapshots | head -3
lrwxrwxrwx. 1 root root 59 ... 1 -> /var/lib/image-cache/io.containerd.snapshotter.v1.overlayfs/snapshots/1
lrwxrwxrwx. 1 root root 59 ... 2 -> /var/lib/image-cache/io.containerd.snapshotter.v1.overlayfs/snapshots/2

# ctr image ls -q
public.ecr.aws/docker/library/alpine:3.20
registry.example.internal/log-router:1.0
registry.example.internal/python-runtime:3.12

# ctr run --rm registry.example.internal/python-runtime:3.12 demo python /opt/runtime/bootstrap.py serve
python-runtime: python 3.12.15 on x86_64
python-runtime: requests 2.32.3 preinstalled
python-runtime: would exec workload: serve

# ctr run --rm registry.example.internal/log-router:1.0 router
log-router: forwarding /var/log/workload
log-router: /var/log/workload does not exist, nothing to forward
```

The instance has no route to registry.example.internal; the images ran from the root volume. Reboot
the instance and `check-workload-volume.service` finds the filesystem, so the volume is neither
reformatted nor re-seeded and containerd keeps whatever it wrote in the meantime.

## Cleaning up

```sh
cargo make clean               # build/ and target/
cargo make clean-image-cache   # the tarball and image-cache/build/
```

## License

Apache-2.0 OR MIT, like Bottlerocket. The build scripts, specs and units under `packages/`,
`sources/settings-*`, `tools/` and `sbkeys/` are derived from the
[bottlerocket](https://github.com/bottlerocket-os/bottlerocket) repository.
