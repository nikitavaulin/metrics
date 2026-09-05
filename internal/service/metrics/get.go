package metricsservice

func (s *MetricsService) Get(name string) (any, error) {
	return s.storage.Get(name)
}

func (s *MetricsService) GetList() (map[string]any, error) {
	return s.storage.GetList()
}
