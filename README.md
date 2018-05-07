# factd
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/85ca8315d2884dbca0b716a800310103)](https://www.codacy.com?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=twhiston/factd&amp;utm_campaign=Badge_Grade)
[![Codacy Badge](https://api.codacy.com/project/badge/Coverage/85ca8315d2884dbca0b716a800310103)](https://www.codacy.com?utm_source=github.com&utm_medium=referral&utm_content=twhiston/factd&utm_campaign=Badge_Coverage)
[![pipeline status](https://gitlab.com/twhiston/factd/badges/master/pipeline.svg)](https://gitlab.com/twhiston/factd/commits/master)

Factd is a fact reporting daemon, similar to puppetlabs facte.
It is intended to be run as a process on the target system, via systemd.
Note that in a container it will need some extra privilages and mounts to be able to access the host information

Facts are served to an http endpoint, as well as prometheus metrics and pprof data.

factd is under heavy development, expect breaking changes until v1.0
factd currently is only fully supported in linux, open to pull requests that change this.

## Environmental Variables

* `HOST_ETC` - specify alternative path to `/etc` directory
* `HOST_PROC` - specify alternative path to `/proc` mountpoint
* `HOST_SYS` - specify alternative path to `/sys` mountpoint


## Adding New Plugins

Resolving the plugins to load is done at the command level. tbd.....

## Factd/Facter Parity

Factd aims to have as much parity with facter's (modern facts)[https://docs.puppet.com/facter/3.3/core_facts.html#modern-facts] set as is reasonable,
note that this does not extend to the names or structure of the facts returned, only the feature set.

| facter                | factd   | notes                 |
|-----------------------|---------|-----------------------|
| augeas                |         |                       |
| disks                 | disks   |                       |
| dmi                   |         |                       |
| ec2_metadata          |         |                       |
| ec2_userdata          |         |                       |
| env_windows_installdir|         | windows not supported |
| facterversion         | version |                       |
| filesystems           | disks   |                       |
| gce                   |         |                       |
| identity              | user    |                       |
| is_virtual            | host    |                       |
| kernel                | host    |                       |
| kernelmajversion      | host    |                       |
| kernelrelease         | host    |                       |
| kernelversion         | host    |                       |
| ldom                  |         | solaris not supported |
| load_averages         | load    |                       |
| memory                | mem     |                       |
| mountpoints           | disks   |                       |
| networking            | net     | partially implemented |
| os                    | host    |                       |
| partitions            | disks   |                       |
| path                  | host    |                       |
| processors            | cpu     |                       |
| ruby                  |         |                       |
| solaris_zones         |         | solaris not supported |
| ssh                   |         |                       |
| system_profiler       |         | osx not supported     |
| system_uptime         | host    |                       |
| timezone              | host    |                       |
| virtual               | host    |                       |
| xen                   |         |                       |
| zfs_featurenumbers    |         | solaris not supported |
| zfs_version           |         | solaris not supported |
| zpool_featurenumbers  |         | solaris not supported |
| zpool_version         |         | solaris not supported |

## Additional Facts

In addition to the facter facts mentioned above factd provides

| name  | notes |
|-------|-------|
| docker| currently only running container and image information is returned      |

## Monitoring

`/metrics` Prometheus Endpoint

`/debug/pprof` pprof Endpoints
