package version

var (
	version   = "unknown"
	gitSHA    = "unknown"
	buildTime = "unknown"
)

func Version() string   { return version }
func GitSHA() string    { return gitSHA }
func BuildTime() string { return buildTime }
