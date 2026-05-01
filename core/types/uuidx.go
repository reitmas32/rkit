package types

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// Bit allocation for module and entity codes in the UUID.
const (
	NumModuleBits = 8
	NumEntityBits = 8
)

const (
	numModuleHex  = NumModuleBits / 4 // 4 bits = 1 hex char
	numEntityHex  = NumEntityBits / 4
	numTotalHex   = numModuleHex + numEntityHex
	numTotalBytes = numTotalHex / 2
)

// ModuleCode represents a module identifier (hex string).
type ModuleCode string

func (m ModuleCode) String() string { return string(m) }

// EntityCode represents an entity type identifier (hex string).
type EntityCode string

func (e EntityCode) String() string { return string(e) }

// TypedUUID is the interface for UUIDs with embedded module/entity codes.
type TypedUUID interface {
	Module() ModuleCode
	Entity() EntityCode
	UUID() uuid.UUID
	String() string
}

// DomainID is the default implementation of TypedUUID.
type DomainID struct {
	raw    uuid.UUID
	module ModuleCode
	entity EntityCode
}

// Module returns the module code.
func (d DomainID) Module() ModuleCode { return d.module }

// Entity returns the entity code.
func (d DomainID) Entity() EntityCode { return d.entity }

// UUID returns the underlying UUID.
func (d DomainID) UUID() uuid.UUID { return d.raw }

// String returns the UUID string representation.
func (d DomainID) String() string { return d.raw.String() }

// IsZero returns true if the DomainID is empty/nil.
func (d DomainID) IsZero() bool { return d.raw == uuid.Nil }

// NewDomainID creates a new DomainID with the given module and entity codes.
// Both moduleHex and entityHex must be 2-character hex strings.
func NewDomainID(moduleHex, entityHex string) (DomainID, error) {
	moduleHex = strings.ToLower(strings.TrimSpace(moduleHex))
	entityHex = strings.ToLower(strings.TrimSpace(entityHex))

	if len(moduleHex) != numModuleHex {
		return DomainID{}, fmt.Errorf("module must be %d hex chars, got %d", numModuleHex, len(moduleHex))
	}
	if len(entityHex) != numEntityHex {
		return DomainID{}, fmt.Errorf("entity must be %d hex chars, got %d", numEntityHex, len(entityHex))
	}

	u, err := makeTypedUUID(moduleHex, entityHex)
	if err != nil {
		return DomainID{}, err
	}

	m, e := ExtractModuleEntity(u)
	return DomainID{raw: u, module: m, entity: e}, nil
}

// ParseDomainID parses a UUID string and extracts the module/entity codes.
func ParseDomainID(s string) (DomainID, error) {
	u, err := uuid.Parse(s)
	if err != nil {
		return DomainID{}, fmt.Errorf("invalid UUID: %w", err)
	}
	m, e := ExtractModuleEntity(u)
	return DomainID{raw: u, module: m, entity: e}, nil
}

// ExtractModuleEntity extracts module and entity codes from a UUID.
func ExtractModuleEntity(id uuid.UUID) (ModuleCode, EntityCode) {
	h := fmt.Sprintf("%x", id[16-numTotalBytes:])
	return ModuleCode(h[:numModuleHex]), EntityCode(h[numModuleHex:numTotalHex])
}

// makeTypedUUID creates a UUID v4 with embedded module/entity codes in the last bytes.
func makeTypedUUID(moduleHex, entityHex string) (uuid.UUID, error) {
	suffix := moduleHex + entityHex
	pb, err := hex.DecodeString(suffix)
	if err != nil || len(pb) != numTotalBytes {
		return uuid.Nil, fmt.Errorf("invalid suffix hex: %w", err)
	}

	u := uuid.New()

	// Overwrite last bytes with module+entity
	copy(u[16-numTotalBytes:], pb)

	// Ensure UUID v4 version and variant bits
	u[6] = (u[6] & 0x0F) | 0x40 // version 4
	u[8] = (u[8] & 0x3F) | 0x80 // variant RFC 4122

	return u, nil
}
