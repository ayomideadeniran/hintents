package simulator

type SimulationRequest struct {
	EnvelopeXdr   string            `json:"envelope_xdr"`
	ResultMetaXdr string            `json:"result_meta_xdr"`
	LedgerEntries map[string]string `json:"ledger_entries,omitempty"`
}

type CategorizedEvent struct {
	EventType  string   `json:"event_type"`
	ContractID *string  `json:"contract_id,omitempty"`
	Topics     []string `json:"topics"`
	Data       string   `json:"data"`
}

type SecurityViolation struct {
	Type        string                 `json:"type"`
	Severity    string                 `json:"severity"`
	Description string                 `json:"description"`
	Contract    string                 `json:"contract"`
	Details     map[string]interface{} `json:"details,omitempty"`
}

type SimulationResponse struct {
	Status             string              `json:"status"`
	Error              string              `json:"error,omitempty"`
	Events             []string            `json:"events,omitempty"`
	CategorizedEvents  []CategorizedEvent  `json:"categorized_events,omitempty"`
	Logs               []string            `json:"logs,omitempty"`
	SecurityViolations []SecurityViolation `json:"security_violations,omitempty"`
}
