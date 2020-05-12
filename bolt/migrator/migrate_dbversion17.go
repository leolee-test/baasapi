package migrator

func (m *Migrator) updateExtensionsToDBVersion18() error {
	legacySettings, err := m.settingsService.Settings()
	if err != nil {
		return err
	}

	if legacySettings.Baask8sInterval == "" {
		legacySettings.Baask8sInterval = "1m"
	}

	return m.settingsService.UpdateSettings(legacySettings)
}
