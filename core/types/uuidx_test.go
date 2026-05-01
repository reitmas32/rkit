package types

import (
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestModuleCode_String(t *testing.T) {
	tests := []struct {
		name     string
		code     ModuleCode
		expected string
	}{
		{"simple code", ModuleCode("ab"), "ab"},
		{"empty code", ModuleCode(""), ""},
		{"long code", ModuleCode("abcdef"), "abcdef"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.code.String(); got != tt.expected {
				t.Errorf("ModuleCode.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestEntityCode_String(t *testing.T) {
	tests := []struct {
		name     string
		code     EntityCode
		expected string
	}{
		{"simple code", EntityCode("cd"), "cd"},
		{"empty code", EntityCode(""), ""},
		{"long code", EntityCode("cdefgh"), "cdefgh"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.code.String(); got != tt.expected {
				t.Errorf("EntityCode.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestDomainID_Module(t *testing.T) {
	did := DomainID{
		raw:    uuid.New(),
		module: ModuleCode("ab"),
		entity: EntityCode("cd"),
	}

	if got := did.Module(); got != ModuleCode("ab") {
		t.Errorf("DomainID.Module() = %v, want %v", got, "ab")
	}
}

func TestDomainID_Entity(t *testing.T) {
	did := DomainID{
		raw:    uuid.New(),
		module: ModuleCode("ab"),
		entity: EntityCode("cd"),
	}

	if got := did.Entity(); got != EntityCode("cd") {
		t.Errorf("DomainID.Entity() = %v, want %v", got, "cd")
	}
}

func TestDomainID_UUID(t *testing.T) {
	testUUID := uuid.New()
	did := DomainID{
		raw:    testUUID,
		module: ModuleCode("ab"),
		entity: EntityCode("cd"),
	}

	if got := did.UUID(); got != testUUID {
		t.Errorf("DomainID.UUID() = %v, want %v", got, testUUID)
	}
}

func TestDomainID_String(t *testing.T) {
	testUUID := uuid.New()
	did := DomainID{
		raw:    testUUID,
		module: ModuleCode("ab"),
		entity: EntityCode("cd"),
	}

	if got := did.String(); got != testUUID.String() {
		t.Errorf("DomainID.String() = %v, want %v", got, testUUID.String())
	}
}

func TestDomainID_IsZero(t *testing.T) {
	tests := []struct {
		name     string
		did      DomainID
		expected bool
	}{
		{"zero UUID", DomainID{raw: uuid.Nil}, true},
		{"non-zero UUID", DomainID{raw: uuid.New()}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.did.IsZero(); got != tt.expected {
				t.Errorf("DomainID.IsZero() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestNewDomainID(t *testing.T) {
	tests := []struct {
		name        string
		moduleHex   string
		entityHex   string
		wantErr     bool
		checkModule bool
		checkEntity bool
	}{
		{"valid codes", "ab", "cd", false, true, true},
		{"valid codes with spaces", " ab ", " cd ", false, true, true},
		{"valid codes uppercase", "AB", "CD", false, true, true},
		{"invalid module length", "a", "cd", true, false, false},
		{"invalid entity length", "ab", "c", true, false, false},
		{"empty module", "", "cd", true, false, false},
		{"empty entity", "ab", "", true, false, false},
		{"module too long", "abc", "cd", true, false, false},
		{"entity too long", "ab", "cde", true, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewDomainID(tt.moduleHex, tt.entityHex)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDomainID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if tt.checkModule && got.Module() != ModuleCode(strings.ToLower(strings.TrimSpace(tt.moduleHex))) {
					t.Errorf("NewDomainID() module = %v, want %v", got.Module(), tt.moduleHex)
				}
				if tt.checkEntity && got.Entity() != EntityCode(strings.ToLower(strings.TrimSpace(tt.entityHex))) {
					t.Errorf("NewDomainID() entity = %v, want %v", got.Entity(), tt.entityHex)
				}
				if got.IsZero() {
					t.Error("NewDomainID() returned zero DomainID")
				}
			}
		})
	}
}

func TestParseDomainID(t *testing.T) {
	tests := []struct {
		name    string
		uuidStr string
		wantErr bool
	}{
		{"valid UUID", uuid.New().String(), false},
		{"invalid UUID", "not-a-uuid", true},
		{"empty string", "", true},
		{"malformed UUID", "12345", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDomainID(tt.uuidStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDomainID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.String() != tt.uuidStr {
					t.Errorf("ParseDomainID() UUID = %v, want %v", got.String(), tt.uuidStr)
				}
				if got.IsZero() {
					t.Error("ParseDomainID() returned zero DomainID")
				}
			}
		})
	}
}

func TestExtractModuleEntity(t *testing.T) {
	// Create a UUID and extract module/entity
	testUUID := uuid.New()
	module, entity := ExtractModuleEntity(testUUID)

	// Verify they are 2 characters each (numModuleHex = 2, numEntityHex = 2)
	if len(module) != numModuleHex {
		t.Errorf("ExtractModuleEntity() module length = %v, want %v", len(module), numModuleHex)
	}
	if len(entity) != numEntityHex {
		t.Errorf("ExtractModuleEntity() entity length = %v, want %v", len(entity), numEntityHex)
	}

	// Test with zero UUID
	moduleZero, entityZero := ExtractModuleEntity(uuid.Nil)
	if len(moduleZero) != numModuleHex {
		t.Errorf("ExtractModuleEntity() module length (zero) = %v, want %v", len(moduleZero), numModuleHex)
	}
	if len(entityZero) != numEntityHex {
		t.Errorf("ExtractModuleEntity() entity length (zero) = %v, want %v", len(entityZero), numEntityHex)
	}
}

func TestMakeTypedUUID(t *testing.T) {
	tests := []struct {
		name      string
		moduleHex string
		entityHex string
		wantErr   bool
	}{
		{"valid hex", "ab", "cd", false},
		{"invalid hex", "xy", "cd", true},
		{"invalid suffix", "ab", "xy", true},
		{"empty module", "", "cd", true},
		{"empty entity", "ab", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := makeTypedUUID(tt.moduleHex, tt.entityHex)
			if (err != nil) != tt.wantErr {
				t.Errorf("makeTypedUUID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got == uuid.Nil {
					t.Error("makeTypedUUID() returned nil UUID")
				}
				// Verify version and variant bits
				version := (got[6] & 0xF0) >> 4
				if version != 4 {
					t.Errorf("makeTypedUUID() version = %v, want 4", version)
				}
				variant := (got[8] & 0xC0) >> 6
				if variant != 2 {
					t.Errorf("makeTypedUUID() variant = %v, want 2 (RFC 4122)", variant)
				}
			}
		})
	}
}

func TestNewDomainID_ExtractRoundTrip(t *testing.T) {
	moduleHex := "ab"
	entityHex := "cd"

	did, err := NewDomainID(moduleHex, entityHex)
	if err != nil {
		t.Fatalf("NewDomainID() error = %v", err)
	}

	extractedModule, extractedEntity := ExtractModuleEntity(did.UUID())
	if extractedModule != ModuleCode(moduleHex) {
		t.Errorf("ExtractModuleEntity() module = %v, want %v", extractedModule, moduleHex)
	}
	if extractedEntity != EntityCode(entityHex) {
		t.Errorf("ExtractModuleEntity() entity = %v, want %v", extractedEntity, entityHex)
	}
}
