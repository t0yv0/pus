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

You can save an object of interest such as the located function to a variable:

    » $getZone = $schema fn cloudflare:index/getZone:getZone

## Language

See [complang](https://github.com/t0yv0/complang) for more info on the underlying language.
