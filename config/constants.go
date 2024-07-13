package config

const (
	// ApplicationName name of the base application
	ApplicationName string = "gextend-bash"
	// ConfigDirEnv environment pointing to the main config dir
	ConfigDirEnv string = "XTRA_BASH_CONFIG_DIR"
	// ConfigDirName config dirname under ~/.config
	ConfigDirName string = "gextend-bash"
	// ConstDirMode permissions on newly created configuration directories
	ConstDirMode int = 0755
	// ConstFileMode permissions on newly created configuration files
	ConstFileMode int = 0644
)
