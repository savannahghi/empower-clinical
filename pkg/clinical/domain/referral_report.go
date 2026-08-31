package domain

// The models below will be used to support the generation of a patient's referral report PDF file from Golang templates.

type Patient struct {
	Name        string `json:"name,omitempty"`
	EmpowerID   string `json:"empowerID,omitempty"`
	NationalID  string `json:"nationalID,omitempty"`
	PhoneNumber string `json:"phoneNumber,omitempty"`
	DateOfBirth string `json:"dateOfBirth,omitempty"`
	Age         int    `json:"age,omitempty"`
	Sex         string `json:"sex,omitempty"`
}

type NextOfKin struct {
	Name         string `json:"name,omitempty"`
	Relationship string `json:"relationship,omitempty"`
	PhoneNumber  string `json:"phoneNumber,omitempty"`
}

// ReferralReason represents the reason for the referral, including specific tests, the date, and any additional notes.
type ReferralReason struct {
	Reason string `json:"reason,omitempty"`
	Test   string `json:"test,omitempty"`
	Date   string `json:"date,omitempty"`
	Note   string `json:"note,omitempty"`
}

type Test struct {
	Name    string `json:"name,omitempty"`
	Results string `json:"results,omitempty"`
	Date    string `json:"date,omitempty"`
}

type MedicalHistory struct {
	Procedure  string `json:"procedure,omitempty"`
	Medication string `json:"medication,omitempty"`
	Tests      []Test `json:"tests,omitempty"`
}

type Footer struct {
	Phone   string `json:"phone,omitempty"`
	Email   string `json:"email,omitempty"`
	Address string `json:"address,omitempty"`
}

// ReferralReport models the details that are required on a patient's referral form
type ReferralReport struct {
	ReceivingFacility ReceivingFacility `json:"receivingFacility,omitempty"`
	Patient           Patient           `json:"patient,omitempty"`
	NextOfKin         NextOfKin         `json:"nextOfKin,omitempty"`
	Facility          string            `json:"facility,omitempty"`
	ReferralReason    ReferralReason    `json:"referralReason,omitempty"`
	MedicalHistory    MedicalHistory    `json:"medicalHistory,omitempty"`
	Footer            Footer            `json:"footer,omitempty"`
}

// PatientReferralDetails models the patient referral details. It is mainly what is found on a referral report
type PatientReferralDetails struct {
	PatientName        string            `json:"patientName,omitempty"`
	ReceivingFacility  ReceivingFacility `json:"receivingFacility,omitempty"`
	ReferralReason     ReferralReason    `json:"referralReason,omitempty"`
	MedicalHistory     MedicalHistory    `json:"medicalHistory,omitempty"`
	ReferralReportLink string            `json:"referralReportLink,omitempty"`
}
