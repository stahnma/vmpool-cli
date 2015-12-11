#NAME

**vmpool** -- interact with a vmpooler


#SYNOPSIS

vmpool <grab|delete|list|status|summary|token|version|vm>

#DESCRIPTION

vmpool is a fast multi-platform API client for a vmpooler fronted infrastructure. vmpool implements much of the v1 specification for the vmpooler.

#SUBCOMMANDS

##grab
Grab uses an HTTP POST to request one or more VMs from the pool(s) specified.

_Note:_ As of verison 0.2.0 all VMs are tagged with User: LDAP_USERNAME as to be tracable.

Usage:

    vmpool grab <pool name> <pool name>

Exit Codes:

0 - Instance Reserved

1 - General Error

##delete
Destroy an instance from the vmpooler.

Usage

    vmpool delete <hostname>
    vmpool vm delete <hostname>

Exit Codes:

0 - Host destroyed

1 - General Error

##list
List available pool names on stdout.

Usage

    vmpool list [filter]

Exit Codes:

0 - Pools Listed

1 - General Error


##status
Display vmpooler health information via the status endpoint.

    vmpool status

Exit Codes:

0 - Summary Listed

1 - General Error

##summary
Display summary information for the vmpooler.

_Warning_: This can be long and verbose.

Usage

    vmpool summary
    
Exit Codes:

0 - Summary Listed

1 - General Error

##token
Interact with tokens. Tokens are the authentication key used by the vmpooler. In most cases, the vmpool client will aquire/retrieve a token from the vmpool API without direct user interaction.

All token interactions require LDAP_USERNAME and LDAP_PASSWORD environment variables to be set.

Usage:
     
Delete a the specified token.

    vmpool token delete <token>
    
List all tokens assigned to LDAP_USERNAME.
    
    vmpool token list
    
Remove or invalidate any/all tokens assigned to LDAP_USERNAME.
    
    vmpool token purge
    
Request a new token for LDAP_USERNAME.    
    
    vmpool token request

##version

Display version of vmpool.

##vm
Interact with specific vms assigned from the pooler.

Usage:

Destroy an instance from the vmpooler.
    
     vmpool vm delete <hostname>

Request a VM from the pool specified. 

     vmpool vm grab <pool name>

Display detailed information for a vm.

     vmpool vm info <hostname>

Alter the lifetime of a VM.


    vmpool vm lifetime <hostname> <TLL in hours>


#ENVIRONMENT
**DEBUG** - if set to 1, lots of developer DEBUG information is displayed on stdout. By default debugging mode is not enabled.

**LDAP_USERNAME** - Username for the LDAP directory that the vmpooler is configured to use. There is no default value for this setting. vmpool may error out with it being set.

**LDAP_PASSWORD** - Password for `LDAP_USERNAME`.

**VMPOOL_LOGFILE** - Location of the application log for vmpool. By default this is in `$HOME/.vmpool.log`.

**VMPOOLER_TOKEN** - If you already have a TOKEN for the vmpooler, you may specify it here. This will reduce HTTP calls to the pooler and may cause a very slight performance boost. 

**VMPOOL_URL** - This is the URL of the vmpooler you are connecting to. This defaults to https://vmpooler.delivery.puppetlabs.net


#EXAMPLES

Get a VM and adjust the lifetime to 2 hours.

    $  vmpool grab fedora-22-i386
    fedora-22-i386: t7g9t6kim9phvli.delivery.puppetlabs.net
    {
       "ok": true,
       "t7g9t6kim9phvli": {
       "template": "fedora-22-i386",
       "lifetime": 12,
       "running": 0.01,
       "state": "running",
       "tags": {
         "user": "stahnma",
         "client": "vmpool-cli-0.1.0.65.gbe0d847"
        },
        "domain": "delivery.puppetlabs.net"
      }
    }
    $ vmpool vm lifetime t7g9t6kim9phvli 2
    Lifetime changed to 2 on t7g9t6kim9phvli

#SEE ALSO

https://github.com/puppetlabs/vmpooler/blob/master/API.md - Official API documetnation


#AUTHORS

Written primarily by Michael Stahnke with significant patches from Andrew Roetker.

