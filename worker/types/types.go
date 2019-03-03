package types

// SubscriberInformation is blockchain concept
type SubscriberInformation struct {
	Roles []string `json:"roles"`
	ID    string   `json:"id"`
}

// Publish is blockchain operation
type Publish struct {
	ProcessID              string                `json:"processID"`
	ApplicationRelatedData []string              `json:"applicationRelatedData"`
	ProcessRelatedData     string                `json:"processRelatedData"`
	SubscriberInformation  SubscriberInformation `json:"subscriberInformation"`
}

// ChangeCondition is blockchain operation
type ChangeCondition struct {
	ProcessID string `json:"processID"`
	Condition string `json:"condition"`
}

// ChangeProcessData is blockchain operation
type ChangeProcessData struct {
	ProcessID                       string   `json:"processID"`
	IsProcessRelatedDataChanged     bool     `json:"isProcessRelatedDataChanged"`
	ApplicationRelatedData          []string `json:"applicationRelatedData"`
	IsApplicationRelatedDataChanged bool     `json:"isApplicationRelatedDataChanged"`
	ProcessRelatedData              string   `json:"processRelatedData"`
}

// ChangeSubscriberInformation is blockchain operation
type ChangeSubscriberInformation struct {
	ProcessID             string                `json:"processID"`
	SubscriberInformation SubscriberInformation `json:"subscriberInformation"`
}
