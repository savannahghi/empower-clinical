package domain

type PublicationStatus string

const (
	PublicationStatusDraft   PublicationStatus = "draft"
	PublicationStatusActive  PublicationStatus = "active"
	PublicationStatusRetired PublicationStatus = "retired"
	PublicationStatusUnknown PublicationStatus = "unknown"
)

func (code PublicationStatus) Code() string {
	switch code {
	case PublicationStatusDraft:
		return "draft"
	case PublicationStatusActive:
		return "active"
	case PublicationStatusRetired:
		return "retired"
	case PublicationStatusUnknown:
		return "unknown"
	}

	return "<unknown>"
}

func (code PublicationStatus) Display() string {
	switch code {
	case PublicationStatusDraft:
		return "Draft"
	case PublicationStatusActive:
		return "Active"
	case PublicationStatusRetired:
		return "Retired"
	case PublicationStatusUnknown:
		return "Unknown"
	}

	return "<unknown>"
}

func (code PublicationStatus) Definition() string {
	switch code {
	case PublicationStatusDraft:
		return "This resource is still under development and is not yet considered to be ready for normal use."
	case PublicationStatusActive:
		return "This resource is ready for normal use."
	case PublicationStatusRetired:
		return "This resource has been withdrawn or superseded and should no longer be used."
	case PublicationStatusUnknown:
		return "The authoring system does not know which of the status values currently applies for this resource.  Note: This concept is not to be used for \"other\" - one of the listed statuses is presumed to apply, it's just not known which one."
	}

	return "<unknown>"
}
