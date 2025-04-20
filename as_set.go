package rpsl

import (
	"regexp"

	"go.mdl.wtf/rpsl/internal/serialize"
	"go.mdl.wtf/rpsl/internal/value"
)

var startsWithASDash = regexp.MustCompile(`^[AaSs]{2}\-.*$`)
var startsWithAS = regexp.MustCompile(`^[AaSs]{2}.*$`)

// ASSetName creates an as-set member type, e.g. AS-ACME or AS65000.
func ASSetName(name string) value.V {
	if startsWithASDash.MatchString(name) {
		name = name[3:]
	} else if startsWithAS.MatchString(name) {
		name = name[2:]
	}
	return value.V("AS-" + name)
}

// ASSet is an RPSL 'as-set class' object. An as-set defines a set of ASNs, aut-num objects, or
// other as-set objects.
type ASSet struct {
	// Name of the as-set object.
	//    *Required
	ASSet string `rpsl:"as-set"`
	// Description for the as-set object.
	Description Description `rpsl:"descr,omitempty"`
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
	// Members of the set; ASNs, aut-num object names, or other as-set names are accepted.
	//
	// Use rpsl.ASSetMembers, rpsl.ASNName, and rpsl.ASSetName functions to ensure proper formatting.
	Members value.VNewline `rpsl:"members,omitempty"`
	// Private container for extra attributes
	Extra map[string]string `rpsl:"-"`
	// Registry Source. Most registries require this field.
	Source string `rpsl:"source,omitempty"`
}

// Add extra pre-formatted attributes to the as-set object.
func (a *ASSet) AddExtra(key, value string) {
	if a.Extra == nil {
		a.Extra = make(map[string]string)
	}
	a.Extra[key] = value
}

// String representation of the as-set in RPSL format. E.g. AS-ACME.
func (a *ASSet) String() string {
	return string(ASSetName(a.ASSet))
}

// RPSL represents the as-set object in RPSL format.
func (a *ASSet) RPSL() (string, error) {
	return serialize.RPSL(a)
}
