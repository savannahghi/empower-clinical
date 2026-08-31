package base

import (
	"context"
	"fmt"
	"strings"

	"github.com/savannahghi/scalarutils"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/common"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/common/helpers"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/utils"
	"github.com/savannahghi/empower-clinical/pkg/clinical/domain"
)

// ProvisionTenant upserts a tenant in FHIR. This endpoint is idempotent: if the tenant
// already exists it is updated with the supplied payload; if it does not exist it is created
// with the supplied tenant_id used as the FHIR Organization ID.
func (c *BaseImpl) ProvisionTenant(ctx context.Context, input dto.ProvisionTenantInput) (*dto.ProvisionTenantOutput, error) {
	if err := helpers.Validate(input); err != nil {
		return nil, err
	}

	if !input.Status.IsValid() {
		return nil, utils.NewCustomError(
			fmt.Errorf("invalid status: %s", input.Status),
			fmt.Sprintf("%s is not a valid tenant status", input.Status),
		)
	}

	orgInput := buildProvisionOrganizationInput(input)

	// Check whether the tenant already exists.
	existing, err := c.FHIR.GetFHIROrganization(ctx, input.TenantID)

	switch {
	case err != nil && strings.Contains(err.Error(), "HTTP 404"):
		// If not found, create it.
		created, err := c.FHIR.PutFHIROrganization(ctx, input.TenantID, *orgInput)
		if err != nil {
			return nil, fmt.Errorf("failed to create tenant organization: %w", err)
		}

		return extractProvisionOutputFromOrganization(created.Resource), nil

	case err != nil:
		return nil, fmt.Errorf("failed to check for existing tenant: %w", err)

	default:
		// If found, update it to match control plane.
		updated, err := c.FHIR.PutFHIROrganization(ctx, *existing.Resource.ID, *orgInput)
		if err != nil {
			return nil, fmt.Errorf("failed to update tenant organization: %w", err)
		}

		return extractProvisionOutputFromOrganization(updated.Resource), nil
	}
}

// GetTenantProvisioningStatus retrieves the provisioning status for a tenant.
func (c *BaseImpl) GetTenantProvisioningStatus(ctx context.Context, tenantID string) (*dto.ProvisionTenantOutput, error) {
	if tenantID == "" {
		return nil, utils.NewCustomError(fmt.Errorf("tenant_id is required"), "please provide tenant_id")
	}

	result, err := c.FHIR.GetFHIROrganization(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("tenant with ID %s not found: %w", tenantID, err)
	}

	return extractProvisionOutputFromOrganization(result.Resource), nil
}

// buildProvisionOrganizationInput constructs a FHIR Organization input from a ProvisionTenantInput,
// setting the FHIR resource ID to the tenant_id and encoding the legacy identifier and parent ref.
func buildProvisionOrganizationInput(input dto.ProvisionTenantInput) *domain.FHIROrganizationInput {
	active := input.Status == dto.TenantStatusActive
	id := input.TenantID

	org := domain.FHIROrganizationInput{
		ID:           &id,
		ResourceType: "Organization",
		Name:         &input.Name,
		Active:       &active,
	}

	// Tenant ID as an identifier
	org.Identifier = append(org.Identifier, mapIdentifierToFHIRIdentifierInput(
		dto.FHIRID,
		input.TenantID,
	))

	// Legacy identifier (SLADE_CODE, etc.)
	if input.LegacyIdentifierType != "" && input.LegacyIdentifierValue != "" && input.LegacyIdentifierType.IsValid() {
		org.Identifier = append(org.Identifier, mapIdentifierToFHIRIdentifierInput(
			dto.SladeCode,
			input.LegacyIdentifierValue,
		))
	}

	// Parent organisation reference
	if input.ParentID != "" {
		ref := fmt.Sprintf("Organization/%s", input.ParentID)
		org.PartOf = &domain.FHIRReference{
			Reference: &ref,
		}
	}

	return &org
}

// extractProvisionOutputFromOrganization maps a FHIR Organization to a ProvisionTenantOutput.
// The Active flag determines the status: true → ACTIVE, false → INACTIVE.
func extractProvisionOutputFromOrganization(org *domain.FHIROrganization) *dto.ProvisionTenantOutput {
	if org == nil {
		return &dto.ProvisionTenantOutput{
			ProvisioningStatus: dto.TenantProvisioningStatusComplete,
		}
	}

	output := &dto.ProvisionTenantOutput{
		ProvisioningStatus: dto.TenantProvisioningStatusComplete,
	}

	if org.ID != nil {
		output.TenantID = *org.ID
	}

	if org.Name != nil {
		output.Name = *org.Name
	}

	if org.Active != nil {
		if *org.Active {
			output.Status = dto.TenantStatusActive
		} else {
			output.Status = dto.TenantStatusInactive
		}
	}

	return output
}

// This file holds all the business logic for creating a FHIR organization. We have a notion of tenants and facilities
// The tenant ID will be used as a logical partitioning key since we want to show that this data resource belongs to this patient who is part of a certain organization(tenant).
// RegisterTenant is used to create an organisation/tenant in the FHIR stores. The tenant ID will be used for logical
// partitioning of data
func (c *BaseImpl) RegisterTenant(ctx context.Context, input dto.OrganizationInput) (*dto.Organization, error) {
	if len(input.Identifiers) == 0 {
		err := fmt.Errorf("expected at least one tenant identifier")
		message := "please provide at least one identifier"

		return nil, utils.NewCustomError(err, message)
	}

	if input.Name == "" {
		err := fmt.Errorf("expected name to be defined")
		message := "please provide the tenant name"

		return nil, utils.NewCustomError(err, message)
	}

	payload := mapOrganizationInputToFHIROrganizationInput(input)

	organisationPayload, err := c.FHIR.CreateFHIROrganization(ctx, *payload)
	if err != nil {
		return nil, err
	}

	return mapFHIROrganizationToDTOOrganization(organisationPayload.Resource), nil
}

func mapIdentifierToFHIRIdentifierInput(idType dto.OrganizationIdentifierType, value string) *domain.FHIRIdentifierInput {
	identifier := domain.FHIRIdentifierInput{
		Use:    domain.IdentifierUseEnumOfficial,
		System: helpers.CodeSystem(common.OrganisationCodeSystemIdentifier),
		Value:  value,
		Period: common.DefaultPeriodInput(),
		Type: &domain.FHIRCodeableConceptInput{
			Coding: []*domain.FHIRCodingInput{},
		},
	}

	identifierTypeCodingInput := domain.FHIRCodeableConceptInput{
		Coding: []*domain.FHIRCodingInput{
			{
				Code:    scalarutils.Code(idType),
				System:  helpers.CodeSystem(common.OrganisationCodeSystemIdentifier),
				Display: idType.ToString(),
			},
		},
	}
	identifier.Type.Coding = append(identifier.Type.Coding, identifierTypeCodingInput.Coding...)

	return &identifier
}

func mapPhoneNumberToFHIRContactPointInput(phoneNumber string) *domain.FHIRContactPointInput {
	use := domain.ContactPointUseEnumWork
	rank := int64(1)
	phoneSystem := domain.ContactPointSystemEnumPhone

	return &domain.FHIRContactPointInput{
		System: &phoneSystem,
		Value:  &phoneNumber,
		Use:    &use,
		Rank:   &rank,
		Period: common.DefaultPeriodInput(),
	}
}

func mapOrganizationInputToFHIROrganizationInput(organization dto.OrganizationInput) *domain.FHIROrganizationInput {
	active := true

	org := domain.FHIROrganizationInput{
		Name:   &organization.Name,
		Active: &active,
	}

	contact := &domain.FHIROrganizationContactInput{}

	if organization.PhoneNumber != "" {
		telecom := mapPhoneNumberToFHIRContactPointInput(organization.PhoneNumber)

		contact.Telecom = append(contact.Telecom, *telecom)

		org.Contact = append(org.Contact, *contact)
	}

	for _, id := range organization.Identifiers {
		identifier := mapIdentifierToFHIRIdentifierInput(id.Type, id.Value)
		org.Identifier = append(org.Identifier, identifier)
	}

	return &org
}

func mapFHIROrganizationToDTOOrganization(organisation *domain.FHIROrganization) *dto.Organization {
	org := &dto.Organization{
		ID:           *organisation.ID,
		Active:       *organisation.Active,
		Name:         *organisation.Name,
		Identifiers:  make([]dto.OrganizationIdentifier, 0),
		PhoneNumbers: make([]string, 0),
	}

	for _, identifier := range organisation.Identifier {
		org.Identifiers = append(org.Identifiers, dto.OrganizationIdentifier{
			Value: identifier.Value,
		})
	}

	org.PhoneNumbers = append(org.PhoneNumbers, organisation.GetOrganizationContacts()...)

	return org
}

// RegisterFacility creates a facility in FHIR. The facility represents the healthcare provider that a service is using.
// E.g if SladeAdvantage are running their program in Nairobi Hospital, then Nairobi hospital will be the facility in this context.
func (c *BaseImpl) RegisterFacility(ctx context.Context, input dto.OrganizationInput) (*dto.Organization, error) {
	found := true

	for _, identifier := range input.Identifiers {
		if !dto.OrganizationIdentifierType.IsValid(identifier.Type) {
			found = false
			break
		}
	}

	if !found {
		err := fmt.Errorf("at least one identifier of type slade code or mfl code is required")
		message := "please provide a SladeCode or MFLCode identifier"

		return nil, utils.NewCustomError(err, message)
	}

	if input.Name == "" {
		err := fmt.Errorf("expected name to be defined")
		message := "please provide the facility name."

		return nil, utils.NewCustomError(err, message)
	}

	payload := mapOrganizationInputToFHIROrganizationInput(input)

	organisationPayload, err := c.FHIR.CreateFHIROrganization(ctx, *payload)
	if err != nil {
		return nil, err
	}

	return mapFHIROrganizationToDTOOrganization(organisationPayload.Resource), nil
}

// CreatePubsubTenant registers a tenant in this service
func (c *BaseImpl) CreatePubsubTenant(ctx context.Context, data dto.OrganizationInput) error {
	if data.Identifiers[0].Type != dto.OrganizationIdentifierType("MCHProgram") {
		return fmt.Errorf("invalid identifier type %v", data.Identifiers[0].Type)
	}

	organization, err := c.RegisterTenant(ctx, data)
	if err != nil {
		return err
	}

	err = c.Pubsub.NotifyProgramFHIRIDUpdate(ctx, dto.UpdateProgramFHIRID{
		ProgramID:    data.Identifiers[0].Value,
		FHIRTenantID: organization.ID,
	})
	if err != nil {
		return err
	}

	return nil
}
