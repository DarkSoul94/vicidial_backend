package vicidial_backend

import "context"

// Usecase ...
type Usecase interface {
	HelloWorld(ctx context.Context)
}
