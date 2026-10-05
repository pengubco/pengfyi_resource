%global _cross_first_party 1
%undefine _debugsource_packages

Name: %{_cross_os}image-cache-tools
Version: 0.0
Release: 1%{?dist}
Summary: Tools to serve the image cache to containerd
License: Apache-2.0 OR MIT
URL: https://github.com/pengubco/pengfyi_resource/tree/main/bottlerocket-image-cache

BuildRequires: %{_cross_os}glibc-devel

%description
%{summary}.

%prep
%setup -T -c
cp -r %{_builddir}/sources/image-cache-tools/* .

%build
%set_cross_go_flags
go build -ldflags="${GOLDFLAGS}" -o image-cache-link ./cmd/image-cache-link

%install
install -d %{buildroot}%{_cross_bindir}
install -p -m 0755 image-cache-link %{buildroot}%{_cross_bindir}

%files
%{_cross_bindir}/image-cache-link
