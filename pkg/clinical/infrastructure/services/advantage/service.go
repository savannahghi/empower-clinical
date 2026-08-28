package advantage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/savannahghi/serverutils"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/dto"
	auth "github.com/savannahghi/empower-clinical/pkg/clinical/infrastructure/services/authutils"
)

var (
	AdvantageBaseURL  = serverutils.MustGetEnvVar("ADVANTAGE_BASE_URL")
	segmentationPath  = "/api/segments/segment/clinical/"
	smsPath           = "/api/notifications/send_sms/"
	schedulePath      = "/api/scheduling/schedules/"
	slotsPath         = "/api/scheduling/slots/"
	createCheckInPath = "/api/scheduling/appointments/"
)

// AdvantageService represents methods that can be used to communicate with the advantage server
type AdvantageService interface {
	SegmentPatient(ctx context.Context, payload dto.SegmentationPayload) error
	SendSMS(ctx context.Context, workstationID, branchID string, payload dto.SMSPayload) error
	GetSchedules(ctx context.Context, headers *dto.AdvantageHeaders) ([]*dto.Schedule, error)
	GetSlots(ctx context.Context, startDate string, scheduleID string, headers *dto.AdvantageHeaders) ([]*dto.Slot, error)
	CreateCheckin(ctx context.Context, checkIn *dto.Checkin, headers *dto.AdvantageHeaders) error
}

// ServiceAdvantageImpl represents advantage server's implementations
type ServiceAdvantageImpl struct {
	client auth.OAuthClientService
}

// NewServiceAdvantage is the advantage server's service constructor
func NewServiceAdvantage(authUtils auth.OAuthClientService) *ServiceAdvantageImpl {
	return &ServiceAdvantageImpl{
		client: authUtils,
	}
}

// SegmentPatient is used to create segmentation information in advantage
func (s *ServiceAdvantageImpl) SegmentPatient(ctx context.Context, payload dto.SegmentationPayload) error {
	url := fmt.Sprintf("%s/%s", AdvantageBaseURL, segmentationPath)

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	body := bytes.NewReader(payloadBytes)

	req, err := s.newAuthenticatedRequest(http.MethodPost, url, body, nil)
	if err != nil {
		return err
	}

	resp, err := s.makeRequest(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}

// SendSMS is used to send an SMS message to a patient.
func (s *ServiceAdvantageImpl) SendSMS(ctx context.Context, workstationID, branchID string, payload dto.SMSPayload) error {
	url := fmt.Sprintf("%s%s", AdvantageBaseURL, smsPath)

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	body := bytes.NewReader(payloadBytes)

	headers := &dto.AdvantageHeaders{
		Workstation: workstationID,
		Branch:      branchID,
	}

	req, err := s.newAuthenticatedRequest(http.MethodPost, url, body, headers)
	if err != nil {
		return err
	}

	resp, err := s.makeRequest(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}

// GetSchedules is used to retrieve days schedule
func (s *ServiceAdvantageImpl) GetSchedules(ctx context.Context, headers *dto.AdvantageHeaders) ([]*dto.Schedule, error) {
	path := fmt.Sprintf("%s%s?actor=PRACTITIONER&fields=id,description,specialty,practitioner_data&page_size=1000", AdvantageBaseURL, schedulePath)

	req, err := s.newAuthenticatedRequest(http.MethodGet, path, nil, headers)
	if err != nil {
		return nil, err
	}

	resp, err := s.makeRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response *dto.AdvantageResponse

	if err := json.Unmarshal(data, &response); err != nil {
		return nil, err
	}

	var output []*dto.Schedule

	err = mapstructure.Decode(response.Results, &output)
	if err != nil {
		return nil, err
	}

	return output, nil
}

// GetSlots is used to retrieve the available slots for a given schedule and day
func (s *ServiceAdvantageImpl) GetSlots(ctx context.Context, startDate string, scheduleID string, headers *dto.AdvantageHeaders) ([]*dto.Slot, error) {
	url := fmt.Sprintf("%s%s?start=%s&fields=id,start,end&schedule_id=%s&ordering=start&status=FREE", AdvantageBaseURL, slotsPath, startDate, scheduleID)

	req, err := s.newAuthenticatedRequest(http.MethodGet, url, nil, headers)
	if err != nil {
		return nil, err
	}

	resp, err := s.makeRequest(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response *dto.AdvantageResponse

	if err := json.Unmarshal(data, &response); err != nil {
		return nil, err
	}

	var slot []*dto.Slot

	err = mapstructure.Decode(response.Results, &slot)
	if err != nil {
		return nil, err
	}

	return slot, nil
}

// CreateCheckin is used to create a checkin appointment
func (s *ServiceAdvantageImpl) CreateCheckin(ctx context.Context, checkIn *dto.Checkin, headers *dto.AdvantageHeaders) error {
	url := fmt.Sprintf("%s%s", AdvantageBaseURL, createCheckInPath)

	payloadBytes, err := json.Marshal(checkIn)
	if err != nil {
		return err
	}

	body := bytes.NewReader(payloadBytes)

	req, err := s.newAuthenticatedRequest(http.MethodPost, url, body, headers)
	if err != nil {
		return err
	}

	resp, err := s.makeRequest(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var output *dto.Checkin

	if err := json.Unmarshal(data, &output); err != nil {
		return fmt.Errorf("%w. %s", err, string(data))
	}

	return nil
}

// makeRequest is a helper function to make a http request to advantage
func (s *ServiceAdvantageImpl) makeRequest(req *http.Request) (*http.Response, error) {
	httpClient := &http.Client{Timeout: time.Second * 30}
	return httpClient.Do(req)
}

// newAuthenticatedRequest creates a new HTTP request, attaches authentication token and other necessary headers.
func (s *ServiceAdvantageImpl) newAuthenticatedRequest(method, url string, body io.Reader, headers *dto.AdvantageHeaders) (*http.Request, error) {
	var request *http.Request

	var err error

	switch method {
	case http.MethodPost:
		request, err = http.NewRequest(http.MethodPost, url, body)
		if err != nil {
			return nil, err
		}
	case http.MethodGet:
		request, err = http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported method %s", method)
	}

	token, err := s.client.Authenticate()
	if err != nil {
		return nil, err
	}

	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

	if headers != nil {
		request.Header.Set("X-Workstation", headers.Workstation)
		request.Header.Set("X-Branch", headers.Branch)
		request.Header.Set("X-Variant", headers.Variant)

		if headers.Department != "" && headers.Cluster != "" {
			request.Header.Set("X-Department", headers.Department)
			request.Header.Set("X-Cluster", headers.Cluster)
		}
	}

	return request, nil
}
