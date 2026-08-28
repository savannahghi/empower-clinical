package dto

// ReferralDetail represents the structured information of a patient's referral
type ReferralDetail struct {
	ID           string `json:"id,omitempty"`
	ReferralDate string `json:"referralDate,omitempty"`
	// This could be Diagnostics, Specialist etc..
	ReferredFor string `json:"referredFor,omitempty"`

	// The actual tests the patient was referred for e.g Mammogram, HPV etc..
	ReferredTests      []string          `json:"referredTests,omitempty"`
	ReferredTo         string            `json:"referredTo,omitempty"`
	ReferralReportLink *string           `json:"referralReportLink,omitempty"`
	PatientName        string            `json:"patientName,omitempty"`
	PatientID          string            `json:"patientID,omitempty"`
	UsageContext       ScreeningTypeEnum `json:"usageContext,omitempty"`
	Status             string            `json:"status"`
}

// ReferralDetailEdge is an referral edge
type ReferralDetailEdge struct {
	Node   ReferralDetail `json:"node,omitempty"`
	Cursor string         `json:"cursor,omitempty"`
}

// ReferralDetailConnection  is an ReferralDetailConnection Connection Type
type ReferralDetailConnection struct {
	TotalCount int                  `json:"totalCount,omitempty"`
	Edges      []ReferralDetailEdge `json:"edges,omitempty"`
	PageInfo   PageInfo             `json:"pageInfo,omitempty"`
}

func CreateReferralDetailConnection(referrals []*ReferralDetail, pageInfo PageInfo, total int) ReferralDetailConnection {
	connection := ReferralDetailConnection{
		TotalCount: total,
		Edges:      []ReferralDetailEdge{},
		PageInfo:   pageInfo,
	}

	for _, referral := range referrals {
		edge := ReferralDetailEdge{
			Node: *referral,
		}

		connection.Edges = append(connection.Edges, edge)
	}

	return connection
}
