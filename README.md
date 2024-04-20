# pus

Loads and explores [Pulumi Package Schema](https://www.pulumi.com/docs/using-pulumi/pulumi-packages/schema/) files.

## Getting Started

Execute `pus` in a context of a checked out repo such as `pulumi/pulumi-cloudflare`. The built-in
`$schema` loads the package schema and allows auto-complete to explore it:

    » $schema <TAB>
    Meta       Diff       Name       Types      Config     Version    License    LogoURL    DiffTag
    Provider   Homepage   Keywords   Language   Functions  Publisher  Resources  Repository

    » $schema Functions
    cloudflare:index/getAccessApplication:getAccessApplication: {}
    cloudflare:index/getAccessIdentityProvider:getAccessIdentityProvider: {}
    cloudflare:index/getAccountRoles:getAccountRoles: {}
    cloudflare:index/getAccounts:getAccounts: {}
    cloudflare:index/getApiTokenPermissionGroups:getApiTokenPermissionGroups: {}
    cloudflare:index/getDevicePostureRules:getDevicePostureRules: {}
    ...

    » $schema Functions cloudflare:index/get<TAB>
    cloudflare:index/getZone:getZone
    cloudflare:index/getUser:getUser
    cloudflare:index/getList:getList
    cloudflare:index/getLists:getLists
    cloudflare:index/getZones:getZones
    cloudflare:index/getTunnel:getTunnel
    cloudflare:index/getRecord:getRecord
    ...

If your workspace has schema edit relative to HEAD, `pus` can help explore the diff hierarchically:

    » $schema Diff resources <TAB>
    cloudflare:index/argo:Argo                    cloudflare:index/zone:Zone
    cloudflare:index/listItem:ListItem            cloudflare:index/zoneHold:ZoneHold
    cloudflare:index/greTunnel:GreTunnel          cloudflare:index/teamsRule:TeamsRule
    ...

Diffs can be drilled down further:

    » $schema Diff resources cloudflare:index/argo:Argo
    inputProperties:
        tieredCaching:
            deprecationMessage:
                rm:
                    text: tiered_caching has been deprecated in favour of using `cloudflare_tiered_cache` resource instead.

Diffs can be performed against a specific tag, for example:

    » $schema DiffTag <TAB>
    v5.2.0 v5.4.0 v5.1.1 v2.2.0 v4.1.0 v5.4.1 v4.5.0 v1.9.0 v2.9.0 v2.0.0 v1.6.0 v3.2.0 v5.5.0 v5.1.0 v3.4.0 v2.4.0 v5.3.0

Every object can be explored further down with auto-complete:

    » $schema Functions cloudflare:index/getZones:getZones outputs properties filter
    $ref:
        text: '#/types/cloudflare:index/getZonesFilter:getZonesFilter'
    description:
        text: |
            One or more values used to look up zone records. If more than one value is given all values must match in order to be in...
    properties:
        accountId:
            description:
                text: |
                    The account identifier to target for the resource.
            type:
                text: string
        lookupType:
            description:
                text: |
                    The type of search to perform for the `name` value when querying the zone API. Available values: `contains`, `ex...
            type:
                text: string
        match:
            description:
                text: |
                    A RE2 compatible regular expression to filter the       results. This is performed client side whereas the `name` and ...
            type:
                text: string
        name:
            description:
                text: |
                    A string value to search for.
            type:
                text: string
        paused:
            description:
    ...

You can save an object of interest such as the located function to a variable:

    » $getZone = $schema Functions cloudflare:index/getZone:getZone

## Features

### Ref Inlining

Pulumi Package Schema supports references to named types, which can come from the local or external
package. References are also critical for working recursive types.

```
$ref: '#/types/cloudflare:index/getZonesFilter:getZonesFilter'
```

Currently `pus` viewer inlines the content of local references while also retaining the `$ref` data.
This makes it easier to view and compare nested object schemata.

Inlining non-local references and proper support for cyclic types is not implemented yet.


### CLI Interface for Completions

To make it easy to build lightweight editor integration, a CLI interface exposes all completion candidates. The idea is
that the editor can then do further interactive filtering of the candidates.


    > ./pus --complete "\$sch"
    $schema

    > ./pus --complete "\$schema x"
    AllowedPackageNames
    Attribution
    Config
    Description
    ...

    > ./pus --complete "\$schema Types x"
    awsx:awsx:Bucket
    awsx:awsx:DefaultBucket
    awsx:awsx:DefaultLogGroup
    awsx:awsx:DefaultRoleWithPolicy
    awsx:awsx:DefaultSecurityGroup
    ...


## Language

See [complang](https://github.com/t0yv0/complang) for more info on the underlying language.
