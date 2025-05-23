# `rpsl`
Go Library for Creating [RPSL](https://datatracker.ietf.org/doc/rfc2622/) Objects

[![Go Reference](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=for-the-badge)](https://pkg.go.dev/go.mdl.wtf/rpsl)

[![GitHub Tag](https://img.shields.io/github/v/tag/thatmattlove/go-rpsl?style=for-the-badge&label=Version)](https://github.com/thatmattlove/go-rpsl/tags) [![Test Status](https://img.shields.io/github/actions/workflow/status/thatmattlove/go-rpsl/test.yml?style=for-the-badge)](https://github.com/thatmattlove/go-rpsl/actions/workflows/test.yml) [![Test Coverage](https://img.shields.io/codecov/c/github/thatmattlove/go-rpsl?style=for-the-badge)](https://codecov.io/gh/thatmattlove/go-rpsl)

`rpsl` is a library for creating and serializing Routing Policy Specification Language (RPSL) objects, which are used by Internet Routing Registry (IRR) databases for route origin validation.

This library tends to lean towards compatibility with [ARIN's IRR spec](https://www.arin.net/resources/manage/irr/), but should be compatible with any RPSL-compliant IRR provider.

> [!NOTE]
> While the `import`, `export`, `default` (and `mp-` versions) are provided, they are not validated.
> If there is interest, this library could also house a RPSL-compliant policy builder, but does not currently.

### Reference

- [RFC 2622: Routing Policy Specification Language (RPSL)](https://datatracker.ietf.org/doc/rfc2622/)
- [RFC 4012: Routing Policy Specification Language next generation (RPSLng)](https://datatracker.ietf.org/doc/html/rfc4012)

## Installation

```
go get -d go.mdl.wtf/rpsl
```

## Usage Examples

### `route`

```go
route := &rpsl.Route{
    Route:       "192.0.2.0/24",
    Origin:      65000,
    Description: "test",
    AdminPOC:    "TEST-ADMIN",
    TechPOC:     "TEST-TECH",
    MntBy:       "MNT-TEST",
}
formatted, _ := rpsl.MarshalBinary(&route)
fmt.Println(string(formatted))
/*
route: 192.0.2.0/24
origin: AS65000
descr: test
admin-c: TEST-ADMIN
tech-c: TEST-TECH
mnt-by: MNT-TEST
*/
```

### `route6`

```go
route := &rpsl.Route6{
    Route6:       "2001:db8::/32",
    Origin:      65000,
    Description: "test",
    AdminPOC:    "TEST-ADMIN",
    TechPOC:     "TEST-TECH",
    MntBy:       "MNT-TEST",
}
formatted, _ := rpsl.MarshalBinary(&route)
fmt.Println(string(formatted))
/*
route6: 2001:db8::/32
origin: AS65000
descr: test
admin-c: TEST-ADMIN
tech-c: TEST-TECH
mnt-by: MNT-TEST
*/
```

### `route-set`

```go
route_set := &rpsl.RouteSet{
    RouteSet: "RS-ACME",
    Members:  rpsl.RSMembers(netip.MustParsePrefix("192.0.2.0/24"), "RS-ACME"), // with mixed types
    Members:  []string{"192.0.2.0/24", "RS-ACME"}, // or just a string slice
}
formatted, _ := rpsl.MarshalBinary(&route_set)
fmt.Println(string(formatted))
/*
route-set: RS-ACME
members: 192.0.2.0/24,RS-CORP
*/
```

### `aut-num`

```go
aut_num := &rpsl.AutNum{
    AutNum: rpsl.ASN(65000),
    ASName: "AS-65000",
    MemberOf: rpsl.AutNumMembers( // with mixed types
        rpsl.ASN(65001), // rpsl.ASN type
        65002, // uint32
        "AS-ACME", // as-set as string
    ),
    MemberOf: []string{"65001", "AS65002", "AS-ACME"}, // or just a string slice
}
formatted, _ := rpsl.MarshalBinary(&aut_num)
fmt.Println(string(formatted))
/*
aut-num: AS65000
as-name: AS-65000
member-of: AS65001, AS65002, AS-ACME
*/
```

### `as-set`

```go
as_set := &rpsl.ASSet{
    ASSet:   "AS-ACME",
    Members: rpsl.ASSetMembers( // with mixed types
        rpsl.ASN(65000), // rpsl.ASN type
        65001, // uint32
        "AS-CORP", // AS-Set name
    ),
    Members: []string{"65000", "AS65001", "AS-CORP"}, // or just a string slice
}
formatted, _ := rpsl.MarshalBinary(&as_set)
fmt.Println(string(formatted))
/*
as-set: AS-ACME
members: AS65000
members: AS65001
members: AS-CORP
*/

// ↑ AS-Sets list members as separate lines, this is handled appropriately.
```

### Decode

`rpsl` can also decode an RPSL blob:

```go
b := []byte(`
route: 192.0.2.0/24
origin: AS65000
some-unknown-attribute: ur a n3rd
`)

var route rpsl.Route
_ := rpsl.UnmarshalBinary(b, &route)
fmt.Println(route.Route)
// 192.0.2.0/24
fmt.Println(route.Origin.String())
// AS65000
fmt.Println(route.Extra["some-unknown-attribute"]) // ← unknown attributes are placed into a map[string]string accessible via .Extra
// ur a n3rd
```

![License](https://img.shields.io/github/license/thatmattlove/go-rpsl?color=000&style=for-the-badge)
