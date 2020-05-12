Name:           baasapi
Version:        1.20.2
Release:        0
License:        Zlib
Summary:        A lightweight docker management UI
Url:            https://baasapi.io
Group:          BLAH
Source0:        https://github.com/baasapi/baasapi/releases/download/%{version}/baasapi-%{version}-linux-amd64.tar.gz
Source1:        baasapi.service
BuildRoot:      %{_tmppath}/%{name}-%{version}-build
%if 0%{?suse_version}
BuildRequires:  help2man
%endif
Requires:       docker
%{?systemd_requires}
BuildRequires: systemd

## HowTo build ## 
# You can use spectool to fetch sources
# spectool -g -R distribution/baasapi.spec 
# Then build with 'rpmbuild -ba distribution/baasapi.spec' 


%description
BaaSapi is a lightweight management UI which allows you to easily manage
your different Docker environments (Docker hosts or Swarm clusters).
BaaSapi is meant to be as simple to deploy as it is to use.
It consists of a single container that can run on any Docker engine
(can be deployed as Linux container or a Windows native container).
BaaSapi allows you to manage your Docker containers, images, volumes,
networks and more ! It is compatible with the standalone Docker engine and with Docker Swarm mode.

%prep
%setup -qn baasapi

%build
%if 0%{?suse_version}
help2man -N --no-discard-stderr ./baasapi  > baasapi.1
%endif

%install
# Create directory structure
install -D -m 0755 baasapi %{buildroot}%{_sbindir}/baasapi
install -d -m 0755 %{buildroot}%{_datadir}/baasapi/public
install -d -m 0755 %{buildroot}%{_localstatedir}/lib/baasapi
install -D -m 0644 %{S:1} %{buildroot}%{_unitdir}/baasapi.service
%if 0%{?suse_version}
install -D -m 0644 baasapi.1 %{buildroot}%{_mandir}/man1/baasapi.1
( cd  %{buildroot}%{_sbindir} ; ln -s service rcbaasapi )
%endif
# populate
# don't install docker binary with package use system wide installed one
cp -ra public/ %{buildroot}%{_datadir}/baasapi/

%pre
%if 0%{?suse_version}
%service_add_pre baasapi.service
#%%else # this does not work on rhel 7?
#%%systemd_pre baasapi.service
true
%endif

%preun
%if 0%{?suse_version}
%service_del_preun baasapi.service
%else
%systemd_preun baasapi.service
%endif

%post
%if 0%{?suse_version}
%service_add_post baasapi.service
%else
%systemd_post baasapi.service
%endif

%postun
%if 0%{?suse_version}
%service_del_postun baasapi.service
%else
%systemd_postun_with_restart baasapi.service
%endif


%files
%defattr(-,root,root)
%{_sbindir}/baasapi
%{_datadir}/baasapi/public
%dir %{_localstatedir}/lib/baasapi/
%{_unitdir}/baasapi.service
%if 0%{?suse_version}
%{_mandir}/man1/baasapi.1*
%{_sbindir}/rcbaasapi
%endif
