package helpers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/asaskevich/govalidator"
	"github.com/go-playground/validator"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/common"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/empower-clinical/pkg/clinical/domain"
	"golang.org/x/net/html"

	"github.com/savannahghi/converterandformatter"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/scalarutils"
	"github.com/savannahghi/serverutils"
	log "github.com/sirupsen/logrus"
)

var (
	baseURL        = serverutils.MustGetEnvVar("HAPI_FHIR_BASE_URL")
	QuantitySystem = "http://unitsofmeasure.org"
)

const (
	// NinentyYears is the number of hours in a (fictional) century of leap years
	NinentyYears       = 878400
	fullAccessLevel    = "FULL_ACCESS"
	partialAccessLevel = "PROFILE_AND_RECENT_VISITS_ACCESS"
	timeFormatStr      = "2006-01-02T15:04:05+03:00"
	// StringTimeParseMonthNameLayout ...
	StringTimeParseMonthNameLayout = "2006-Jan-02"
	// StringTimeParseMonthNumberLayout ...
	StringTimeParseMonthNumberLayout   = "2006-01-02"
	StringTimeParseMonthFullNameLayout = "January 2, 2006"
)

// simple patient registration defaults
// should be reviewed later (ticket created)
const (
	DefaultCountry = "ke"
	DefaultVersion = "0.0.1"
)

// ComposeOneHealthEpisodeOfCare is used to create an episode of care
func ComposeOneHealthEpisodeOfCare(
	validPhone string, fullAccess bool, organizationID, providerCode, patientID string,
) domain.FHIREpisodeOfCare {
	accessLevel := ""
	if fullAccess {
		accessLevel = fullAccessLevel
	} else {
		accessLevel = partialAccessLevel
	}

	now := time.Now()
	farFuture := time.Now().Add(time.Hour * NinentyYears)
	orgIdentifier := &domain.FHIRIdentifier{
		Use:   "official",
		Value: providerCode,
	}
	active := domain.EpisodeOfCareStatusEnumActive
	orgRef := fmt.Sprintf("Organization/%s", organizationID)
	patientRef := fmt.Sprintf("Patient/%s", patientID)
	orgType := scalarutils.URI("Organization")
	patientType := scalarutils.URI("Patient")

	return domain.FHIREpisodeOfCare{
		Status: &active,
		Period: &domain.FHIRPeriod{
			Start: scalarutils.DateTime(now.Format(timeFormatStr)),
			End:   scalarutils.DateTime(farFuture.Format(timeFormatStr)),
		},
		ManagingOrganization: &domain.FHIRReference{
			Reference:  &orgRef,
			Display:    providerCode,
			Type:       &orgType,
			Identifier: orgIdentifier,
		},
		Patient: &domain.FHIRReference{
			Reference: &patientRef,
			Display:   validPhone,
			Type:      &patientType,
		},
		Type: []*domain.FHIRCodeableConcept{
			{
				Text: accessLevel, // FULL_ACCESS or PROFILE_AND_RECENT_VISITS_ACCESS
			},
		},
	}
}

// ParseDate parses a formatted string and returns the time value it represents
// Try different layout due to Inconsistent date formats
// They include "2018-01-01", "2020-09-24T18:02:38.661033Z", "2018-Jan-01"
func ParseDate(date string) time.Time {
	// parses "2020-09-24T18:02:38.661033Z"
	timeValue, err := time.Parse(time.RFC3339, date)
	if err != nil {
		// parses when month is a number e.g "2018-01-01"
		timeValue, err := time.Parse(StringTimeParseMonthNameLayout, date)
		if err != nil {
			// parses when month is a name e.g "2018-Jan-01"
			timeValue, err := time.Parse(StringTimeParseMonthNumberLayout, date)
			if err != nil {
				// parses when month is a name e.g "January 2, 2006"
				timeValue, err := time.Parse(StringTimeParseMonthFullNameLayout, date)
				if err != nil {
					log.Errorf("cannot parse date %v with error %v", date, err)
				}

				return timeValue
			}

			return timeValue
		}

		return timeValue
	}

	return timeValue
}

// IDToIdentifier translates simple identification
// document details to FHIR identifiers
func IDToIdentifier(
	ids []*domain.IdentificationDocument,
	organizationID string,
) ([]*domain.FHIRIdentifierInput, error) {
	if ids == nil {
		return nil, nil
	}

	identifierSystem := CodeSystem(common.PersonCodeSystemIdentifier)
	output := []*domain.FHIRIdentifierInput{}
	userSelected := true
	version := DefaultVersion

	for _, id := range ids {
		if id.DocumentNumber == "" {
			continue
		}

		identifier := &domain.FHIRIdentifierInput{
			Use: domain.IdentifierUseEnumOfficial,
			Type: &domain.FHIRCodeableConceptInput{
				Coding: []*domain.FHIRCodingInput{
					{
						System:       identifierSystem,
						Version:      &version,
						Code:         scalarutils.Code(id.DocumentType.String()),
						Display:      id.DocumentType.GetDisplayName(),
						UserSelected: &userSelected,
					},
				},
				Text: id.DocumentNumber,
			},
			System: identifierSystem,
			Value:  id.DocumentNumber,
			Period: common.DefaultPeriodInput(),
		}
		identifier.AddAssigner(organizationID)
		output = append(output, identifier)
	}

	return output, nil
}

// NameToHumanName translates the simple name input of simple
// patient registration to FHIR human names
func NameToHumanName(names []*domain.NameInput) []*domain.FHIRHumanNameInput {
	if names == nil {
		return nil
	}

	output := []*domain.FHIRHumanNameInput{}

	for _, name := range names {
		otherNames := ""
		if name.OtherNames != nil {
			otherNames = *name.OtherNames
		}

		fullName := fmt.Sprintf(
			"%s, %s %s", name.LastName, name.FirstName, otherNames)
		fullName = strings.TrimSpace(fullName)
		use := domain.HumanNameUseEnumOfficial
		humanName := &domain.FHIRHumanNameInput{
			Given:  []string{name.FirstName},
			Family: name.LastName,
			Use:    use,
			Period: common.DefaultPeriodInput(),
			Text:   fullName,
		}
		output = append(output, humanName)
	}

	return output
}

// PhysicalPostalAddressesToFHIRAddresses translates address inputs to FHIR addresses
func PhysicalPostalAddressesToFHIRAddresses(
	physical []*domain.PhysicalAddress,
	postal []*domain.PostalAddress,
) []*domain.FHIRAddressInput {
	if physical == nil && postal == nil {
		return nil
	}

	output := []*domain.FHIRAddressInput{}
	addrUse := domain.AddressUseEnumHome
	country := DefaultCountry
	physicalAddrType := domain.AddressTypeEnumPhysical
	postalAddrType := domain.AddressTypeEnumPostal

	for _, postal := range postal {
		text := fmt.Sprintf("%s\n%s", *postal.PostalAddress, postal.PostalCode)
		postalCode := scalarutils.Code(postal.PostalCode)
		postalAddr := &domain.FHIRAddressInput{
			Use:        &addrUse,
			Type:       &postalAddrType,
			Country:    &country,
			Period:     common.DefaultPeriodInput(),
			PostalCode: &postalCode,
			Line:       []*string{postal.PostalAddress},
			Text:       text,
		}
		output = append(output, postalAddr)
	}

	for _, physical := range physical {
		text := fmt.Sprintf(
			"%s\n%s", physical.MapsCode, *physical.PhysicalAddress)
		mapsCode := scalarutils.Code(physical.MapsCode)
		physicalAddr := &domain.FHIRAddressInput{
			Use:        &addrUse,
			Type:       &physicalAddrType,
			Country:    &country,
			Period:     common.DefaultPeriodInput(),
			PostalCode: &mapsCode,
			Line:       []*string{physical.PhysicalAddress},
			Text:       text,
		}
		output = append(output, physicalAddr)
	}

	return output
}

// MaritalStatusEnumToCodeableConcept turns the simple enum selected in the
// user interface to a FHIR codeable concept
func MaritalStatusEnumToCodeableConcept(val domain.MaritalStatus) *domain.FHIRCodeableConcept {
	sel := true
	disp := domain.MaritalStatusDisplay(val)
	codingCode := scalarutils.Code(val.String())
	output := &domain.FHIRCodeableConcept{
		Coding: []*domain.FHIRCoding{
			{
				Code:         &codingCode,
				Display:      disp,
				UserSelected: &sel,
			},
		},
		Text: domain.MaritalStatusDisplay(val),
	}

	return output
}

// LanguagesToCommunicationInputs translates the supplied languages to FHIR
// communication preferences
func LanguagesToCommunicationInputs(languages []enumutils.Language) []*domain.FHIRPatientCommunicationInput {
	output := []*domain.FHIRPatientCommunicationInput{}
	preferred := false
	userSelected := true
	system := scalarutils.URI(enumutils.LanguageCodingSystem)
	version := enumutils.LanguageCodingVersion

	for _, language := range languages {
		comm := &domain.FHIRPatientCommunicationInput{
			Language: &domain.FHIRCodeableConceptInput{
				Coding: []*domain.FHIRCodingInput{
					{
						Code:         scalarutils.Code(language.String()),
						Display:      enumutils.LanguageNames[language],
						UserSelected: &userSelected,
						System:       &system,
						Version:      &version,
					},
				},
				Text: enumutils.LanguageNames[language],
			},
			Preferred: &preferred,
		}
		output = append(output, comm)
	}

	return output
}

// PhysicalPostalAddressesToCombinedFHIRAddress translates address inputs to
// a single FHIR address.
//
// It is used for patient contacts (e.g next of kin) where the spec has only
// one address per next of kin.
func PhysicalPostalAddressesToCombinedFHIRAddress(
	physical []*domain.PhysicalAddress,
	postal []*domain.PostalAddress,
) *domain.FHIRAddressInput {
	if physical == nil && postal == nil {
		return nil
	}

	addressUse := domain.AddressUseEnumHome
	postalAddrType := domain.AddressTypeEnumPostal
	country := DefaultCountry

	addr := &domain.FHIRAddressInput{
		Use:     &addressUse,
		Type:    &postalAddrType,
		Country: &country,
		Period:  common.DefaultPeriodInput(),
		Line:    nil, // will be replaced below
		Text:    "",  // will be replaced below
	}

	postalAddressLines := []string{}
	for _, postal := range postal {
		postalAddressLines = append(postalAddressLines, *postal.PostalAddress)
		postalAddressLines = append(postalAddressLines, postal.PostalCode)

		if addr.PostalCode == nil {
			postalCode := scalarutils.Code(postal.PostalCode)
			addr.PostalCode = &postalCode
		}
	}

	combinedPostalAddress := strings.Join(postalAddressLines, "\n")
	addr.Line = []*string{&combinedPostalAddress}

	physicalAddressLines := []string{}
	for _, physical := range physical {
		physicalAddressLines = append(physicalAddressLines, *physical.PhysicalAddress)
		physicalAddressLines = append(physicalAddressLines, physical.MapsCode)
	}

	combinedPhysicalAddress := strings.Join(physicalAddressLines, "\n")
	addr.Text = combinedPhysicalAddress

	return addr
}

// ContactsToContactPoint translates phone and email contacts to
// FHIR contact points
func ContactsToContactPoint(
	_ context.Context,
	phones []*domain.PhoneNumberInput,
	emails []*domain.EmailInput,
	firestoreClient *firestore.Client,
) ([]*domain.FHIRContactPoint, error) {
	if phones == nil && emails == nil {
		return nil, nil
	}

	output := []*domain.FHIRContactPoint{}
	rank := int64(1)
	contactUse := domain.ContactPointUseEnumHome
	emailSystem := domain.ContactPointSystemEnumEmail
	phoneSystem := domain.ContactPointSystemEnumPhone

	for _, phone := range phones {
		normalized, err := converterandformatter.NormalizeMSISDN(phone.Msisdn)
		if err != nil {
			return nil, fmt.Errorf("unable to normalize phone number: %w", err)
		}

		phoneContact := &domain.FHIRContactPoint{
			System: &phoneSystem,
			Use:    &contactUse,
			Rank:   &rank,
			Period: common.DefaultPeriod(),
			Value:  normalized,
		}
		output = append(output, phoneContact)
		rank++
	}

	for _, email := range emails {
		err := ValidateEmail(
			*email.Email, email.CommunicationOptIn, firestoreClient)
		if err != nil {
			return nil, fmt.Errorf("invalid email: %w", err)
		}

		emailContact := &domain.FHIRContactPoint{
			System: &emailSystem,
			Use:    &contactUse,
			Rank:   &rank,
			Period: common.DefaultPeriod(),
			Value:  email.Email,
		}
		output = append(output, emailContact)
		rank++
	}

	return output, nil
}

// ValidateEmail returns an error if the supplied string does not have a
// valid format or resolvable host
func ValidateEmail(
	email string, _ bool, _ *firestore.Client) error {
	if !govalidator.IsEmail(email) {
		return fmt.Errorf("invalid email format")
	}

	return nil
}

// Validate takes in a struct and validates its fields.
func Validate(object interface{}) error {
	v := validator.New()
	err := v.Struct(object)

	return err
}

// MaritalStatusEnumToCodeableConceptInput turns the simple enum selected in the
// user interface to a FHIR codeable concept
func MaritalStatusEnumToCodeableConceptInput(val domain.MaritalStatus) *domain.FHIRCodeableConceptInput {
	userSelected := true
	output := &domain.FHIRCodeableConceptInput{
		Coding: []*domain.FHIRCodingInput{
			{
				Code:         scalarutils.Code(val.String()),
				Display:      domain.MaritalStatusDisplay(val),
				UserSelected: &userSelected,
			},
		},
		Text: domain.MaritalStatusDisplay(val),
	}

	return output
}

func CodeSystem(value string) *scalarutils.URI {
	path := fmt.Sprintf("%s/CodeSystem/%s", baseURL, value)
	return (*scalarutils.URI)(&path)
}

// ExtractTextFromHTML is a utility function that extract the text from scalarutils.XHTML.
// The essence of this is that the text property predominantly used to show "UsageContext" that helpers
// the fronted control workflows within Empower platform
func ExtractTextFromHTML(htmlStr string) string {
	doc, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		return ""
	}

	return extractText(doc)
}

// Helper function to traverse the parsed HTML tree
func extractText(n *html.Node) string {
	if n.Type == html.TextNode {
		return strings.TrimSpace(n.Data)
	}

	var textContent string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		textContent += extractText(c) + " "
	}

	return strings.TrimSpace(textContent)
}

// SetTokenURL sets the tokenUrl based on environment.
func SetTokenURL(docTemplate, currentTokenURL, newTokenURL string) string {
	return strings.ReplaceAll(docTemplate, currentTokenURL, newTokenURL)
}

// CalculateScoresFromResponses Calculates scores for each group in the questionnaire based on responses.
// This function iterates through each item in the QuestionnaireResponse, matching it with its corresponding item in the Questionnaire
// to apply scoring rules based on 'ordinalValue'. This matching is essential to accurately compute scores since it directly ties
// the respondent's answers back to the structured definitions within the Questionnaire, including the consideration of 'ordinalValue'
// for scoring purposes.
func CalculateScoresFromResponses(questionnaire *domain.FHIRQuestionnaire, response *dto.QuestionnaireResponse) map[string]int {
	groupScores := make(map[string]int) // Holds the calculated scores for each group.

	ordinalValues := mapLinkIdsToOrdinalValuesAndDisplay(questionnaire.Item)

	for _, responseGroup := range response.Item {
		groupScore := 0
		groupLinkID := responseGroup.LinkID

		for _, questionnaireGroup := range questionnaire.Item {
			if *questionnaireGroup.LinkID == groupLinkID {
				// Calculate score for this group based on responses
				groupScore = calculateGroupScore(responseGroup.Item, ordinalValues)
				break
			}
		}

		groupScores[groupLinkID] = groupScore
	}

	return groupScores
}

func calculateGroupScore(responseItems []dto.QuestionnaireResponseItem, ordinalValues map[string]int) int {
	score := 0

	for _, responseItem := range responseItems {
		if val, ok := ordinalValues[responseItem.LinkID]; ok {
			for _, answer := range responseItem.Answer {
				if answer.ValueCoding.Display == "Yes" {
					score += val
					break
				}
			}
		}
	}

	return score
}

// mapLinkIdsToOrdinalValuesAndDisplay Maps each linkId in the Questionnaire to its corresponding 'ordinalValue', if available.
// This mapping facilitates a more efficient scoring process by eliminating the need to repeatedly search through the
// questionnaire structure for 'ordinalValue' extensions during scoring.
func mapLinkIdsToOrdinalValuesAndDisplay(items []*domain.FHIRQuestionnaireItem) map[string]int {
	ordinalValues := make(map[string]int)

	for _, item := range items {
		for _, questionItem := range item.Item {
			// Loop through answerOption to find "Yes" answers with ordinalValue
			for _, answerOption := range questionItem.AnswerOption {
				if answerOption.ValueCoding.Display == "Yes" {
					for _, ext := range answerOption.Extension {
						if ext.URL == "http://hl7.org/fhir/StructureDefinition/ordinalValue" {
							ordinalValues[*questionItem.LinkID] = int(*ext.ValueDecimal)
							break
						}
					}
				}
			}
		}
	}

	return ordinalValues
}

func StripID(reference *string) string {
	if reference == nil {
		return ""
	}

	if idx := strings.LastIndex(*reference, "/"); idx != -1 {
		return (*reference)[idx+1:]
	}

	return *reference
}

func StripIDFromReference(ref *domain.FHIRReference) string {
	if ref == nil {
		return ""
	}

	if ref.ID != nil {
		return *ref.ID
	}

	return StripID(ref.Reference)
}
