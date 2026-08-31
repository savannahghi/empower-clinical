package base

import (
	"context"
	"log/slog"
	"sort"
	"strings"

	"github.com/savannahghi/empower-clinical/pkg/clinical/application/dto"
)

// GetPatientTimeline returns a patient's timeline as a minimal, chronologically sorted list of timeline resources.
// It fetches all FHIR resources for the patient, maps them to TimelineResource DTOs, and sorts them by date.
func (c *BaseImpl) GetPatientTimeline(ctx context.Context, patientID string, params *dto.PatientEverythingFilterParams) (*dto.HealthTimeline, error) {
	patientEverything, err := c.GetPatientEverything(ctx, patientID, params)
	if err != nil {
		c.Error("failed to fetch patient everything", slog.String("patientID", patientID), slog.Any("err", err))
		return nil, err
	}

	timeline := c.patientEverythingBuilder(patientEverything)

	return &dto.HealthTimeline{
		Timeline:   timeline,
		TotalCount: len(timeline),
		PageInfo:   patientEverything.PageInfo,
	}, nil
}

func (c *BaseImpl) patientEverythingBuilder(patientEverything *dto.PatientEverythingConnection) []dto.TimelineResource {
	timeline := make([]dto.TimelineResource, 0, len(patientEverything.Edges))

	for i, resource := range patientEverything.Edges {
		if resource["resourceType"] == "Patient" {
			continue
		}

		mapped, err := c.Mapper.ToTimeline(resource)
		if err != nil {
			c.Debug("skipping unmappable resource", slog.Int("index", i), slog.Any("resource", resource), slog.Any("err", err))
			continue
		}

		timeline = append(timeline, *mapped)
	}

	sort.Slice(timeline, func(i, j int) bool {
		d1, d2 := timeline[i].Date, timeline[j].Date

		if d1.Year != d2.Year {
			return d1.Year > d2.Year
		}

		if d1.Month != d2.Month {
			return d1.Month > d2.Month
		}

		if d1.Day != d2.Day {
			return d1.Day > d2.Day
		}

		return timeline[i].TimeRecorded.After(timeline[j].TimeRecorded)
	})

	return timeline
}

// GetPatientBanner returns resources that are sufficient enough to show patient banner info (conditions, allergies & medications) sorted from most recent
func (c *BaseImpl) GetPatientBanner(ctx context.Context, patientID string, params *dto.PatientEverythingFilterParams) (*dto.PatientBanner, error) {
	if params == nil {
		params = &dto.PatientEverythingFilterParams{}
	}

	bannerResources := []string{"Condition", "MedicationStatement", "AllergyIntolerance"}
	params.Type = strings.Join(bannerResources, ",")

	patientEverything, err := c.GetPatientEverything(ctx, patientID, params)
	if err != nil {
		c.Error("failed to fetch patient everything", slog.String("patientID", patientID), slog.Any("err", err))
		return nil, err
	}

	timeline := c.patientEverythingBuilder(patientEverything)

	var (
		conditions   []dto.TimelineResource
		allergies    []dto.TimelineResource
		medications  []dto.TimelineResource
		observations []dto.TimelineResource
	)

	for _, item := range timeline {
		switch item.ResourceType {
		case dto.ResourceTypeCondition:
			conditions = append(conditions, item)
		case dto.ResourceTypeAllergyIntolerance:
			allergies = append(allergies, item)
		case dto.ResourceTypeMedicationStatement:
			medications = append(medications, item)
		case dto.ResourceTypeObservation:
			observations = append(observations, item)
		}
	}

	firstN := func(items []dto.TimelineResource, n int) []dto.TimelineResource {
		if len(items) <= n {
			return items
		}

		return items[:n]
	}

	return &dto.PatientBanner{
		Conditions:  firstN(conditions, 3),
		Allergies:   firstN(allergies, 3),
		Medications: firstN(medications, 3),
		Observation: firstN(observations, 3),
	}, nil
}
