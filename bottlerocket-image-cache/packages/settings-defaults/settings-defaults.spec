%global _cross_first_party 1
%undefine _debugsource_packages

%global cargo_clean %{__cargo_cross_env} %{__cargo} clean

%global _cross_defaultsdir %{_cross_datadir}/storewolf

Name: %{_cross_os}settings-defaults
Version: 0.0
Release: 1%{?dist}
Summary: Settings defaults
License: Apache-2.0 OR MIT
URL: https://github.com/pengubco/pengfyi_resource/tree/main/bottlerocket-image-cache
BuildRequires: %{_cross_os}glibc-devel
Requires: %{_cross_os}settings-defaults(any)

%description
%{summary}.

%package aws-image-cache
Summary: Settings defaults for the aws-image-cache variant
Requires: %{_cross_os}variant(aws-image-cache)
Provides: %{_cross_os}settings-defaults(any)
Provides: %{_cross_os}settings-defaults(aws-image-cache)
Conflicts: %{_cross_os}settings-defaults(any)

%description aws-image-cache
%{summary}.

%prep
%setup -T -c
%cargo_prep

%build
declare -a projects
for defaults in \
  aws-image-cache \
  ;
do
  projects+=( "-p" "settings-defaults-$(echo "${defaults}" | sed -e 's,\.,_,g')" )
done

# Output is written to an unpredictable directory name, so clean it up first to
# avoid reusing any cached artifacts.
%cargo_clean --manifest-path %{_builddir}/sources/Cargo.toml \
  "${projects[@]}" \
  %{nil}

%cargo_build --manifest-path %{_builddir}/sources/Cargo.toml \
  "${projects[@]}" \
  %{nil}

%install
install -d %{buildroot}%{_cross_defaultsdir}
install -d %{buildroot}%{_cross_tmpfilesdir}

for defaults in \
  aws-image-cache \
  ;
do
  crate="$(echo "${defaults}" | sed -e 's,\.,_,g')"
  for f in $(find "${HOME}/.cache" -name "settings-defaults-${crate}.toml") ; do
    install -p -m 0644 "${f}" "%{buildroot}%{_cross_defaultsdir}/${defaults}.toml"
  done
  echo \
    "L+ /etc/storewolf/defaults.toml - - - - %{_cross_defaultsdir}/${defaults}.toml" > \
    "%{buildroot}%{_cross_tmpfilesdir}/storewolf-defaults-${defaults}.conf"
done

%files
%dir %{_cross_defaultsdir}

%files aws-image-cache
%{_cross_defaultsdir}/aws-image-cache.toml
%{_cross_tmpfilesdir}/storewolf-defaults-aws-image-cache.conf
