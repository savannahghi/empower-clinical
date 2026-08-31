package common

import (
	"time"

	"github.com/rs/xid"
	"github.com/savannahghi/scalarutils"
	"github.com/savannahghi/empower-clinical/pkg/clinical/domain"
)

// constants and defaults
const (

	// NinentyYears is the number of hours in a (fictional) 90 years
	NinentyYears                  = 614880
	timeFormatStr                 = "2006-01-02T15:04:05+03:00"
	healthCloudIdentifiers        = "healthcloud.identifiers"
	healthCloudIdentifiersVersion = "0.0.1"

	// CreatePatientTopic is the topic ID where patient data is published to
	CreatePatientTopic = "patient.create"

	// VitalsTopicName is the topic for publishing a patient's vital signs
	VitalsTopicName = "vitals.create"

	// AllergyTopicName is the topic for publishing a patient's allergy
	AllergyTopicName = "allergy.create"

	// MedicationTopicName is the topic for publishing a patient's medication
	MedicationTopicName = "medication.create"

	// TestResultTopicName is the topic for publishing a patient's test results
	TestResultTopicName = "test.result.create"

	// TestOrderTopicName is the topic for publishing a patient's test order
	TestOrderTopicName = "test.order.create"

	// OrganizationTopicName is the topic where organization(facility) details are published to
	OrganizationTopicName = "organization.create"

	// TenantTopicName is the topic where program is registered in clinical as a tenant
	TenantTopicName = "mycarehub.tenant.create"

	// SegmentationTopicName topic sends patient segmentation information to slade advantage
	SegmentationTopicName = "patient.segmentation.create"

	// ReferralTopicName is used to create a clinical task for a patient referrals
	ReferralTopicName = "patient.referral.task.create"

	// ReferralReportNotificationTopic is used to send notifications - Sends SMS to patient and Email to the receiving facility
	ReferralReportNotificationTopic = "referral.report.notification"

	// CreatePatientCarePlanTopic is used to create tasks for tracking cylces in patient's regimen
	CreateCarePlanTopic = "patient.careplan.create"

	// CreateCFollowUpTaskTopic is used to create followup task(s)
	FollowUpTaskTopic = "followup.task.create"

	// MedicalDataCount is the count of medical records
	MedicalDataCount = "3"

	LabResultsInterpretationLOINCTerminologyCode = "56850-1"

	// WeightLOINCTerminologyCode is the terminology code for weight
	WeightLOINCTerminologyCode = "29463-7"

	// HeightLOINCTerminologyCode is the terminology code height
	HeightLOINCTerminologyCode = "8302-2"

	// TemperatureLOINCTerminologyCode is the terminology code for temperature
	TemperatureLOINCTerminologyCode = "8310-5"

	// MuacLOINCTerminologyCode is the terminology code for mid-upper arm circumference (right)
	MuacLOINCTerminologyCode = "56072-2"

	// BloodSugarLOINCTerminologyCode is the terminology code for blood sugar (Serum glucose)
	BloodSugarLOINCTerminologyCode = "2345-7"

	// DiastolicBloodPressureLOINCTerminologyCode is the terminology code for diastolic blood pressure
	DiastolicBloodPressureLOINCTerminologyCode = "8462-4"

	// LastMenstrualPeriodLOINCTerminologyCode is the terminology code for last menstrual period
	LastMenstrualPeriodLOINCTerminologyCode = "8665-2"

	// Spoc2LOINCTerminologyCode is the terminology code oxygen saturation
	OxygenSaturationLOINCTerminologyCode = "59408-5"

	// RespiratoryRateLOINCTerminologyCode is the terminology code for respiratory rate
	RespiratoryRateLOINCTerminologyCode = "9279-1"

	// PulseLOINCTerminologyCode is the terminology code for pulse
	PulseLOINCTerminologyCode = "8867-4"

	// BloodPressureLOINCTerminologyCode is the terminology code for blood pressure
	BloodPressureLOINCTerminologyCode = "8480-6"

	// BMILOINCTerminologyCode is the terminology code for Body Mass Index
	BMILOINCTerminologyCode = "39156-5"

	// ViralLoadLOINCTerminologyCode is the terminology code for Viral Load
	ViralLoadLOINCTerminologyCode = "25836-8"

	// CD4CountLOINCTerminologyCode is the terminology code for CD4 Count
	CD4CountLOINCTerminologyCode = "29350-6"

	// ClinicalServiceName defines the service where the topic is created
	ClinicalServiceName = "clinical"

	// MyCareHubServiceName defines the service where some of the topics have been created
	MyCareHubServiceName = "mycarehub"

	// TestTopicName is a topic that is used for testing purposes
	TestTopicName = "pubsub.testing.topic"

	// TopicVersion defines the topic version. That standard one is `v1`
	TopicVersion = "v1"

	// AddFHIRIDToPatientProfile is the topic name where the details to update a patient's FHIR ID will be posted to
	AddFHIRIDToPatientProfile = "patient.fhirid.update"

	// AddFHIRIDToFacility is the topic where details to update a facility's fhir ID will be published to
	AddFHIRIDToFacility = "facility.fhirid.update"

	// AddFHIRIDToProgram is the topic where details to update a program's fhir ID will be published to.
	AddFHIRIDToProgram = "program.fhirid.update"

	// LOINCProgressNoteCode defines LOINC progress note terminology code
	LOINCProgressNoteCode = "11506-3"

	// LOINCAssessmentPlanCode defines LOINC assessment plan note terminology code
	LOINCAssessmentPlanCode = "51847-2"

	// LOINCHistoryOfPresentingIllness defines LOINC history of presenting illness note terminology code
	LOINCHistoryOfPresentingIllness = "8684-3"

	// LOINCSocialHistory defines LOINC social history note terminology code
	LOINCSocialHistory = "29762-2"

	// LOINCFamilyHistory defines LOINC family history note terminology code
	LOINCFamilyHistory = "10157-6"

	// LOINCExamination defines LOINC Examination note terminology code
	LOINCExamination = "29545-1"

	// LOINCPlanOfCare defines LOINC Plan of care note terminology code
	LOINCPlanOfCare = "18776-5"

	// LOINCProviderUnspecifiedProgressNote defines LOINC Provider unspecified progress note terminology code
	LOINCProviderUnspecifiedProgressNote = "11506-3"

	// LOINCChiefComplaint defines patients chief complaints
	LOINCChiefComplaint = "10154-3"

	// LOINCPastSurgeryHistory is the terminology code for a patient surgery history
	LOINCPastSurgeryHistory = "71424-6"

	// ColposcopyLOINCTerminologyCode is the terminology code for colposcopy findings
	ColposcopyLOINCTerminologyCode = "LA28211-3"

	// VIALOINCCode is the terminology code for a VIA test
	VIALOINCCode = "47527-7"

	// VIAPositiveLOINCCode is the terminology code for a positive VIA test
	VIAResultPositiveLOINCCode = "LA6576-8"

	// VIANegativeLOINCCode is the terminology code for a negative VIA test
	VIAResultNegativeLOINCCode = "LA6577-6"

	// VIASuspiciousOfCancerCIELCode is the terminology code for a suspicious cancer VIA test
	VIAResultSuspiciousOfCancerLOINCCode = "LA13813-3"

	// HPVLOINCTerminologyCode is the terminology code used to represent HPV test.
	HPVLOINCTerminologyCode = "73959-9"

	// HPV_PCR_DNATerminologyCode is the terminology code used to represent HPV PCR test.
	HPV_PCR_DNATerminologyCode = "91852-4"

	// HPV_OncoproteinTerminologyCode is the terminology code used to represent Human papillomavirus, mixed subtype.
	HPV_OncoproteinTerminologyCode = "82354-2"

	// PapSmearTerminologyCode is the terminology code used to represent pap smear test.
	PapSmearTerminologyCode = "LA16047-5"

	// ProstateCancerTerminologyCode is the terminology code used to represent prostate cancer
	ProstateCancerTerminologyCode = "15325-4"

	LOINCMedicalRecordCode = "11503-0"

	// Chemotherapy terminology code
	CancerChemoTerminologyCode = "21967-5"

	CIELTerminologySystem = "https://CIELterminology.org"

	// RiskAssessmentCodeSystem reprsents the Risk Probability codesystem
	RiskAssessmentCodeSystem = "http://terminology.hl7.org/CodeSystem/risk-probability"

	LOINCSystem = "http://loinc.org"

	// MammogramTerminologyCode is the terminology code used to represent mammogram results.
	MammogramTerminologyCode = "LA16046-7"

	// BenignNeoplasmOfBreastOfSkinTerminologyCode is the terminology code used to represent benign of skin results.
	BenignNeoplasmOfBreastOfSkinTerminologyCode = "LA6675-8"

	// BiopsyTerminologySystem is the terminology code used to represent Biopsy of cervix.
	BiopsyTerminologySystem = "52121-1"

	// MRITerminologySystem is the terminology code used to represent MRI scan of the breast
	MRITerminologySystem = "30794-2"

	// LeftBreastUltrasoundTerminologySystem is the terminology code used to represent left breast ultrasound scan
	LeftBreastUltrasoundTerminologySystem = "26215-4"

	// RightBreastUltrasoundTerminologySystem is the terminology code used to represent right breast scan
	RightBreastUltrasoundTerminologySystem = "26216-2"

	// ChestUltrasoundTerminologySystem is the terminology code used to represent chest ultrasound
	ChestUltrasoundTerminologySystem = "24630-6"

	// BilateralConceptTerminologySystem is the terminology code used to represent miscellaneous bilateral concepts
	BilateralConceptTerminologySystem = "LA25377-5"

	// BreastExaminationLOINCTerminologySystem is the terminology code used to represent breast examination concept.
	BreastExaminationLOINCTerminologySystem = "32422-8"

	// ReferralLOINCTerminologySystem is the system code used to represent a more general referral in LOINC terminology system
	ReferralLOINCTerminologySystem = "57133-1"

	// ReferForMedicalConsultationCIELTerminology is the system code used to represent CIEL's terminology code for patient seeking medical consultations
	ReferForMedicalConsultationLOINCTerminology = "57133-1"

	// ReferralReasonLOINCode represents the referral reason LOINC code
	ReferralReasonLOINCode = "42349-1"

	// IntraReferralLOINCCode represents the LOINC code for referral within the same facility
	IntraReferralLOINCCode = "LA9328-1"

	// ImmunoHistoChemistryLOINCCode represents IHC's ciel code
	ImmunoHistoChemistryLOINCCode = "55229-9"

	// IHCProgesteroneReceptorLOINCCode represents Progesterone Receptor's ciel code
	IHCProgesteroneReceptorLOINCCode = "85339-0"

	// IHCHer2ReceptorLOINCCode represents HER2's ciel code
	IHCEstrogenReceptorLOINCCode = "85337-4"

	// HER2LOINCCode represents HER2's ciel code
	HER2LOINCCode = "85319-2"

	// Ki67LOINCCode represents Ki67's ciel code
	Ki67LOINCCode = "33055-5"

	// PostCoitalBleedingCIELCode represents PCB's ciel code
	PostCoitalBleedingCIELCode = "129452"

	// LaboratoryOrderLOINCCode represents lab orders concept code
	LaboratoryOrderLOINCCode = "26436-6"

	// TestsOrderedLOINCCode represents a labaratory procedure
	TestsOrderedLOINCCode = "LP7839-6"

	CancerProgressionLOINCCode = "97509-4"

	// WholeBloodTerminologyCode represents the whole blood terminology code for whole blood
	WholeBloodTerminologyCode = "906-8"

	// LOINCPastMedicalSurgeryHistory defines LOINC Past Medical Surgery History LOINC terminology code
	LOINCPastMedicalSurgeryHistory = "71424-6"

	// LOINCGeneralStatusNarrative code is used to represent usecases where EMR general examination ins to be recorded
	LOINCGeneralStatusNarrative = "10210-3"
)

var (
	SystemURLICD11 = "http://id.who.int/icd/release/11/mms"
	SystemURLICHI  = "http://id.who.int/ichi"
)

var (
	FacilitySystem     = scalarutils.URI("http://mycarehub/tenant-identification/facility")
	OrganisationSystem = scalarutils.URI("http://mycarehub/tenant-identification/organisation")
)

// DefaultIdentifier assigns a patient a code to function as their
// medical record number.
func DefaultIdentifier() *domain.FHIRIdentifierInput {
	xid := xid.New().String()
	system := scalarutils.URI(healthCloudIdentifiers)
	version := healthCloudIdentifiersVersion
	userSelected := false

	return &domain.FHIRIdentifierInput{
		Use:   domain.IdentifierUseEnumOfficial,
		Value: xid,
		Type: &domain.FHIRCodeableConceptInput{
			Text: "MR",
			Coding: []*domain.FHIRCodingInput{
				{
					System:       &system,
					Version:      &version,
					Code:         scalarutils.Code(xid),
					Display:      xid,
					UserSelected: &userSelected,
				},
			},
		},
		System: &system,
		Period: DefaultPeriodInput(),
	}
}

// DefaultPeriodInput sets up a period input covering roughly a century from when it's run
func DefaultPeriodInput() *domain.FHIRPeriodInput {
	now := time.Now()
	farFuture := time.Now().Add(time.Hour * NinentyYears)

	startDate := scalarutils.DateTime(now.Format(timeFormatStr))
	endDate := scalarutils.DateTime(farFuture.Format(timeFormatStr))

	return &domain.FHIRPeriodInput{
		Start: &startDate,
		End:   &endDate,
	}
}

// DefaultPeriod sets up a period input covering roughly a century from when it's run
func DefaultPeriod() *domain.FHIRPeriod {
	now := time.Now()
	farFuture := time.Now().Add(time.Hour * NinentyYears)

	return &domain.FHIRPeriod{
		Start: scalarutils.DateTime(now.Format(timeFormatStr)),
		End:   scalarutils.DateTime(farFuture.Format(timeFormatStr)),
	}
}

var LoincSystemURL = scalarutils.URI(LOINCSystem)
