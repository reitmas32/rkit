package criteria

import "regexp"

// MaxIdentifierLength is the maximum accepted length for a field identifier.
const MaxIdentifierLength = 128

// identifierPattern matches a safe, optionally table-qualified SQL/Mongo
// identifier: a letter or underscore followed by letters, digits or
// underscores, with at most one dotted qualifier (e.g. "users.email").
//
// Because the charset excludes spaces, quotes, parentheses, commas and
// semicolons, a value matching this pattern cannot carry SQL/NoSQL injection
// even when interpolated into a query as an identifier.
var identifierPattern = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*(\.[A-Za-z_][A-Za-z0-9_]*)?$`)

// IsValidIdentifier reports whether field is a safe identifier (see
// identifierPattern). Repository implementations use this as the secure default
// when no explicit allow-list is configured.
func IsValidIdentifier(field string) bool {
	return field != "" && len(field) <= MaxIdentifierLength && identifierPattern.MatchString(field)
}

// FieldPolicy controls which field names a repository will accept in filters
// and sorts. The zero value is safe: with an empty Allowed list, every field is
// validated against IsValidIdentifier, which blocks injection while keeping the
// repository usable without extra configuration.
//
// Set Allowed to lock the repository down to a known set of columns (the
// recommended setting when field names may originate from client input).
type FieldPolicy struct {
	// Allowed, when non-empty, is the exhaustive set of acceptable field names.
	// Names are matched verbatim. When empty, IsValidIdentifier is used instead.
	Allowed []string
}

// Permits reports whether field is acceptable under the policy. It always
// requires a syntactically safe identifier; if Allowed is non-empty the field
// must additionally be a member of that set.
func (p FieldPolicy) Permits(field string) bool {
	if !IsValidIdentifier(field) {
		return false
	}
	if len(p.Allowed) == 0 {
		return true
	}
	for _, a := range p.Allowed {
		if a == field {
			return true
		}
	}
	return false
}
