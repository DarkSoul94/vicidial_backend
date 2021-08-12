package vicidial_backend

// Usecase ...
type Usecase interface {
	AddLeads(leads []map[string]interface{})
	UpdateLead(lead map[string]interface{})
}
