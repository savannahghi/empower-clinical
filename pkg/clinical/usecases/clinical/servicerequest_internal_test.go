package clinical

import (
	"testing"

	"github.com/savannahghi/scalarutils"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/common"
	"github.com/savannahghi/empower-clinical/pkg/clinical/domain"
)

// coding builds a *domain.FHIRCoding from a code and display for use in table tests.
func coding(code, display string) *domain.FHIRCoding {
	c := scalarutils.Code(code)

	return &domain.FHIRCoding{Code: &c, Display: display}
}

// serviceRequestWithCodings wraps the supplied codings in a service request's Code block.
func serviceRequestWithCodings(codings ...*domain.FHIRCoding) domain.FHIRServiceRequest {
	return domain.FHIRServiceRequest{
		Code: &domain.FHIRCodeableReference{
			Concept: &domain.FHIRCodeableConcept{
				Coding: codings,
			},
		},
	}
}

func Test_hasOnlyReferralNarrativeCoding(t *testing.T) {
	tests := []struct {
		name           string
		serviceRequest domain.FHIRServiceRequest
		want           bool
	}{
		{
			name:           "only the referral-narrative coding is present",
			serviceRequest: serviceRequestWithCodings(coding(common.ReferralReasonLOINCode, "Reason for referral (narrative)")),
			want:           true,
		},
		{
			name: "multiple codings but all are referral-narrative",
			serviceRequest: serviceRequestWithCodings(
				coding(common.ReferralReasonLOINCode, "Reason for referral (narrative)"),
				coding(common.ReferralReasonLOINCode, "Reason for referral (narrative)"),
			),
			want: true,
		},
		{
			name: "nil coding entries around a lone narrative coding",
			serviceRequest: serviceRequestWithCodings(
				nil,
				coding(common.ReferralReasonLOINCode, "Reason for referral (narrative)"),
				&domain.FHIRCoding{Display: "no code"},
			),
			want: true,
		},
		{
			name: "referral-narrative plus an ordered test coding is a real lab order",
			serviceRequest: serviceRequestWithCodings(
				coding(common.ReferralReasonLOINCode, "Reason for referral (narrative)"),
				coding(common.TestsOrderedLOINCCode, "IHC_HER2"),
			),
			want: false,
		},
		{
			name: "referral-narrative plus a resolved LOINC test coding is a real lab order",
			serviceRequest: serviceRequestWithCodings(
				coding(common.ReferralReasonLOINCode, "Reason for referral (narrative)"),
				coding(common.IHCEstrogenReceptorLOINCCode, "IHC Estrogen Receptor"),
			),
			want: false,
		},
		{
			name:           "self/intra lab order with a real test coding only",
			serviceRequest: serviceRequestWithCodings(coding(common.HER2LOINCCode, "HER2 Ag [Presence] in Breast cancer specimen by Immune stain")),
			want:           false,
		},
		{
			name:           "empty coding array",
			serviceRequest: serviceRequestWithCodings(),
			want:           false,
		},
		{
			name:           "nil Code block",
			serviceRequest: domain.FHIRServiceRequest{},
			want:           false,
		},
		{
			name: "Code present but Concept is nil",
			serviceRequest: domain.FHIRServiceRequest{
				Code: &domain.FHIRCodeableReference{},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasOnlyReferralNarrativeCoding(tt.serviceRequest); got != tt.want {
				t.Errorf("hasOnlyReferralNarrativeCoding() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getOrderedTestCoding(t *testing.T) {
	tests := []struct {
		name           string
		serviceRequest domain.FHIRServiceRequest
		wantCode       string
		wantName       string
		wantOK         bool
	}{
		{
			name: "referral order under the tests-ordered marker with the raw enum in Display",
			serviceRequest: serviceRequestWithCodings(
				coding(common.ReferralReasonLOINCode, "Reason for referral (narrative)"),
				coding(common.TestsOrderedLOINCCode, "IHC_HER2"),
			),
			wantCode: common.HER2LOINCCode,
			wantName: "IHC HER2",
			wantOK:   true,
		},
		{
			name: "referral order under the tests-ordered marker with the humanised label in Display",
			serviceRequest: serviceRequestWithCodings(
				coding(common.ReferralReasonLOINCode, "Reason for referral (narrative)"),
				coding(common.TestsOrderedLOINCCode, "IHC HER2"),
			),
			wantCode: common.HER2LOINCCode,
			wantName: "IHC HER2",
			wantOK:   true,
		},
		{
			name: "referral order storing the resolved LOINC code alongside the narrative coding",
			serviceRequest: serviceRequestWithCodings(
				coding(common.ReferralReasonLOINCode, "Reason for referral (narrative)"),
				coding(common.IHCEstrogenReceptorLOINCCode, "IHC Estrogen Receptor"),
			),
			wantCode: common.IHCEstrogenReceptorLOINCCode,
			wantName: "IHC Estrogen Receptor",
			wantOK:   true,
		},
		{
			name:           "self/intra lab order carries the real code and display directly",
			serviceRequest: serviceRequestWithCodings(coding(common.HER2LOINCCode, "HER2 Ag [Presence] in Breast cancer specimen by Immune stain")),
			wantCode:       common.HER2LOINCCode,
			wantName:       "HER2 Ag [Presence] in Breast cancer specimen by Immune stain",
			wantOK:         true,
		},
		{
			name:           "narrative-only referral yields no test coding",
			serviceRequest: serviceRequestWithCodings(coding(common.ReferralReasonLOINCode, "Reason for referral (narrative)")),
			wantCode:       "",
			wantName:       "",
			wantOK:         false,
		},
		{
			name:           "nil Code block yields no coding",
			serviceRequest: domain.FHIRServiceRequest{},
			wantCode:       "",
			wantName:       "",
			wantOK:         false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCode, gotName, gotOK := getOrderedTestCoding(tt.serviceRequest)
			if gotCode != tt.wantCode || gotName != tt.wantName || gotOK != tt.wantOK {
				t.Errorf("getOrderedTestCoding() = (%q, %q, %v), want (%q, %q, %v)",
					gotCode, gotName, gotOK, tt.wantCode, tt.wantName, tt.wantOK)
			}
		})
	}
}
