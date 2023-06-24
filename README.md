# pus

Loads and explores [Pulumi Package Schema](https://www.pulumi.com/docs/using-pulumi/pulumi-packages/schema/) files.

## Getting Started

Execute `pus` in a context of a checked out repo such as `pulumi/pulumi-cloudflare`. Try auto-complete on functions:

    » schema fn cloudflare:index/get[TAB]
    cloudflare:index/getZone:getZone
    cloudflare:index/getLists:getLists
    cloudflare:index/getRecord:getRecord
    cloudflare:index/getRulesets:getRulesets
    cloudflare:index/getIpRanges:getIpRanges
    cloudflare:index/getAccountRoles:getAccountRoles
    cloudflare:index/getAccessIdentityProvider:getAccessIdentityProvider
    cloudflare:index/getApiTokenPermissionGroups:getApiTokenPermissionGroups

You can save the function of interest to a variable:

    » getZone = schema fn cloudflare:index/getZone:getZone

## Language

Executing `pus` opens a REPL where schema operations are embedded in a simple language `pus-lang`. The language is
optimized for completion. As you explore, press TAB at any point to code-complete the current expression.

Expressions are space-separted tokens interpreted as message sends as in Smalltalk, that is `a b c` is interpreted as
`a.b().c()` blub. Expressions operate on one flat mutable global environment.

Strings are self-evaluating:

    » something
    something

Assignment evaluates to the expression on the right, modifying the global environment:

    » myvar = myvalue
    myvalue

A bound expression resolves to its value, and no longer self-evaluates:

    » myvar
    myvalue
