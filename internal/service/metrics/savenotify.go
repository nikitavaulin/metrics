package metricsservice

func (s *MetricsService) notifySave() {
	select {
	case s.saveToFileTrigger <- struct{}{}:
	default:
	}
}
