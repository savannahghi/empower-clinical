package clinical_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/empower-clinical/pkg/clinical/domain"
	usecaseMock "github.com/savannahghi/empower-clinical/pkg/clinical/usecases/mock"
)

func TestUseCasesClinicalImpl_CreateQuestionnaireResponse(t *testing.T) {
	ctx := context.Background()
	ID := gofakeit.UUID()

	tenantIDs := &dto.TenantIdentifiers{
		OrganizationID: uuid.NewString(),
		FacilityID:     uuid.NewString(),
	}

	encounter := &domain.FHIREncounterRelayPayload{
		Resource: &domain.FHIREncounter{
			ID:     &ID,
			Status: domain.EncounterStatusEnumInProgress,
			Subject: &domain.FHIRReference{
				ID:      &ID,
				Display: ID,
			},
		},
	}

	orgName := mock.Anything
	organization := &domain.FHIROrganizationRelayPayload{
		Resource: &domain.FHIROrganization{
			ID:   &ID,
			Name: &orgName,
		},
	}

	questionnaireName := "Cervical Cancer Screening"
	questionnaireRelayPayload := &domain.FHIRQuestionnaireRelayPayload{
		Resource: &domain.FHIRQuestionnaire{
			ID:    &ID,
			Name:  &questionnaireName,
			Title: &questionnaireName,
		},
	}

	patient := &domain.FHIRPatientRelayPayload{
		Resource: &domain.FHIRPatient{
			ID: &ID,
		},
	}

	riskassessment := &domain.FHIRRiskAssessmentRelayPayload{
		Resource: &domain.FHIRRiskAssessment{
			ID: &ID,
			Prediction: []domain.FHIRRiskAssessmentPrediction{
				{
					QualitativeRisk: &domain.FHIRCodeableConcept{
						Text: gofakeit.BeerName(),
					},
				},
			},
		},
	}

	type args struct {
		ctx             context.Context
		input           dto.QuestionnaireResponse
		questionnaireID string
		encounterID     string
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Sad case: unable to get tenant tags",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, input: dto.QuestionnaireResponse{}, questionnaireID: ID, encounterID: ID}
			},
			wantErr: true,
		},

		{
			name: "Sad Case: invalid encounter id",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, input: dto.QuestionnaireResponse{}, questionnaireID: ID, encounterID: ""}
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to create questionnaire response",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					})
				mh.FHIR.EXPECT().CreateFHIRQuestionnaireResponse(mock.Anything, mock.Anything).
					RunAndReturn(func(context1 context.Context, input *domain.FHIRQuestionnaireResponse) (*domain.FHIRQuestionnaireResponse, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, input: dto.QuestionnaireResponse{}, questionnaireID: ID, encounterID: ID}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Attempt to record questionnaire response in a finished encounter",
			setup: func(mh *usecaseMock.Mocks) args {
				encounter := &domain.FHIREncounterRelayPayload{
					Resource: &domain.FHIREncounter{
						ID:     &ID,
						Status: domain.EncounterStatusEnumCompleted,
					},
				}

				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})

				return args{ctx: ctx, input: dto.QuestionnaireResponse{}, questionnaireID: ID, encounterID: ID}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get fhir questionnaire",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					})
				mh.FHIR.EXPECT().CreateFHIRQuestionnaireResponse(mock.Anything, mock.Anything).
					RunAndReturn(func(context1 context.Context, input *domain.FHIRQuestionnaireResponse) (*domain.FHIRQuestionnaireResponse, error) {
						score := 3
						questionnaireReponse := &domain.FHIRQuestionnaireResponse{
							ID: &ID,
							Item: []domain.FHIRQuestionnaireResponseItem{
								{
									LinkID: "symptoms",
									Item: []domain.FHIRQuestionnaireResponseItem{
										{
											LinkID: "symptoms-score",
											Answer: []domain.FHIRQuestionnaireResponseItemAnswer{
												{
													ValueInteger: &score,
												},
											},
										},
									},
								},
								{
									LinkID: "risk-factors",
									Item: []domain.FHIRQuestionnaireResponseItem{
										{
											LinkID: "risk-factors-score",
											Answer: []domain.FHIRQuestionnaireResponseItemAnswer{
												{
													ValueInteger: &score,
												},
											},
										},
									},
								},
							},
						}
						return questionnaireReponse, nil
					})
				mh.FHIR.EXPECT().GetFHIRQuestionnaire(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRQuestionnaireRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, input: dto.QuestionnaireResponse{}, questionnaireID: ID, encounterID: ID}
			},
			wantErr: true,
		},
		{
			name: "Happy Case - Create questionnaire response and generate review summary - Cervical Cancer - High Risk",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					})
				mh.FHIR.EXPECT().CreateFHIRQuestionnaireResponse(mock.Anything, mock.Anything).
					RunAndReturn(func(context1 context.Context, input *domain.FHIRQuestionnaireResponse) (*domain.FHIRQuestionnaireResponse, error) {
						score := 3
						questionnaireReponse := &domain.FHIRQuestionnaireResponse{
							ID: &ID,
							Item: []domain.FHIRQuestionnaireResponseItem{
								{
									LinkID: "symptoms",
									Item: []domain.FHIRQuestionnaireResponseItem{
										{
											LinkID: "symptoms-score",
											Answer: []domain.FHIRQuestionnaireResponseItemAnswer{
												{
													ValueCoding: &domain.FHIRCoding{
														Display: "Yes",
													},
													ValueInteger: &score,
												},
											},
										},
									},
								},
								{
									LinkID: "risk-factors",
									Item: []domain.FHIRQuestionnaireResponseItem{
										{
											LinkID: "risk-factors-score",
											Answer: []domain.FHIRQuestionnaireResponseItemAnswer{
												{
													ValueCoding: &domain.FHIRCoding{
														Display: "Yes",
													},
													ValueInteger: &score,
												},
											},
										},
									},
								},
							},
						}

						return questionnaireReponse, nil
					})
				mh.FHIR.EXPECT().GetFHIRQuestionnaire(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRQuestionnaireRelayPayload, error) {
						questionnaireName := "Cervical Cancer Screening"
						symptoms := "symptoms"
						symptomsScore := "symptoms-score"
						riskFactors := "risk-factors"
						riskFactorsScore := "risk-factors-score"
						score := 3.0
						questionnaireRelayPayload := &domain.FHIRQuestionnaireRelayPayload{
							Resource: &domain.FHIRQuestionnaire{
								ID:    &ID,
								Name:  &questionnaireName,
								Title: &questionnaireName,
								Item: []*domain.FHIRQuestionnaireItem{
									{
										LinkID: &symptoms,
										Item: []*domain.FHIRQuestionnaireItem{
											{
												LinkID: &symptomsScore,
												AnswerOption: []*domain.FHIRQuestionnaireItemAnswerOption{
													{
														ValueCoding: &domain.FHIRCoding{
															Display: "Yes",
														},
														Extension: []*domain.Extension{
															{
																URL:          "http://hl7.org/fhir/StructureDefinition/ordinalValue",
																ValueDecimal: &score,
															},
														},
													},
												},
											},
										},
									},
									{
										LinkID: &riskFactors,
										Item: []*domain.FHIRQuestionnaireItem{
											{
												LinkID: &riskFactorsScore,
												AnswerOption: []*domain.FHIRQuestionnaireItemAnswerOption{
													{
														ValueCoding: &domain.FHIRCoding{
															Display: "Yes",
														},
														Extension: []*domain.Extension{
															{
																URL:          "http://hl7.org/fhir/StructureDefinition/ordinalValue",
																ValueDecimal: &score,
															},
														},
													},
												},
											},
										},
									},
								},
							},
						}
						return questionnaireRelayPayload, nil
					})
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return patient, nil
					})
				mh.PubSub.EXPECT().NotifySegmentation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, data dto.SegmentationPayload) error {
						return nil
					})
				mh.FHIR.EXPECT().CreateFHIRRiskAssessment(mock.Anything, mock.Anything).
					RunAndReturn(func(context1 context.Context, input *domain.FHIRRiskAssessmentInput) (*domain.FHIRRiskAssessmentRelayPayload, error) {
						return riskassessment, nil
					})

				return args{ctx: ctx, input: dto.QuestionnaireResponse{}, questionnaireID: ID, encounterID: ID}
			},
			wantErr: false,
		},
		{
			name: "Happy Case - Create questionnaire response and generate review summary - Cervical Cancer - Low Risk",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					})
				mh.FHIR.EXPECT().CreateFHIRQuestionnaireResponse(mock.Anything, mock.Anything).
					RunAndReturn(func(context1 context.Context, input *domain.FHIRQuestionnaireResponse) (*domain.FHIRQuestionnaireResponse, error) {
						score := 0
						questionnaireReponse := &domain.FHIRQuestionnaireResponse{
							ID: &ID,
							Item: []domain.FHIRQuestionnaireResponseItem{
								{
									LinkID: "symptoms",
									Item: []domain.FHIRQuestionnaireResponseItem{
										{
											LinkID: "symptoms-score",
											Answer: []domain.FHIRQuestionnaireResponseItemAnswer{
												{
													ValueCoding: &domain.FHIRCoding{
														Display: "Yes",
													},
													ValueInteger: &score,
												},
											},
										},
									},
								},
								{
									LinkID: "risk-factors",
									Item: []domain.FHIRQuestionnaireResponseItem{
										{
											LinkID: "risk-factors-score",
											Answer: []domain.FHIRQuestionnaireResponseItemAnswer{
												{
													ValueCoding: &domain.FHIRCoding{
														Display: "Yes",
													},
													ValueInteger: &score,
												},
											},
										},
									},
								},
							},
						}

						return questionnaireReponse, nil
					})
				mh.FHIR.EXPECT().GetFHIRQuestionnaire(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRQuestionnaireRelayPayload, error) {
						questionnaireName := "Cervical Cancer Screening"
						symptoms := "symptoms"
						symptomsScore := "symptoms-score"
						riskFactors := "risk-factors"
						riskFactorsScore := "risk-factors-score"
						score := 0.0
						questionnaireRelayPayload := &domain.FHIRQuestionnaireRelayPayload{
							Resource: &domain.FHIRQuestionnaire{
								ID:    &ID,
								Name:  &questionnaireName,
								Title: &questionnaireName,
								Item: []*domain.FHIRQuestionnaireItem{
									{
										LinkID: &symptoms,
										Item: []*domain.FHIRQuestionnaireItem{
											{
												LinkID: &symptomsScore,
												AnswerOption: []*domain.FHIRQuestionnaireItemAnswerOption{
													{
														ValueCoding: &domain.FHIRCoding{
															Display: "Yes",
														},
														Extension: []*domain.Extension{
															{
																URL:          "http://hl7.org/fhir/StructureDefinition/ordinalValue",
																ValueDecimal: &score,
															},
														},
													},
												},
											},
										},
									},
									{
										LinkID: &riskFactors,
										Item: []*domain.FHIRQuestionnaireItem{
											{
												LinkID: &riskFactorsScore,
												AnswerOption: []*domain.FHIRQuestionnaireItemAnswerOption{
													{
														ValueCoding: &domain.FHIRCoding{
															Display: "Yes",
														},
														Extension: []*domain.Extension{
															{
																URL:          "http://hl7.org/fhir/StructureDefinition/ordinalValue",
																ValueDecimal: &score,
															},
														},
													},
												},
											},
										},
									},
								},
							},
						}
						return questionnaireRelayPayload, nil
					})
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return patient, nil
					})
				mh.PubSub.EXPECT().NotifySegmentation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, data dto.SegmentationPayload) error {
						return nil
					})
				mh.FHIR.EXPECT().CreateFHIRRiskAssessment(mock.Anything, mock.Anything).
					RunAndReturn(func(context1 context.Context, input *domain.FHIRRiskAssessmentInput) (*domain.FHIRRiskAssessmentRelayPayload, error) {
						return riskassessment, nil
					})

				return args{ctx: ctx, input: dto.QuestionnaireResponse{}, questionnaireID: ID, encounterID: ID}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to record risk assessment - Low Risk",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					})
				mh.FHIR.EXPECT().CreateFHIRQuestionnaireResponse(mock.Anything, mock.Anything).
					RunAndReturn(func(context1 context.Context, input *domain.FHIRQuestionnaireResponse) (*domain.FHIRQuestionnaireResponse, error) {
						score := 1
						questionnaireReponse := &domain.FHIRQuestionnaireResponse{
							ID: &ID,
							Item: []domain.FHIRQuestionnaireResponseItem{
								{
									LinkID: "symptoms",
									Item: []domain.FHIRQuestionnaireResponseItem{
										{
											LinkID: "symptoms-score",
											Answer: []domain.FHIRQuestionnaireResponseItemAnswer{
												{
													ValueCoding: &domain.FHIRCoding{
														Display: "Yes",
													},
													ValueInteger: &score,
												},
											},
										},
									},
								},
								{
									LinkID: "risk-factors",
									Item: []domain.FHIRQuestionnaireResponseItem{
										{
											LinkID: "risk-factors-score",
											Answer: []domain.FHIRQuestionnaireResponseItemAnswer{
												{
													ValueCoding: &domain.FHIRCoding{
														Display: "Yes",
													},
													ValueInteger: &score,
												},
											},
										},
									},
								},
							},
						}

						return questionnaireReponse, nil
					})
				mh.FHIR.EXPECT().GetFHIRQuestionnaire(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRQuestionnaireRelayPayload, error) {
						questionnaireName := "Cervical Cancer Screening"
						symptoms := "symptoms"
						symptomsScore := "symptoms-score"
						riskFactors := "risk-factors"
						riskFactorsScore := "risk-factors-score"
						score := 1.0
						questionnaireRelayPayload := &domain.FHIRQuestionnaireRelayPayload{
							Resource: &domain.FHIRQuestionnaire{
								ID:    &ID,
								Name:  &questionnaireName,
								Title: &questionnaireName,
								Item: []*domain.FHIRQuestionnaireItem{
									{
										LinkID: &symptoms,
										Item: []*domain.FHIRQuestionnaireItem{
											{
												LinkID: &symptomsScore,
												AnswerOption: []*domain.FHIRQuestionnaireItemAnswerOption{
													{
														ValueCoding: &domain.FHIRCoding{
															Display: "Yes",
														},
														Extension: []*domain.Extension{
															{
																URL:          "http://hl7.org/fhir/StructureDefinition/ordinalValue",
																ValueDecimal: &score,
															},
														},
													},
												},
											},
										},
									},
									{
										LinkID: &riskFactors,
										Item: []*domain.FHIRQuestionnaireItem{
											{
												LinkID: &riskFactorsScore,
												AnswerOption: []*domain.FHIRQuestionnaireItemAnswerOption{
													{
														ValueCoding: &domain.FHIRCoding{
															Display: "Yes",
														},
														Extension: []*domain.Extension{
															{
																URL:          "http://hl7.org/fhir/StructureDefinition/ordinalValue",
																ValueDecimal: &score,
															},
														},
													},
												},
											},
										},
									},
								},
							},
						}
						return questionnaireRelayPayload, nil
					})
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return patient, nil
					})
				mh.FHIR.EXPECT().CreateFHIRRiskAssessment(mock.Anything, mock.Anything).
					RunAndReturn(func(context1 context.Context, input *domain.FHIRRiskAssessmentInput) (*domain.FHIRRiskAssessmentRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, input: dto.QuestionnaireResponse{}, questionnaireID: ID, encounterID: ID}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to record risk assessment - High Risk",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					})
				mh.FHIR.EXPECT().CreateFHIRQuestionnaireResponse(mock.Anything, mock.Anything).
					RunAndReturn(func(context1 context.Context, input *domain.FHIRQuestionnaireResponse) (*domain.FHIRQuestionnaireResponse, error) {
						score := 3
						questionnaireReponse := &domain.FHIRQuestionnaireResponse{
							ID: &ID,
							Item: []domain.FHIRQuestionnaireResponseItem{
								{
									LinkID: "symptoms",
									Item: []domain.FHIRQuestionnaireResponseItem{
										{
											LinkID: "symptoms-score",
											Answer: []domain.FHIRQuestionnaireResponseItemAnswer{
												{
													ValueCoding: &domain.FHIRCoding{
														Display: "Yes",
													},
													ValueInteger: &score,
												},
											},
										},
									},
								},
								{
									LinkID: "risk-factors",
									Item: []domain.FHIRQuestionnaireResponseItem{
										{
											LinkID: "risk-factors-score",
											Answer: []domain.FHIRQuestionnaireResponseItemAnswer{
												{
													ValueCoding: &domain.FHIRCoding{
														Display: "Yes",
													},
													ValueInteger: &score,
												},
											},
										},
									},
								},
							},
						}

						return questionnaireReponse, nil
					})
				mh.FHIR.EXPECT().GetFHIRQuestionnaire(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRQuestionnaireRelayPayload, error) {
						questionnaireName := "Cervical Cancer Screening"
						symptoms := "symptoms"
						symptomsScore := "symptoms-score"
						riskFactors := "risk-factors"
						riskFactorsScore := "risk-factors-score"
						score := 3.0
						questionnaireRelayPayload := &domain.FHIRQuestionnaireRelayPayload{
							Resource: &domain.FHIRQuestionnaire{
								ID:    &ID,
								Name:  &questionnaireName,
								Title: &questionnaireName,
								Item: []*domain.FHIRQuestionnaireItem{
									{
										LinkID: &symptoms,
										Item: []*domain.FHIRQuestionnaireItem{
											{
												LinkID: &symptomsScore,
												AnswerOption: []*domain.FHIRQuestionnaireItemAnswerOption{
													{
														ValueCoding: &domain.FHIRCoding{
															Display: "Yes",
														},
														Extension: []*domain.Extension{
															{
																URL:          "http://hl7.org/fhir/StructureDefinition/ordinalValue",
																ValueDecimal: &score,
															},
														},
													},
												},
											},
										},
									},
									{
										LinkID: &riskFactors,
										Item: []*domain.FHIRQuestionnaireItem{
											{
												LinkID: &riskFactorsScore,
												AnswerOption: []*domain.FHIRQuestionnaireItemAnswerOption{
													{
														ValueCoding: &domain.FHIRCoding{
															Display: "Yes",
														},
														Extension: []*domain.Extension{
															{
																URL:          "http://hl7.org/fhir/StructureDefinition/ordinalValue",
																ValueDecimal: &score,
															},
														},
													},
												},
											},
										},
									},
								},
							},
						}
						return questionnaireRelayPayload, nil
					})
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return patient, nil
					})
				mh.FHIR.EXPECT().CreateFHIRRiskAssessment(mock.Anything, mock.Anything).
					RunAndReturn(func(context1 context.Context, input *domain.FHIRRiskAssessmentInput) (*domain.FHIRRiskAssessmentRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, input: dto.QuestionnaireResponse{}, questionnaireID: ID, encounterID: ID}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - non-existent fhir questionnaire",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					})
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return patient, nil
					})
				mh.FHIR.EXPECT().CreateFHIRQuestionnaireResponse(mock.Anything, mock.Anything).
					RunAndReturn(func(context1 context.Context, input *domain.FHIRQuestionnaireResponse) (*domain.FHIRQuestionnaireResponse, error) {
						score := 1
						questionnaireReponse := &domain.FHIRQuestionnaireResponse{
							ID: &ID,
							Item: []domain.FHIRQuestionnaireResponseItem{
								{
									LinkID: "symptoms",
									Item: []domain.FHIRQuestionnaireResponseItem{
										{
											LinkID: "symptoms-score",
											Answer: []domain.FHIRQuestionnaireResponseItemAnswer{
												{
													ValueInteger: &score,
												},
											},
										},
									},
								},
								{
									LinkID: "risk-factors",
									Item: []domain.FHIRQuestionnaireResponseItem{
										{
											LinkID: "risk-factors-score",
											Answer: []domain.FHIRQuestionnaireResponseItemAnswer{
												{
													ValueInteger: &score,
												},
											},
										},
									},
								},
							},
						}

						return questionnaireReponse, nil
					})
				mh.FHIR.EXPECT().GetFHIRQuestionnaire(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRQuestionnaireRelayPayload, error) {
						randomName := gofakeit.BeerName()
						linkID := "risk-assessment"
						questionLinkID := "risk-assessment-question-one"
						valueDecimalOne := 1.0
						questionnaireRelayPayload := &domain.FHIRQuestionnaireRelayPayload{
							Resource: &domain.FHIRQuestionnaire{
								ID:    &ID,
								Name:  &randomName,
								Title: &randomName,
								Item: []*domain.FHIRQuestionnaireItem{
									{
										LinkID: &linkID,
										Item: []*domain.FHIRQuestionnaireItem{
											{
												ID:     &questionLinkID,
												LinkID: &questionLinkID,
												Meta:   &domain.FHIRMeta{},
												AnswerOption: []*domain.FHIRQuestionnaireItemAnswerOption{
													{
														Extension: []*domain.Extension{
															{
																URL:          "http://hl7.org/fhir/StructureDefinition/ordinalValue",
																ValueDecimal: &valueDecimalOne,
															},
														},
														ValueCoding: &domain.FHIRCoding{
															Display: "Yes",
														},
													},
												},
											},
										},
									},
								},
							},
						}
						return questionnaireRelayPayload, nil
					})

				return args{ctx: ctx, input: dto.QuestionnaireResponse{}, questionnaireID: ID, encounterID: ID}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - fail to get tenant meta tags",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, input: dto.QuestionnaireResponse{}, questionnaireID: ID, encounterID: ID}
			},
			wantErr: true,
		},
		{
			name: "Happy Case - Create questionnaire response and generate review summary - Breast Cancer - High Risk",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					})
				mh.FHIR.EXPECT().CreateFHIRQuestionnaireResponse(mock.Anything, mock.Anything).
					RunAndReturn(func(context1 context.Context, input *domain.FHIRQuestionnaireResponse) (*domain.FHIRQuestionnaireResponse, error) {
						score := 1
						questionnaireReponse := &domain.FHIRQuestionnaireResponse{
							ID: &ID,
							Item: []domain.FHIRQuestionnaireResponseItem{
								{
									LinkID: "risk-assessment",
									Item: []domain.FHIRQuestionnaireResponseItem{
										{
											LinkID: "risk-factors-score",
											Answer: []domain.FHIRQuestionnaireResponseItemAnswer{
												{
													ValueCoding: &domain.FHIRCoding{
														Display: "Yes",
													},
													ValueInteger: &score,
												},
											},
										},
									},
								},
							},
						}

						return questionnaireReponse, nil
					})
				mh.FHIR.EXPECT().GetFHIRQuestionnaire(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRQuestionnaireRelayPayload, error) {
						questionnaireName := "Breast Cancer Screening"
						linkID := "risk-assessment"
						questionLinkID := "risk-assessment-question-one"
						valueDecimalOne := 1.0
						questionnaireRelayPayload := &domain.FHIRQuestionnaireRelayPayload{
							Resource: &domain.FHIRQuestionnaire{
								ID:    &ID,
								Name:  &questionnaireName,
								Title: &questionnaireName,
								Item: []*domain.FHIRQuestionnaireItem{
									{
										LinkID: &linkID,
										Item: []*domain.FHIRQuestionnaireItem{
											{
												ID:     &questionLinkID,
												LinkID: &questionLinkID,
												Meta:   &domain.FHIRMeta{},
												AnswerOption: []*domain.FHIRQuestionnaireItemAnswerOption{
													{
														Extension: []*domain.Extension{
															{
																URL:          "http://hl7.org/fhir/StructureDefinition/ordinalValue",
																ValueDecimal: &valueDecimalOne,
															},
														},
														ValueCoding: &domain.FHIRCoding{
															Display: "Yes",
														},
													},
												},
											},
										},
									},
								},
							},
						}
						return questionnaireRelayPayload, nil
					})
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return patient, nil
					})
				mh.PubSub.EXPECT().NotifySegmentation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, data dto.SegmentationPayload) error {
						return nil
					})
				mh.FHIR.EXPECT().CreateFHIRRiskAssessment(mock.Anything, mock.Anything).
					RunAndReturn(func(context1 context.Context, input *domain.FHIRRiskAssessmentInput) (*domain.FHIRRiskAssessmentRelayPayload, error) {
						return riskassessment, nil
					})

				return args{ctx: ctx, input: dto.QuestionnaireResponse{}, questionnaireID: ID, encounterID: ID}
			},
			wantErr: false,
		},
		{
			name: "Happy Case - Create questionnaire response and generate review summary - Breast Cancer - Low Risk",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					})
				mh.FHIR.EXPECT().CreateFHIRQuestionnaireResponse(mock.Anything, mock.Anything).
					RunAndReturn(func(context1 context.Context, input *domain.FHIRQuestionnaireResponse) (*domain.FHIRQuestionnaireResponse, error) {
						score := 1
						questionnaireReponse := &domain.FHIRQuestionnaireResponse{
							ID: &ID,
							Item: []domain.FHIRQuestionnaireResponseItem{
								{
									LinkID: "risk-assessment",
									Item: []domain.FHIRQuestionnaireResponseItem{
										{
											LinkID: "risk-factors-score",
											Answer: []domain.FHIRQuestionnaireResponseItemAnswer{
												{
													ValueCoding: &domain.FHIRCoding{
														Display: "Yes",
													},
													ValueInteger: &score,
												},
											},
										},
									},
								},
							},
						}

						return questionnaireReponse, nil
					})
				mh.FHIR.EXPECT().GetFHIRQuestionnaire(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRQuestionnaireRelayPayload, error) {
						questionnaireName := "Breast Cancer Screening"
						linkID := "risk-assessment"
						questionLinkID := "risk-assessment-question-one"
						valueDecimalOne := 1.0
						questionnaireRelayPayload := &domain.FHIRQuestionnaireRelayPayload{
							Resource: &domain.FHIRQuestionnaire{
								ID:    &ID,
								Name:  &questionnaireName,
								Title: &questionnaireName,
								Item: []*domain.FHIRQuestionnaireItem{
									{
										LinkID: &linkID,
										Item: []*domain.FHIRQuestionnaireItem{
											{
												ID:     &questionLinkID,
												LinkID: &questionLinkID,
												Meta:   &domain.FHIRMeta{},
												AnswerOption: []*domain.FHIRQuestionnaireItemAnswerOption{
													{
														Extension: []*domain.Extension{
															{
																URL:          "http://hl7.org/fhir/StructureDefinition/ordinalValue",
																ValueDecimal: &valueDecimalOne,
															},
														},
														ValueCoding: &domain.FHIRCoding{
															Display: "Yes",
														},
													},
												},
											},
										},
									},
								},
							},
						}
						return questionnaireRelayPayload, nil
					})
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return patient, nil
					})
				mh.PubSub.EXPECT().NotifySegmentation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, data dto.SegmentationPayload) error {
						return nil
					})
				mh.FHIR.EXPECT().CreateFHIRRiskAssessment(mock.Anything, mock.Anything).
					RunAndReturn(func(context1 context.Context, input *domain.FHIRRiskAssessmentInput) (*domain.FHIRRiskAssessmentRelayPayload, error) {
						return riskassessment, nil
					})

				return args{ctx: ctx, input: dto.QuestionnaireResponse{}, questionnaireID: ID, encounterID: ID}
			},
			wantErr: false,
		},
		{
			name: "Happy Case - Create questionnaire response and generate review summary - Prostate Cancer - High Risk",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					})
				mh.FHIR.EXPECT().CreateFHIRQuestionnaireResponse(mock.Anything, mock.Anything).
					RunAndReturn(func(context1 context.Context, input *domain.FHIRQuestionnaireResponse) (*domain.FHIRQuestionnaireResponse, error) {
						score := 1
						questionnaireReponse := &domain.FHIRQuestionnaireResponse{
							ID: &ID,
							Item: []domain.FHIRQuestionnaireResponseItem{
								{
									LinkID: "risk-assessment",
									Item: []domain.FHIRQuestionnaireResponseItem{
										{
											LinkID: "risk-factors-score",
											Answer: []domain.FHIRQuestionnaireResponseItemAnswer{
												{
													ValueCoding: &domain.FHIRCoding{
														Display: "Yes",
													},
													ValueInteger: &score,
												},
											},
										},
									},
								},
							},
						}

						return questionnaireReponse, nil
					})
				mh.FHIR.EXPECT().GetFHIRQuestionnaire(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRQuestionnaireRelayPayload, error) {
						questionnaireName := "Prostate Cancer Screening"
						linkID := "risk-assessment"
						questionLinkID := "risk-assessment-question-one"
						valueDecimalOne := 0.0
						questionnaireRelayPayload := &domain.FHIRQuestionnaireRelayPayload{
							Resource: &domain.FHIRQuestionnaire{
								ID:    &ID,
								Name:  &questionnaireName,
								Title: &questionnaireName,
								Item: []*domain.FHIRQuestionnaireItem{
									{
										LinkID: &linkID,
										Item: []*domain.FHIRQuestionnaireItem{
											{
												ID:     &questionLinkID,
												LinkID: &questionLinkID,
												Meta:   &domain.FHIRMeta{},
												AnswerOption: []*domain.FHIRQuestionnaireItemAnswerOption{
													{
														Extension: []*domain.Extension{
															{
																URL:          "http://hl7.org/fhir/StructureDefinition/ordinalValue",
																ValueDecimal: &valueDecimalOne,
															},
														},
														ValueCoding: &domain.FHIRCoding{
															Display: "Yes",
														},
													},
												},
											},
										},
									},
								},
							},
						}
						return questionnaireRelayPayload, nil
					})
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return patient, nil
					})
				mh.PubSub.EXPECT().NotifySegmentation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, data dto.SegmentationPayload) error {
						return nil
					})
				mh.FHIR.EXPECT().CreateFHIRRiskAssessment(mock.Anything, mock.Anything).
					RunAndReturn(func(context1 context.Context, input *domain.FHIRRiskAssessmentInput) (*domain.FHIRRiskAssessmentRelayPayload, error) {
						return riskassessment, nil
					})

				return args{ctx: ctx, input: dto.QuestionnaireResponse{}, questionnaireID: ID, encounterID: ID}
			},
			wantErr: false,
		},
		{
			name: "Happy Case - Create questionnaire response and generate review summary - Prostate Cancer - Low Risk",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					})
				mh.FHIR.EXPECT().CreateFHIRQuestionnaireResponse(mock.Anything, mock.Anything).
					RunAndReturn(func(context1 context.Context, input *domain.FHIRQuestionnaireResponse) (*domain.FHIRQuestionnaireResponse, error) {
						score := 0
						questionnaireReponse := &domain.FHIRQuestionnaireResponse{
							ID: &ID,
							Item: []domain.FHIRQuestionnaireResponseItem{
								{
									LinkID: "risk-assessment",
									Item: []domain.FHIRQuestionnaireResponseItem{
										{
											LinkID: "risk-factors-score",
											Answer: []domain.FHIRQuestionnaireResponseItemAnswer{
												{
													ValueCoding: &domain.FHIRCoding{
														Display: "Yes",
													},
													ValueInteger: &score,
												},
											},
										},
									},
								},
							},
						}

						return questionnaireReponse, nil
					})
				mh.FHIR.EXPECT().GetFHIRQuestionnaire(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRQuestionnaireRelayPayload, error) {
						questionnaireName := "Prostate Cancer Screening"
						linkID := "risk-assessment"
						questionLinkID := "risk-assessment-question-one"
						valueDecimalOne := 0.0
						questionnaireRelayPayload := &domain.FHIRQuestionnaireRelayPayload{
							Resource: &domain.FHIRQuestionnaire{
								ID:    &ID,
								Name:  &questionnaireName,
								Title: &questionnaireName,
								Item: []*domain.FHIRQuestionnaireItem{
									{
										LinkID: &linkID,
										Item: []*domain.FHIRQuestionnaireItem{
											{
												ID:     &questionLinkID,
												LinkID: &questionLinkID,
												Meta:   &domain.FHIRMeta{},
												AnswerOption: []*domain.FHIRQuestionnaireItemAnswerOption{
													{
														Extension: []*domain.Extension{
															{
																URL:          "http://hl7.org/fhir/StructureDefinition/ordinalValue",
																ValueDecimal: &valueDecimalOne,
															},
														},
														ValueCoding: &domain.FHIRCoding{
															Display: "Yes",
														},
													},
												},
											},
										},
									},
								},
							},
						}
						return questionnaireRelayPayload, nil
					})
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return patient, nil
					})
				mh.PubSub.EXPECT().NotifySegmentation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, data dto.SegmentationPayload) error {
						return nil
					})
				mh.FHIR.EXPECT().CreateFHIRRiskAssessment(mock.Anything, mock.Anything).
					RunAndReturn(func(context1 context.Context, input *domain.FHIRRiskAssessmentInput) (*domain.FHIRRiskAssessmentRelayPayload, error) {
						return riskassessment, nil
					})
				return args{ctx: ctx, input: dto.QuestionnaireResponse{}, questionnaireID: ID, encounterID: ID}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - fail to get patient",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					})
				mh.FHIR.EXPECT().CreateFHIRQuestionnaireResponse(mock.Anything, mock.Anything).
					RunAndReturn(func(context1 context.Context, input *domain.FHIRQuestionnaireResponse) (*domain.FHIRQuestionnaireResponse, error) {
						score := 0
						questionnaireReponse := &domain.FHIRQuestionnaireResponse{
							ID: &ID,
							Item: []domain.FHIRQuestionnaireResponseItem{
								{
									LinkID: "symptoms",
									Item: []domain.FHIRQuestionnaireResponseItem{
										{
											LinkID: "symptoms-score",
											Answer: []domain.FHIRQuestionnaireResponseItemAnswer{
												{
													ValueInteger: &score,
												},
											},
										},
									},
								},
								{
									LinkID: "risk-factors",
									Item: []domain.FHIRQuestionnaireResponseItem{
										{
											LinkID: "risk-factors-score",
											Answer: []domain.FHIRQuestionnaireResponseItemAnswer{
												{
													ValueInteger: &score,
												},
											},
										},
									},
								},
							},
						}

						return questionnaireReponse, nil
					})
				mh.FHIR.EXPECT().GetFHIRQuestionnaire(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRQuestionnaireRelayPayload, error) {
						return questionnaireRelayPayload, nil
					})
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, input: dto.QuestionnaireResponse{}, questionnaireID: ID, encounterID: ID}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - fail to publish to pusbsub",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					})
				mh.FHIR.EXPECT().CreateFHIRQuestionnaireResponse(mock.Anything, mock.Anything).
					RunAndReturn(func(context1 context.Context, input *domain.FHIRQuestionnaireResponse) (*domain.FHIRQuestionnaireResponse, error) {
						score := 0
						questionnaireReponse := &domain.FHIRQuestionnaireResponse{
							ID: &ID,
							Item: []domain.FHIRQuestionnaireResponseItem{
								{
									LinkID: "risk-assessment",
									Item: []domain.FHIRQuestionnaireResponseItem{
										{
											LinkID: "risk-factors-score",
											Answer: []domain.FHIRQuestionnaireResponseItemAnswer{
												{
													ValueCoding: &domain.FHIRCoding{
														Display: "Yes",
													},
													ValueInteger: &score,
												},
											},
										},
									},
								},
							},
						}

						return questionnaireReponse, nil
					})
				mh.FHIR.EXPECT().GetFHIRQuestionnaire(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRQuestionnaireRelayPayload, error) {
						questionnaireName := "Prostate Cancer Screening"
						linkID := "risk-assessment"
						questionLinkID := "risk-assessment-question-one"
						valueDecimalOne := 0.0
						questionnaireRelayPayload := &domain.FHIRQuestionnaireRelayPayload{
							Resource: &domain.FHIRQuestionnaire{
								ID:    &ID,
								Name:  &questionnaireName,
								Title: &questionnaireName,
								Item: []*domain.FHIRQuestionnaireItem{
									{
										LinkID: &linkID,
										Item: []*domain.FHIRQuestionnaireItem{
											{
												ID:     &questionLinkID,
												LinkID: &questionLinkID,
												Meta:   &domain.FHIRMeta{},
												AnswerOption: []*domain.FHIRQuestionnaireItemAnswerOption{
													{
														Extension: []*domain.Extension{
															{
																URL:          "http://hl7.org/fhir/StructureDefinition/ordinalValue",
																ValueDecimal: &valueDecimalOne,
															},
														},
														ValueCoding: &domain.FHIRCoding{
															Display: "Yes",
														},
													},
												},
											},
										},
									},
								},
							},
						}
						return questionnaireRelayPayload, nil
					})
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return patient, nil
					})
				mh.FHIR.EXPECT().CreateFHIRRiskAssessment(mock.Anything, mock.Anything).
					RunAndReturn(func(context1 context.Context, input *domain.FHIRRiskAssessmentInput) (*domain.FHIRRiskAssessmentRelayPayload, error) {
						return riskassessment, nil
					})
				mh.PubSub.EXPECT().NotifySegmentation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, data dto.SegmentationPayload) error {
						return fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, input: dto.QuestionnaireResponse{}, questionnaireID: ID, encounterID: ID}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			_, err := clinicalUsecase.CreateQuestionnaireResponse(args.ctx, args.questionnaireID, args.encounterID, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.CreateQuestionnaireResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCasesClinicalImpl_SimpleQuestionnaireResponse(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx                     context.Context
		questionnaireResponseID string
	}

	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		want    []*domain.SimpleQuestionnaireResponse
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully get the questionnaire response",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRQuestionnaireResponse(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRQuestionnaireResponseRelayPayload, error) {
						ID := uuid.NewString()
						score := 2
						response := &domain.FHIRQuestionnaireResponseRelayPayload{
							Resource: &domain.FHIRQuestionnaireResponse{
								ID: &ID,
								Item: []domain.FHIRQuestionnaireResponseItem{
									{
										LinkID: "symptoms",
										Item: []domain.FHIRQuestionnaireResponseItem{
											{
												LinkID: "symptoms-score",
												Answer: []domain.FHIRQuestionnaireResponseItemAnswer{
													{
														ValueCoding: &domain.FHIRCoding{
															Display: "Yes",
														},
														ValueInteger: &score,
													},
												},
											},
										},
									},
								},
							},
						}
						return response, nil
					})
				return args{ctx: ctx, questionnaireResponseID: uuid.NewString()}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Missing questionnaire response ID",
			setup: func(mh *usecaseMock.Mocks) args {
				return args{ctx: ctx}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get questionnaire",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRQuestionnaireResponse(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRQuestionnaireResponseRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, questionnaireResponseID: uuid.NewString()}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			_, err := clinicalUsecase.SimpleQuestionnaireResponse(args.ctx, args.questionnaireResponseID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.SimpleQuestionnaireResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
