package rpsl

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
)

// ASN is an autonomous system number, 2-byte or 4-byte.
type ASN uint32

// String represents the ASN in RPSL format, e.g. AS65000.
func (a ASN) String() string {
	n := strconv.FormatUint(uint64(a), 10)
	return "AS" + n
}

// UnmarshalBinary parses a byte string to an ASN type.
func (a ASN) UnmarshalBinary(b []byte) (ASN, error) {
	if bytes.HasPrefix(b, []byte{0x41, 0x53}) {
		b = b[2:]
	}
	u, err := strconv.ParseUint(string(b), 10, 64)
	if err != nil {
		err := errors.Join(fmt.Errorf("value '%s' could not be parsed as uint32", string(b)), err)
		return 0, err
	}
	return ASN(u), nil
}

// ASName creates an ASN object name from an ASN uint32 number, e.g. AS65000.
func ASNName(a uint32) string {
	return ASN(a).String()
}

// AutNumMembersOf creates a list of aut-num member-of objects for use by [rpsl.AutNum].
func AutNumMembersOf(vals ...any) []string {
	return ASSetMembers(vals)
}

// AutNum is an RPSL 'aut-num class' object. Routing policies are specified using the aut-num class.
// The as-name attribute is a symbolic name (in RPSL name syntax) of the AS. The import, export and
// default routing policies of the AS are specified using import, export and default attributes
// respectively.
type AutNum struct {
	// AS Number.
	//    *Required
	AutNum ASN `rpsl:"aut-num"`
	// aut-num object name.
	//    *Required
	ASName string `rpsl:"as-name"`
	// Description for the aut-num object.
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
	// Import policy expression. See RFC2622 section 6.1.
	Import string `rpsl:"import,omitempty"`
	// Export policy expression. See RFC2622 section 6.1.
	Export string `rpsl:"export,omitempty"`
	// Multi-protocol import policy expression. See RFC4012 section 2.5.
	MPImport string `rpsl:"mp-import,omitempty"`
	// Multi-protocol export policy expression. See RFC4012 section 2.5.
	MPExport string `rpsl:"mp-export,omitempty"`
	// MemberOf can be a list of other aut-num objects or as-set objects of which this aut-num
	// object is a member.
	MemberOf []string `rpsl:"member-of,omitempty" as:"comma-space"`
	// MembersByRef is a list of maintainer names or the keyword ANY. If this attribute is used,
	// the AS set also includes ASes whose aut-num objects are registered by one of these
	// maintainers and whose member-of attribute refers to the name of this AS set. If the value
	// of a mbrs-by-ref attribute is ANY, any AS object referring to the AS set is a member of the
	// set.  If the mbrs-by-ref attribute is missing, only the ASes listed in the members attribute
	// are members of the set.
	MembersByRef []string `rpsl:"mbrs-by-ref,omitempty" as:"comma-space"`
	// Default routing policies. See RFC 2622 section 6.5.
	Default string `rpsl:"default,omitempty"`
	// Multi-protocol default routing policies. See RFC 4012 section 2.5.
	MPDefault string `rpsl:"mp-default,omitempty"`
	// Private container for extra attributes
	Extra map[string]string `rpsl:"-"`
	// Registry Source. Most registries require this field.
	Source string `rpsl:"source,omitempty"`
}

// Add extra pre-formatted attributes to the aut-num object.
func (a *AutNum) AddExtra(key, value string) {
	if a.Extra == nil {
		a.Extra = make(map[string]string)
	}
	a.Extra[key] = value
}

// String representation of the aut-num in RPSL format. E.g. AS65000.
func (a *AutNum) String() string {
	return a.AutNum.String()
}
