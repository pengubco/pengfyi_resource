%global _cross_first_party 1
%undefine _debugsource_packages

Name: %{_cross_os}migrations
Version: 0.0
Release: 1%{?dist}
Summary: Settings migrations
License: Apache-2.0 OR MIT
URL: https://github.com/pengubco/pengfyi_resource/tree/main/bottlerocket-image-cache

# The build system expects to find a package named "bottlerocket-migrations".
Provides: %{_cross_os}settings-migrations

%description
%{summary}.

%prep
%{nil}

%build
%{nil}

%install
install -d %{buildroot}%{_cross_datadir}/migrations

%files
%dir %{_cross_datadir}/migrations
