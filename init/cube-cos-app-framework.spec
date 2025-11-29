Name:           appctl
Version:        %{version}
Release:        1%{?dist}.%{build_number}
Summary:        App Framework Binary for CubeCOS

License:        Apache License 2.0
URL:            https://github.com/bigstack-oss/cube-cos-app-framework
Source0:        https://github.com/bigstack-oss/cube-cos-app-framework/tree/%{build_number}

BuildRequires: systemd golang

%description
The App Framework Binary for CubeCOS.

%prep
rm -rf ./*
cp %{_topdir}/SOURCES/"cube-cos-app-framework-%{version}.tar.gz" .
tar -xzf "cube-cos-app-framework-%{version}.tar.gz"
rm "cube-cos-app-framework-%{version}.tar.gz"
find ./source/ -mindepth 1 -maxdepth 1 -name  '*' -exec mv -t . {} +
rmdir source

%build
GOWORK=off go mod tidy -v
go clean -v
GOWORK=off GOOS=linux GOARCH=amd64 go build -o %{name} -v main.go

%install
rm -rf $RPM_BUILD_ROOT
mkdir -p $RPM_BUILD_ROOT/usr/local/bin
mv %{name} $RPM_BUILD_ROOT/usr/local/bin
mkdir -p $RPM_BUILD_ROOT/%{_sysconfdir}/cube/app-framework
cp ./configs/cube-cos-app-framework.yaml.template $RPM_BUILD_ROOT/%{_sysconfdir}/cube/app-framework/cube-cos-app-framework.yaml
mkdir -p $RPM_BUILD_ROOT/%{_unitdir}
mkdir -p $RPM_BUILD_ROOT/%{_datadir}/cube/app-framework
cp LICENSE $RPM_BUILD_ROOT/%{_datadir}/cube/app-framework

%files
/usr/local/bin/%{name}
%{_sysconfdir}/cube/app-framework/cube-cos-app-framework.yaml
%{_datadir}/cube/app-framework/LICENSE
