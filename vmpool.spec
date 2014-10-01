%define debug_package %{nil}

Name:      vmpool
Version:   ==VERSION==
Release:   1%{?dist}
Summary:   CLI Utility for vmpooler

Group:     Development/Tools
License:   ASL 2.0
URL:       https://github.com/stahnma/vmpool-cli

# Tarball given to you via make
Source0:   %{name}-%{version}.tar.gz

BuildRequires: golang, make

%description
VMPooler CLI for usage with https://github.com/puppetlabs/vmpooler

Source for this package at: https://github.com/stahnma/vmpool-cli


%prep
%setup -q


%build
VERSION=%{version} make %{?_smp_mflags}


%install
make install DESTDIR=%{buildroot}


%files
%doc
/usr/local/bin/vmpool



%changelog
* Tue Sep 30 2014 Michael Stahnke <stahnma@puppetlabs.com>
- Initial package
