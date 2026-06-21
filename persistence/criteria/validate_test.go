package criteria_test

import (
	"strings"
	"testing"

	"github.com/reitmas32/rkit/persistence/criteria"
)

func TestIsValidIdentifier_AcceptsSafeNames(t *testing.T) {
	valid := []string{"id", "user_id", "created_at", "_private", "Col1", "users.email"}
	for _, f := range valid {
		if !criteria.IsValidIdentifier(f) {
			t.Errorf("expected %q to be a valid identifier", f)
		}
	}
}

func TestIsValidIdentifier_RejectsInjection(t *testing.T) {
	// Each of these is a classic SQL/NoSQL identifier-injection attempt and must
	// be rejected so it can never be interpolated into a query.
	bad := []string{
		"",
		"1=1",
		"id; DROP TABLE users",
		"id OR 1=1",
		"id)",
		"name'",
		"col--",
		"$where",
		"a b",
		"(SELECT 1)",
		"id,name",
		"a.b.c",
		strings.Repeat("a", criteria.MaxIdentifierLength+1),
	}
	for _, f := range bad {
		if criteria.IsValidIdentifier(f) {
			t.Errorf("expected %q to be rejected as an identifier", f)
		}
	}
}

func TestFieldPolicy_Permits(t *testing.T) {
	// Empty allow-list: any syntactically safe identifier is permitted, but
	// injection is still rejected.
	open := criteria.FieldPolicy{}
	if !open.Permits("email") {
		t.Error("open policy should permit a safe identifier")
	}
	if open.Permits("id; DROP TABLE users") {
		t.Error("open policy must still reject injection")
	}

	// Explicit allow-list: only listed fields pass, and they must also be safe.
	locked := criteria.FieldPolicy{Allowed: []string{"email", "status"}}
	if !locked.Permits("email") {
		t.Error("locked policy should permit an allowed field")
	}
	if locked.Permits("password") {
		t.Error("locked policy must reject a non-allowed field")
	}
	if locked.Permits("email OR 1=1") {
		t.Error("locked policy must reject injection even if it contains an allowed word")
	}
}
