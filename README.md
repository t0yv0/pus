# pus

Loads and explores [Pulumi Package Schema](https://www.pulumi.com/docs/using-pulumi/pulumi-packages/schema/) files.

## Getting Started

Execute `pus` in a context of a checked out repo such as `pulumi/pulumi-cloudflare`. Try auto-complete on functions:

    » $schema
    <schema:97rs/15fn/267ty>

    » $schema fn
    <functions>

    » $schema fn cloudflare:index/get<TAB>
    cloudflare:index/getZone:getZone
    cloudflare:index/getLists:getLists
    cloudflare:index/getRecord:getRecord
    cloudflare:index/getRulesets:getRulesets
    cloudflare:index/getIpRanges:getIpRanges
    cloudflare:index/getAccountRoles:getAccountRoles
    cloudflare:index/getAccessIdentityProvider:getAccessIdentityProvider
    cloudflare:index/getApiTokenPermissionGroups:getApiTokenPermissionGroups

Inspecting an object further:

    » $schema fn cloudflare:index/getZones:getZones schema outputs properties filter
    $ref: '#/types/cloudflare:index/getZonesFilter:getZonesFilter'
    description: |
        One or more values used to look up zone records. If more than one value is given all values must match in order to be included.
    properties:
        accountId:
            description: |
                The account identifier to target for the resource.
            type: string
        lookupType:
            description: |
                The type of search to perform for the `name` value when querying the zone API. Available values: `contains`, `exact`. Defaults to `exact`.
            type: string
        match:
            description: |
                A RE2 compatible regular expression to filter the   results. This is performed client side whereas the `name` and `lookup_type`     are performed on the Cloudflare server side.
            type: string
        name:
            description: |
                A string value to search for.
            type: string
        paused:
            description: |
                Paused status of the zone to lookup. Defaults to `false`.
            type: boolean
        status:
            description: |
                Status of the zone to lookup.
            type: string
    type: object

You can save an object of interest such as the located function to a variable:

    » $getZone = $schema fn cloudflare:index/getZone:getZone

## Language

See [complang](https://github.com/t0yv0/complang) for more info on the underlying language.
