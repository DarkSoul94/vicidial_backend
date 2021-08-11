package mock

// RepoMock ...
type RepoMock struct {
}

// NewRepo ...
func NewRepo() *RepoMock {
	return &RepoMock{}
}

// Close ...
func (r *RepoMock) Close() error {
	return nil
}
