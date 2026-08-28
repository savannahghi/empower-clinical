package domain

import (
	"fmt"
	"strings"
)

// OperationOutcome is documented here http://hl7.org/fhir/StructureDefinition/OperationOutcome
type OperationOutcome struct {
	Issue []OperationOutcomeIssue `json:"issue,omitempty"`
}

type OperationOutcomeIssue struct {
	Severity    IssueSeverity        `json:"severity,omitempty"`
	Code        IssueType            `json:"code,omitempty"`
	Details     *FHIRCodeableConcept `json:"details,omitempty"`
	Diagnostics *string              `json:"diagnostics,omitempty"`
	Expression  []string             `json:"expression,omitempty"`
}

func (oo OperationOutcome) Error() string {
	var builder strings.Builder

	for _, issue := range oo.Issue {
		builder.WriteString(formatIssue(issue))
		builder.WriteString("\n\n")
	}

	return strings.TrimSpace(builder.String())
}

// formatIssue formats a single OperationOutcomeIssue into a string
func formatIssue(issue OperationOutcomeIssue) string {
	return fmt.Sprintf(
		"Severity: %s\nCode: %s\nDetails: %s\nDiagnostics: %s\nExpression: %s",
		issue.Severity.String(),
		issue.Code.String(),
		issue.Details.Text,
		*issue.Diagnostics,
		strings.Join(issue.Expression, ", "),
	)
}
