package staging

type Store interface {
	List() ([]*StagingData, error)
}
