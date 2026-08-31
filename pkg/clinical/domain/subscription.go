package domain

type FHIRSubscriptionInput struct {
	ID        *string                 `json:"id,omitempty"`
	Meta      *FHIRMetaInput          `json:"meta,omitempty"`
	Extension []Extension             `json:"extension,omitempty"`
	Status    SubscriptionTypeEnum    `json:"status"`
	Contact   []FHIRContactPointInput `json:"contact,omitempty"`
	End       *string                 `json:"end,omitempty"`
	Reason    string                  `json:"reason"`
	Criteria  string                  `json:"criteria"`
	Error     *string                 `json:"error,omitempty"`
	Channel   FHIRSubscriptionChannel `json:"channel"`
}

// FHIRSubscriptionChannel the channel on which to report matches to the criteria
type FHIRSubscriptionChannelInput struct {
	ID       *string  `json:"id,omitempty"`
	Type     *string  `json:"type"`
	Endpoint *string  `json:"endpoint,omitempty"`
	Payload  *string  `json:"payload,omitempty"`
	Header   []string `json:"header,omitempty"`
}

// FHIRSubscription models a subscription output
type FHIRSubscription struct {
	ID                *string                 `json:"id,omitempty"`
	Meta              *FHIRMeta               `json:"meta,omitempty"`
	ImplicitRules     *string                 `json:"implicitRules,omitempty"`
	Language          *string                 `json:"language,omitempty"`
	Text              *FHIRNarrative          `json:"text,omitempty"`
	Extension         []*Extension            `json:"extension,omitempty"`
	ModifierExtension []*Extension            `json:"modifierExtension,omitempty"`
	Identifier        []*FHIRIdentifier       `json:"identifier,omitempty"`
	Status            SubscriptionTypeEnum    `json:"status"`
	Contact           []FHIRContactPoint      `json:"contact,omitempty"`
	End               *string                 `json:"end,omitempty"`
	Reason            string                  `json:"reason"`
	Criteria          string                  `json:"criteria"`
	Error             *string                 `json:"error,omitempty"`
	Channel           FHIRSubscriptionChannel `json:"channel"`
}

// FHIRSubscriptionChannel ...
type FHIRSubscriptionChannel struct {
	ID                *string           `json:"id,omitempty"`
	Meta              *FHIRMeta         `json:"meta,omitempty"`
	ImplicitRules     *string           `json:"implicitRules,omitempty"`
	Language          *string           `json:"language,omitempty"`
	Text              *FHIRNarrative    `json:"text,omitempty"`
	Extension         []*Extension      `json:"extension,omitempty"`
	ModifierExtension []*Extension      `json:"modifierExtension,omitempty"`
	Identifier        []*FHIRIdentifier `json:"identifier,omitempty"`
	Type              *string           `json:"type"`
	Endpoint          *string           `json:"endpoint,omitempty"`
	Payload           *string           `json:"payload,omitempty"`
	Header            []string          `json:"header,omitempty"`
}
