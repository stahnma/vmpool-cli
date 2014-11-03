

# Installation

## Homebrew

for usage: clone this project

    brew install --HEAD formula/vmpool.rb


## RPM
For RPM based systems


for usage: clone this project


    make srpm

# put the srpm in mock against the target you want to build for

# Grabbing multiple VMs

To fetch multiple vms simply provide whatever valid platforms
you desire as arguements to the grab command.

For example:
```
~> vmpool grab debian-7-x86_64 debian-7-x86_64 centos-7-i386
centos-7-i386: qv7uij5ofqj7pqs.delivery.puppetlabs.net
debian-7-x86_64:
    qer906rgo66zhwp.delivery.puppetlabs.net
    f2uie2kr56hjg34.delivery.puppetlabs.net
```
