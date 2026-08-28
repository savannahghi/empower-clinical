package domain

import (
	"time"

	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/scalarutils"
)

// FHIRComposition definition: a set of healthcare-related information that is assembled together into a single logical package that provides a single coherent statement of meaning, establishes its own context and that has clinical attestation with regard to who is making the statement. a composition defines the structure and narrative content necessary for a document. however, a composition alone does not constitute a document. rather, the composition must be the first entry in a bundle where bundle.type=document, and any other resources referenced from composition must be included as subsequent entries in the bundle (for example patient, practitioner, encounter, etc.).
type FHIRComposition struct {
	// The logical id of the resource, as used in the URL for the resource. Once assigned, this value never changes.
	ID *string `json:"id,omitempty"`

	// A human-readable narrative that contains a summary of the resource and can be used to represent the content of the resource to a human. The narrative need not encode all the structured data, but is required to contain sufficient detail to make it "clinically safe" for a human to just read the narrative. Resource definitions may define what content should be represented in the narrative to ensure clinical safety.
	Text *FHIRNarrative `json:"text,omitempty"`

	Name *string `json:"name,omitempty"`

	// A version-independent identifier for the Composition. This identifier stays constant as the composition is changed over time.
	Identifier []*FHIRIdentifier `json:"identifier,omitempty"`

	// The workflow/clinical status of this composition. The status is a marker for the clinical standing of the document.
	Status *CompositionStatusEnum `json:"status,omitempty"`

	// Specifies the particular kind of composition (e.g. History and Physical, Discharge Summary, Progress Note). This usually equates to the purpose of making the composition.
	Type *FHIRCodeableConcept `json:"type,omitempty"`

	// A categorization for the type of the composition - helps for indexing and searching. This may be implied by or derived from the code specified in the Composition Type.
	Category []*FHIRCodeableConcept `json:"category,omitempty"`

	// Who or what the composition is about. The composition can be about a person, (patient or healthcare practitioner), a device (e.g. a machine) or even a group of subjects (such as a document about a herd of livestock, or a set of patients that share a common exposure).
	Subject []*FHIRReference `json:"subject,omitempty"`

	// Describes the clinical encounter or type of care this documentation is associated with.
	Encounter *FHIRReference `json:"encounter,omitempty"`

	// The composition editing time, when the composition was last logically changed by the author.
	Date *scalarutils.Date `json:"date,omitempty"`

	// Identifies who is responsible for the information in the composition, not necessarily who typed it in.
	Author []*FHIRReference `json:"author,omitempty"`

	// Official human-readable label for the composition.
	Title *string `json:"title,omitempty"`

	// The code specifying the level of confidentiality of the Composition.
	Confidentiality *scalarutils.Code `json:"confidentiality,omitempty"`

	// A participant who has attested to the accuracy of the composition/document.
	Attester []*FHIRCompositionAttester `json:"attester,omitempty"`

	// Identifies the organization or group who is responsible for ongoing maintenance of and access to the composition/document information.
	Custodian *FHIRReference `json:"custodian,omitempty"`

	// Relationships that this composition has with other compositions or documents that already exist.
	RelatesTo []*FHIRCompositionRelatesto `json:"relatesTo,omitempty"`

	// The clinical service, such as a colonoscopy or an appendectomy, being documented.
	Event []*FHIRCompositionEvent `json:"event,omitempty"`

	// The root of the sections that make up the composition.
	Section []*FHIRCompositionSection `json:"section,omitempty"`

	// Meta stores more information about the resource
	Meta *FHIRMeta `json:"meta,omitempty"`

	// Extension is an optional element that provides additional information not captured in the basic resource definition
	Extension []*FHIRExtension `json:"extension,omitempty"`
}

// FHIRCompositionAttester definition: a set of healthcare-related information that is assembled together into a single logical package that provides a single coherent statement of meaning, establishes its own context and that has clinical attestation with regard to who is making the statement. a composition defines the structure and narrative content necessary for a document. however, a composition alone does not constitute a document. rather, the composition must be the first entry in a bundle where bundle.type=document, and any other resources referenced from composition must be included as subsequent entries in the bundle (for example patient, practitioner, encounter, etc.).
type FHIRCompositionAttester struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// The type of attestation the authenticator offers.
	Mode *CompositionAttesterModeEnum `json:"mode,omitempty"`

	// When the composition was attested by the party.
	Time *time.Time `json:"time,omitempty"`

	// Who attested the composition in the specified way.
	Party *FHIRReference `json:"party,omitempty"`
}

// FHIRCompositionAttesterInput is the input type for CompositionAttester
type FHIRCompositionAttesterInput struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// The type of attestation the authenticator offers.
	Mode *CompositionAttesterModeEnum `json:"mode,omitempty"`

	// When the composition was attested by the party.
	Time *scalarutils.DateTime `json:"time,omitempty"`

	// Who attested the composition in the specified way.
	Party *FHIRReferenceInput `json:"party,omitempty"`
}

// FHIRCompositionEvent definition: a set of healthcare-related information that is assembled together into a single logical package that provides a single coherent statement of meaning, establishes its own context and that has clinical attestation with regard to who is making the statement. a composition defines the structure and narrative content necessary for a document. however, a composition alone does not constitute a document. rather, the composition must be the first entry in a bundle where bundle.type=document, and any other resources referenced from composition must be included as subsequent entries in the bundle (for example patient, practitioner, encounter, etc.).
type FHIRCompositionEvent struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// This list of codes represents the main clinical acts, such as a colonoscopy or an appendectomy, being documented. In some cases, the event is inherent in the typeCode, such as a "History and Physical Report" in which the procedure being documented is necessarily a "History and Physical" act.
	Code *scalarutils.Code `json:"code,omitempty"`

	// The period of time covered by the documentation. There is no assertion that the documentation is a complete representation for this period, only that it documents events during this time.
	Period *FHIRPeriod `json:"period,omitempty"`

	// The description and/or reference of the event(s) being documented. For example, this could be used to document such a colonoscopy or an appendectomy.
	Detail []*FHIRReference `json:"detail,omitempty"`
}

// FHIRCompositionEventInput is the input type for CompositionEvent
type FHIRCompositionEventInput struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// This list of codes represents the main clinical acts, such as a colonoscopy or an appendectomy, being documented. In some cases, the event is inherent in the typeCode, such as a "History and Physical Report" in which the procedure being documented is necessarily a "History and Physical" act.
	Code *scalarutils.Code `json:"code,omitempty"`

	// The period of time covered by the documentation. There is no assertion that the documentation is a complete representation for this period, only that it documents events during this time.
	Period *FHIRPeriodInput `json:"period,omitempty"`

	// The description and/or reference of the event(s) being documented. For example, this could be used to document such a colonoscopy or an appendectomy.
	Detail []*FHIRReferenceInput `json:"detail,omitempty"`
}

// FHIRCompositionInput is the input type for Composition
type FHIRCompositionInput struct {
	// The logical id of the resource, as used in the URL for the resource. Once assigned, this value never changes.
	ID *string `json:"id,omitempty"`

	Text *FHIRNarrativeInput `json:"text,omitempty"`

	Name *string `json:"name,omitempty"`

	// A version-independent identifier for the Composition. This identifier stays constant as the composition is changed over time.
	Identifier *FHIRIdentifierInput `json:"identifier,omitempty"`

	// The workflow/clinical status of this composition. The status is a marker for the clinical standing of the document.
	Status *CompositionStatusEnum `json:"status,omitempty"`

	// Specifies the particular kind of composition (e.g. History and Physical, Discharge Summary, Progress Note). This usually equates to the purpose of making the composition.
	Type *FHIRCodeableConceptInput `json:"type,omitempty"`

	// A categorization for the type of the composition - helps for indexing and searching. This may be implied by or derived from the code specified in the Composition Type.
	Category []*FHIRCodeableConceptInput `json:"category,omitempty"`

	// Who or what the composition is about. The composition can be about a person, (patient or healthcare practitioner), a device (e.g. a machine) or even a group of subjects (such as a document about a herd of livestock, or a set of patients that share a common exposure).
	Subject []*FHIRReferenceInput `json:"subject,omitempty"`

	// Describes the clinical encounter or type of care this documentation is associated with.
	Encounter *FHIRReferenceInput `json:"encounter,omitempty"`

	// The composition editing time, when the composition was last logically changed by the author.
	Date *scalarutils.Date `json:"date,omitempty"`

	// Identifies who is responsible for the information in the composition, not necessarily who typed it in.
	Author []*FHIRReferenceInput `json:"author,omitempty"`

	// Official human-readable label for the composition.
	Title *string `json:"title,omitempty"`

	// The code specifying the level of confidentiality of the Composition.
	Confidentiality *scalarutils.Code `json:"confidentiality,omitempty"`

	// A participant who has attested to the accuracy of the composition/document.
	Attester []*FHIRCompositionAttesterInput `json:"attester,omitempty"`

	// Identifies the organization or group who is responsible for ongoing maintenance of and access to the composition/document information.
	Custodian *FHIRReferenceInput `json:"custodian,omitempty"`

	// Relationships that this composition has with other compositions or documents that already exist.
	RelatesTo []*FHIRCompositionRelatestoInput `json:"relatesTo,omitempty"`

	// The clinical service, such as a colonoscopy or an appendectomy, being documented.
	Event []*FHIRCompositionEventInput `json:"event,omitempty"`

	// The root of the sections that make up the composition.
	Section []*FHIRCompositionSectionInput `json:"section,omitempty"`

	// Meta stores more information about the resource
	Meta *FHIRMetaInput `json:"meta,omitempty"`

	// Extension is an optional element that provides additional information not captured in the basic resource definition
	Extension []*FHIRExtension `json:"extension,omitempty"`
}

// FHIRCompositionRelatesto definition: a set of healthcare-related information that is assembled together into a single logical package that provides a single coherent statement of meaning, establishes its own context and that has clinical attestation with regard to who is making the statement. a composition defines the structure and narrative content necessary for a document. however, a composition alone does not constitute a document. rather, the composition must be the first entry in a bundle where bundle.type=document, and any other resources referenced from composition must be included as subsequent entries in the bundle (for example patient, practitioner, encounter, etc.).
type FHIRCompositionRelatesto struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// The type of relationship that this composition has with anther composition or document.
	Code *scalarutils.Code `json:"code,omitempty"`

	// The target composition/document of this relationship.
	TargetIdentifier *FHIRIdentifier `json:"targetIdentifier,omitempty"`

	// The target composition/document of this relationship.
	TargetReference *FHIRReference `json:"targetReference,omitempty"`
}

// FHIRCompositionRelatestoInput is the input type for CompositionRelatesto
type FHIRCompositionRelatestoInput struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// The type of relationship that this composition has with anther composition or document.
	Code *scalarutils.Code `json:"code,omitempty"`

	// The target composition/document of this relationship.
	TargetIdentifier *FHIRIdentifierInput `json:"targetIdentifier,omitempty"`

	// The target composition/document of this relationship.
	TargetReference *FHIRReferenceInput `json:"targetReference,omitempty"`
}

// FHIRCompositionSection definition: a set of healthcare-related information that is assembled together into a single logical package that provides a single coherent statement of meaning, establishes its own context and that has clinical attestation with regard to who is making the statement. a composition defines the structure and narrative content necessary for a document. however, a composition alone does not constitute a document. rather, the composition must be the first entry in a bundle where bundle.type=document, and any other resources referenced from composition must be included as subsequent entries in the bundle (for example patient, practitioner, encounter, etc.).
type FHIRCompositionSection struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// The label for this particular section.  This will be part of the rendered content for the document, and is often used to build a table of contents.
	Title *string `json:"title,omitempty"`

	// A code identifying the kind of content contained within the section. This must be consistent with the section title.
	Code *FHIRCodeableConceptInput `json:"code,omitempty"`

	// Identifies who is responsible for the information in this section, not necessarily who typed it in.
	Author []*FHIRReference `json:"author,omitempty"`

	// The actual focus of the section when it is not the subject of the composition, but instead represents something or someone associated with the subject such as (for a patient subject) a spouse, parent, fetus, or donor. If not focus is specified, the focus is assumed to be focus of the parent section, or, for a section in the Composition itself, the subject of the composition. Sections with a focus SHALL only include resources where the logical subject (patient, subject, focus, etc.) matches the section focus, or the resources have no logical subject (few resources).
	Focus *FHIRReference `json:"focus,omitempty"`

	// A human-readable narrative that contains the attested content of the section, used to represent the content of the resource to a human. The narrative need not encode all the structured data, but is required to contain sufficient detail to make it "clinically safe" for a human to just read the narrative.
	Text *FHIRNarrative `json:"text,omitempty"`

	// How the entry list was prepared - whether it is a working list that is suitable for being maintained on an ongoing basis, or if it represents a snapshot of a list of items from another source, or whether it is a prepared list where items may be marked as added, modified or deleted.
	Mode *scalarutils.Code `json:"mode,omitempty"`

	// Specifies the order applied to the items in the section entries.
	OrderedBy *FHIRCodeableConcept `json:"orderedBy,omitempty"`

	// A reference to the actual resource from which the narrative in the section is derived.
	Entry []*FHIRReference `json:"entry,omitempty"`

	// If the section is empty, why the list is empty. An empty section typically has some text explaining the empty reason.
	EmptyReason *FHIRCodeableConcept `json:"emptyReason,omitempty"`

	// A nested sub-section within this section.
	Section []*FHIRCompositionSection `json:"section,omitempty"`
}

// FHIRCompositionSectionInput is the input type for CompositionSection
type FHIRCompositionSectionInput struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// The label for this particular section.  This will be part of the rendered content for the document, and is often used to build a table of contents.
	Title *string `json:"title,omitempty"`

	// A code identifying the kind of content contained within the section. This must be consistent with the section title.
	Code *FHIRCodeableConceptInput `json:"code,omitempty"`

	// Identifies who is responsible for the information in this section, not necessarily who typed it in.
	Author []*FHIRReferenceInput `json:"author,omitempty"`

	// The actual focus of the section when it is not the subject of the composition, but instead represents something or someone associated with the subject such as (for a patient subject) a spouse, parent, fetus, or donor. If not focus is specified, the focus is assumed to be focus of the parent section, or, for a section in the Composition itself, the subject of the composition. Sections with a focus SHALL only include resources where the logical subject (patient, subject, focus, etc.) matches the section focus, or the resources have no logical subject (few resources).
	Focus *FHIRReferenceInput `json:"focus,omitempty"`

	// A human-readable narrative that contains the attested content of the section, used to represent the content of the resource to a human. The narrative need not encode all the structured data, but is required to contain sufficient detail to make it "clinically safe" for a human to just read the narrative.
	Text *FHIRNarrativeInput `json:"text,omitempty"`

	// How the entry list was prepared - whether it is a working list that is suitable for being maintained on an ongoing basis, or if it represents a snapshot of a list of items from another source, or whether it is a prepared list where items may be marked as added, modified or deleted.
	Mode *scalarutils.Code `json:"mode,omitempty"`

	// Specifies the order applied to the items in the section entries.
	OrderedBy *FHIRCodeableConceptInput `json:"orderedBy,omitempty"`

	// A reference to the actual resource from which the narrative in the section is derived.
	Entry []*FHIRReferenceInput `json:"entry,omitempty"`

	// If the section is empty, why the list is empty. An empty section typically has some text explaining the empty reason.
	EmptyReason *FHIRCodeableConceptInput `json:"emptyReason,omitempty"`

	// A nested sub-section within this section.
	Section []*FHIRCompositionSectionInput `json:"section,omitempty"`
}

// FHIRCompositionRelayConnection is a Relay connection for Composition
type FHIRCompositionRelayConnection struct {
	Edges []*FHIRCompositionRelayEdge `json:"edges,omitempty"`

	PageInfo *firebasetools.PageInfo `json:"pageInfo,omitempty"`
}

// FHIRCompositionRelayEdge is a Relay edge for Composition
type FHIRCompositionRelayEdge struct {
	Cursor *string `json:"cursor,omitempty"`

	Node *FHIRComposition `json:"node,omitempty"`
}

// PagedFHIRComposition ...
type PagedFHIRComposition struct {
	Compositions    []FHIRComposition
	HasNextPage     bool
	NextCursor      string
	HasPreviousPage bool
	PreviousCursor  string
	TotalCount      int
}

// FHIRCompositionRelayPayload is used to return single instances of Composition
type FHIRCompositionRelayPayload struct {
	Resource *FHIRComposition `json:"resource,omitempty"`
}
