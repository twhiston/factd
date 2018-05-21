# factd
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/85ca8315d2884dbca0b716a800310103)](https://www.codacy.com?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=twhiston/factd&amp;utm_campaign=Badge_Grade)
[![Codacy Badge](https://api.codacy.com/project/badge/Coverage/85ca8315d2884dbca0b716a800310103)](https://www.codacy.com?utm_source=github.com&utm_medium=referral&utm_content=twhiston/factd&utm_campaign=Badge_Coverage)
[![Go Report Card](https://goreportcard.com/badge/github.com/twhiston/factd)](https://goreportcard.com/report/github.com/twhiston/factd)
[![pipeline status](https://gitlab.com/twhiston/factd/badges/master/pipeline.svg)](https://gitlab.com/twhiston/factd/commits/master)

Factd is a fact reporting daemon, similar to puppetlabs facte.
It is intended to be run as a process on the target system, via systemd.
Note that in a container it will need some extra privilages and mounts to be able to access the host information

Facts are served to an http endpoint, as well as prometheus metrics and pprof data.

factd is under heavy development, expect breaking changes until v1.0
factd currently is only fully supported in linux, please submit pull requests that change this.

## Config

All command line options can also be given in a `factd.yml` config file or as env vars prepended with `FACTD_`

## Environmental Variables

As factd uses [gopsutil](https://github.com/shirou/gopsutil) under the hood for some functionality you should also be aware of it's env vars

* `HOST_ETC` - specify alternative path to `/etc` directory
* `HOST_PROC` - specify alternative path to `/proc` mountpoint
* `HOST_SYS` - specify alternative path to `/sys` mountpoint
* `HOST_VAR` - specify alternative path to `/var` mountpoint

## Containers

Although a container is provided it is not the preferred way to run factd.

There are some things to be aware of when running in a container:

* It is a super privilaged container in every way!
* The results you get will differ from running it on the host in a few ways
    * Host filesystem mounted at /host and env vars set to use this
    * Extra volumes
    * Different $PATH
    * Different Host Id
    * Different User
    * IsVirtual true in container (docker only)
    * VirtualizationRole guest in container (docker only)
    * VirtualizationSystem set in container (docker only)

Check the taskfile for running and building the container with the awesome, and daemonless, [podman](https://github.com/projectatomic/libpod)

## Adding New Plugins/Formatters

Use the tmpl command to generate new plugin boilerplate.
See tmpl --help for more details.

## Adding New Commands

Use [cobra](https://github.com/spf13/cobra) for command boilerplate
When adding new commands bind your (p)flags to viper and use viper and not the cmd for data, this allows the user to supply a configuration file supplying flag values. gopsutil env vars are also available from viper.

## Factd/Facter Parity

Factd aims to have as much parity with facter's [modern facts](https://docs.puppet.com/facter/3.3/core_facts.html#modern-facts) feature set as is reasonable,
note that this does not extend to the names or structure of the facts returned.

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
| docker| containers & images |

## Monitoring

`/metrics` Prometheus Endpoint

You can disable prometheus monitoring by passing `--prometheus=false`

`/debug/pprof` pprof Endpoints

You can disable pprof by passing `--pprof=false`

`/healthz` Endpoint

