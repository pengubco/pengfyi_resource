# image-cache-release owns the workload volume: the second EBS volume attached
# to the instance as /dev/xvdb. It is formatted on first boot, mounted at
# /mnt/workload, seeded with symlinks into the image cache, and bind mounted
# over /var/lib/containerd before containerd starts. See README.md.
Name: %{_cross_os}image-cache-release
Version: 0.1.0
Release: 1%{?dist}
Summary: Serve the image cache to containerd from an empty workload volume
License: Apache-2.0 OR MIT
URL: https://github.com/pengubco/pengfyi_resource/tree/main/bottlerocket-image-cache

Source1: check-workload-volume.service
Source2: format-workload-volume.service
Source3: mnt-workload.mount
Source4: prepare-workload-volume.service
Source5: var-lib-containerd.mount

Requires: %{_cross_os}image-cache
Requires: %{_cross_os}image-cache-tools

%description
%{summary}.

%prep
%{nil}

%build
%{nil}

%install
install -d %{buildroot}%{_cross_unitdir}
install -p -m 0644 %{S:1} %{S:2} %{S:3} %{S:4} %{S:5} %{buildroot}%{_cross_unitdir}

%files
%{_cross_unitdir}/check-workload-volume.service
%{_cross_unitdir}/format-workload-volume.service
%{_cross_unitdir}/mnt-workload.mount
%{_cross_unitdir}/prepare-workload-volume.service
%{_cross_unitdir}/var-lib-containerd.mount
%{_cross_attribution_file}
