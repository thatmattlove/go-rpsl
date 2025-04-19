package rpsl

import (
	"go.mdl.wtf/rpsl/internal/serialize"
)

// Route is an RPSL 'route class' object. Each interAS route (also referred to as an interdomain
// route) originated by an AS is specified using a route object.
type Route struct {
	// IPv4 address prefix. It is assumed that this value is a valid IPv4 prefix; no validation occurs.
	//    *Required
	Route string `rpsl:"route"`
	// The ASN from which the route originates.
	//    *Required
	Origin ASN `rpsl:"origin"`
	// Description for the route object.
	Description string `rpsl:"descr,omitempty"`
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
	// Private container for extra attributes
	Extra map[string]string `rpsl:"-"`
	// Registry Source. Most registries require this field.
	Source string `rpsl:"source,omitempty"`
}

// Add extra pre-formatted attributes to the route object.
func (r *Route) AddExtra(key, value string) {
	if r.Extra == nil {
		r.Extra = make(map[string]string)
	}
	r.Extra[key] = value
}

// RPSL represents the route object in RPSL format.
func (r *Route) RPSL() (string, error) {
	return serialize.RPSL(r)
}
