package rpsl

import "go.mdl.wtf/rpsl/internal/value"

// RSMembers creates a list of route-set members for use by rpsl.RouteSet.
// Example:
//
//	rpsl.RSMembers(rpsl.RSMember("192.0.2.0/24"), rpsl.RSMember("RS-ACME"))
func RSMembers(vals ...value.V) value.VComma {
	out := make(value.VComma, len(vals))
	copy(out, vals)
	return out
}

// ASSetMembers creates a list of as-set members for use by rpsl.ASSet.
// Example:
//
//	rpsl.ASSetMembers(rpsl.ASNName(65000), rpsl.ASSetName("AS-ACME"))
func ASSetMembers(vals ...value.V) value.VNewline {
	out := make(value.VNewline, len(vals))
	for i := range vals {
		if i == 0 {
			out[i] = vals[i]
		} else {
			out[i] = value.V("members: " + vals[i].String())
		}
	}
	return out
}

// AutNumMembers creates a list of aut-num member-of objects for use by rpsl.AutNum.
// Example:
//
//	rpsl.AutNumMembers(rpsl.ASNName(65001), rpsl.AutNumMember("AS65002"), rpsl.ASSetName("AS-ACME"))
func AutNumMembers(vals ...value.V) value.VCommaSpace {
	out := make(value.VCommaSpace, len(vals))
	copy(out, vals)
	return out
}
