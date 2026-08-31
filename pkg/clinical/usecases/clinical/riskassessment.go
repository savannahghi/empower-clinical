package clinical

import (
	"context"
	"fmt"
	"strconv"

	"github.com/savannahghi/serverutils"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/empower-clinical/pkg/clinical/domain"
)

// RecordRiskAssessment records a risk assessment based on the provided parameters.
func (c *ClinicalImpl) RecordRiskAssessment(
	ctx context.Context,
	riskAssessment *domain.FHIRRiskAssessmentInput,
) (*domain.FHIRRiskAssessment, error) {
	assessment, err := c.FHIR.CreateFHIRRiskAssessment(ctx, riskAssessment)
	if err != nil {
		return nil, err
	}

	return assessment.Resource, nil
}

// ListRiskAssessment is used to retrieve risk assessments results based on provided parameters
func (c *ClinicalImpl) ListRiskAssessment(ctx context.Context, searchID string, filter *dto.RiskAssessmentFilterInput, pagination serverutils.PaginationInput) (*dto.RiskAssessmentConnection, error) {
	riskAssessmentSearchParams := map[string]interface{}{
		"_sort":  "-_lastUpdated",
		"_total": "accurate",
	}

	err := filter.RiskAssessmentFilter(riskAssessmentSearchParams)
	if err != nil {
		return nil, err
	}

	identifiers, err := c.BaseExtension.GetTenantIdentifiers(ctx)
	if err != nil {
		return nil, err
	}

	limit, offset, err := pagination.ToLimitOffset()
	if err != nil {
		return nil, fmt.Errorf("invalid pagination input: %w", err)
	}

	riskAssessmentSearchParams["_count"] = strconv.Itoa(limit)

	var riskAssessmentOutput *domain.PagedFHIRRiskAssessment

	if searchID != "" {
		riskAssessmentSearchParams["_getpages"] = searchID
		riskAssessmentSearchParams["_getpagesoffset"] = strconv.Itoa(offset)
		riskAssessmentSearchParams["_pretty"] = "true"
		riskAssessmentSearchParams["_bundletype"] = "searchset"

		riskAssessmentOutput, err = c.FHIR.SearchFHIRRiskAssessment(ctx, searchID, riskAssessmentSearchParams, *identifiers, pagination)
		if err != nil {
			return nil, err
		}
	} else {
		riskAssessmentOutput, err = c.FHIR.SearchFHIRRiskAssessment(ctx, "", riskAssessmentSearchParams, *identifiers, pagination)
		if err != nil {
			return nil, err
		}
	}

	riskAssessmentList := []dto.RiskAssessment{}

	for _, riskAssessment := range riskAssessmentOutput.RiskAssessment {
		riskAssessmentList = append(riskAssessmentList, mapFHIRRiskAssessmentToRiskAssessmentDTO(riskAssessment))
	}

	connection := serverutils.BuildLimitOffsetConnection(riskAssessmentList, offset, limit, riskAssessmentOutput.TotalCount)

	return &dto.RiskAssessmentConnection{
		SearchID:   riskAssessmentOutput.BundleID,
		TotalCount: connection.TotalCount,
		Edges:      dto.ConvertRiskAssessmentEdges(connection.Edges),
		PageInfo:   connection.PageInfo,
	}, nil
}

// GetRiskAssessment is used to fetch the details of a risk assessment resulting from a questionnaire
func (c *ClinicalImpl) GetRiskAssessmentByID(ctx context.Context, id string) (*dto.RiskAssessment, error) {
	riskAssessment, err := c.FHIR.GetFHIRRiskAssessment(ctx, id)
	if err != nil {
		return nil, err
	}

	output := mapFHIRRiskAssessmentToRiskAssessmentDTO(*riskAssessment.Resource)

	return &output, nil
}
