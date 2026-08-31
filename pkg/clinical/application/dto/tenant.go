package dto

// TenantStatus represents the status of a tenant as per the control plane
type TenantStatus string

const (
	TenantStatusActive   TenantStatus = "ACTIVE"
	TenantStatusInactive TenantStatus = "INACTIVE"
)

// IsValid checks validity of a TenantStatus enum
func (t TenantStatus) IsValid() bool {
	switch t {
	case TenantStatusActive, TenantStatusInactive:
		return true
	}

	return false
}

// TenantProvisioningStatus represents the provisioning status within this system
type TenantProvisioningStatus string

const (
	TenantProvisioningStatusComplete TenantProvisioningStatus = "COMPLETE"
)

// LegacyIdentifierType represents the type of legacy identifier used for migration/mapping
type LegacyIdentifierType string

const (
	LegacyIdentifierTypeSladeCode LegacyIdentifierType = "SLADE_CODE"
)

// IsValid checks validity of a LegacyIdentifierType enum
func (l LegacyIdentifierType) IsValid() bool {
	return l == LegacyIdentifierTypeSladeCode
}

// ToString returns a human-readable representation of LegacyIdentifierType
func (l LegacyIdentifierType) ToString() string {
	if l == LegacyIdentifierTypeSladeCode {
		return "Slade Code"
	}

	return string(l)
}

// ProvisionTenantInput is the request payload for provisioning a tenant
type ProvisionTenantInput struct {
	TenantID              string               `json:"tenant_id" validate:"required"`
	ParentID              string               `json:"parent_id,omitempty"`
	Name                  string               `json:"name" validate:"required"`
	LegacyIdentifierType  LegacyIdentifierType `json:"legacy_identifier_type"`
	LegacyIdentifierValue string               `json:"legacy_identifier_value"`
	Status                TenantStatus         `json:"status"`
}

// ProvisionTenantOutput is the response payload for tenant provisioning
type ProvisionTenantOutput struct {
	Name               string                   `json:"name"`
	TenantID           string                   `json:"tenant_id"`
	Status             TenantStatus             `json:"status"`
	ProvisioningStatus TenantProvisioningStatus `json:"provisioning_status"`
}
