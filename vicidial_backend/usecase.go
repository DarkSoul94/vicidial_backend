package vicidial_backend

import "github.com/DarkSoul94/vicidial_backend/models"

// Usecase ...
type Usecase interface {
	AddLeads(leads []models.Lead)
	UpdateLead(lead models.Lead)
}
