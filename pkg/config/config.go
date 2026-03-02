package config

type Config struct {

	// Configurations (read only)

	Configs struct {

		// Folderpath with tasks ".json" files ("".state.json")
		Path string

		// Seconds between folder read; for new/updated tasks detection.
		SecsBetweenSync int
	}
}
