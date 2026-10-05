# image-cache installs a containerd root directory, captured after containerd
# had pulled the images we want every host to start with, onto the read-only
# root volume at /usr/share/image-cache. image-cache-release wires it into the
# live containerd root on the workload volume at boot.
#
# Everything on the Bottlerocket root partition is labeled os_t, which the
# container runtime is not allowed to manipulate. containerd also expects to
# own the files under its root. Rather than loosen SELinux policy, the cache is
# exposed through a read-only overlayfs mount at /var/lib/image-cache whose
# context= option relabels it as local_t, the label of containerd's own data
# under /var/lib. A lowerdir-only overlay needs no upperdir or workdir.
#
# This is a .mount unit rather than a mount command in a .service so that
# other units can order on it with RequiresMountsFor=.

# The cache contains binaries from other distributions' images, so none of the
# post-install checks (strip, debuginfo, FIPS) apply; keep the files as they are.
%global __os_install_post /usr/bin/true
%undefine cross_check_fips
%global __spec_clean_pre exit 0

Name: %{_cross_os}image-cache
Version: 0.1.0
Release: 1%{?dist}
Summary: Container images cached on the root volume
License: Apache-2.0 OR MIT
URL: https://github.com/pengubco/pengfyi_resource/tree/main/bottlerocket-image-cache

Source0: image-cache.tar
Source1: image-cache.mount

%description
%{summary}.

%prep
%{nil}

%build
%{nil}

%install
install -d %{buildroot}%{_cross_datadir}/image-cache
# SELinux labels and xattrs from the build sandbox do not apply on the host.
tar --no-selinux --no-xattrs -xpf %{S:0} -C %{buildroot}%{_cross_datadir}/image-cache

# rpmbuild --clean has to be able to delete the build root, and the file
# checks have to be able to read it; containerd snapshots contain read-only
# and oddly permissioned entries that would otherwise block both.
chmod +rwX -R %{buildroot}%{_cross_datadir}/image-cache

# systemd requires a .mount unit to be named after its mount point, with the
# path escaped (/var/lib/image-cache -> var-lib-image\x2dcache.mount).
install -d %{buildroot}%{_cross_unitdir}
install -p -m 0644 %{S:1} \
  %{buildroot}%{_cross_unitdir}/$(systemd-escape --path /var/lib/image-cache).mount

%files
%dir %{_cross_datadir}/image-cache
%{_cross_datadir}/image-cache/*
%{_cross_unitdir}/var-lib-image\x2dcache.mount
%{_cross_attribution_file}
