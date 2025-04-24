package rpsl

import (
	"net/netip"
	"regexp"
	"strings"
)

var startsWithRSDash = regexp.MustCompile(`^[RaSs]{2}\-.*$`)
var startsWithRS = regexp.MustCompile(`^[RaSs]{2}.*$`)

// RSName creates an route-set member type, e.g. RS-ACME.
func RSName(name string) string {
	if startsWithRSDash.MatchString(name) {
		name = name[3:]
	} else if startsWithRS.MatchString(name) {
		name = name[2:]
	}
	return "RS-" + name
}

// RSMembers creates a list of route-set members for use by [rpsl.RouteSet].
func RSMembers(vals ...any) []string {
	out := make([]string, 0, len(vals))
	for _, v := range vals {
		switch t := v.(type) {
		case string:
			out = append(out, t)
		case netip.Prefix:
			out = append(out, strings.ToLower(t.String()))
		}
	}
	return out
}

// RouteSet is an RPSL 'route-set class' object. RFC2622 specifies that 'the route-set class is a
// set of route prefixes, not of RPSL route objects'; however, because some registries support the
// use of either prefixes or other route-set objects, both are accepted.
type RouteSet struct {
	// Name of the route-set. Begins with RS- or with an AS that is managed by the organization
	// followed by a colon and RS- (for example, AS65536:RS-ARIZ-SE-5).
	//    *Required
	RouteSet string `rpsl:"route-set"`
	// Description for the route-set object.
	Description string `rpsl:"descr,omitempty" as:"multiline"`
	// Admin Point of Contact handle. For ARIN, this field is the exact POC Handle as shown in
	// Whois/RDAP for the Org ID.
	AdminPOC string `rpsl:"admin-c,omitempty"`
	// Technical Point of Contact handle. For ARIN, this field is the exact POC Handle as shown in
	// Whois/RDAP for the Org ID.
	TechPOC string `rpsl:"tech-c,omitempty"`
	// Maintainer object, the prefix MNT and the Org ID of the organization that configures
	// (maintains) the IRR object and manages the resource that is specified in the route object.
	// It is in the format MNT-OrgID; for example, MNT-EXAMPLECORP.
	MntBy string `rpsl:"mnt-by,omitempty"`
	// Any additional information the creator of the objects wants to provide.
	Remarks string `rpsl:"remarks,omitempty"`
	// Members of the set; IPv4 prefixes or other route-set names are accepted.
	//
	// Use rpsl.RSMembers & rpsl.RSMember functions to ensure proper formatting.
	Members []string `rpsl:"members,omitempty" as:"comma"`
	// Members of the set; IPv4 prefixes, IPv6 prefixes, or other route-set names are accepted.
	//
	// Use rpsl.RSMembers & rpsl.RSMember functions to ensure proper formatting.
	MPMembers []string `rpsl:"mp-members,omitempty" as:"comma"`
	// Private container for extra attributes
	Extra map[string]string `rpsl:"-"`
	// Registry Source. Most registries require this field.
	Source string `rpsl:"source,omitempty"`
}

// Add extra pre-formatted attributes to the route-set object.
func (rs *RouteSet) AddExtra(key, value string) {
	if rs.Extra == nil {
		rs.Extra = map[string]string{key: value}
	}
	rs.Extra[key] = value
}

// String representation of the route-set in RPSL format. E.g. RS-ACME.
func (rs *RouteSet) String() string {
	return RSName(rs.RouteSet)
}
