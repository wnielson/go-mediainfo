package mediainfo

type AnalyzeOptions struct {
	ParseSpeed                 float64
	HasParseSpeed              bool
	TestContinuousFileNames    bool
	HasTestContinuousFileNames bool
}

func defaultAnalyzeOptions() AnalyzeOptions {
	return AnalyzeOptions{ParseSpeed: 0.5, TestContinuousFileNames: false}
}

func normalizeAnalyzeOptions(opts AnalyzeOptions) AnalyzeOptions {
	if !opts.HasParseSpeed {
		opts.ParseSpeed = 0.5
	}
	if !opts.HasTestContinuousFileNames {
		opts.TestContinuousFileNames = false
	}
	if opts.ParseSpeed < 0 {
		opts.ParseSpeed = 0
	}
	if opts.ParseSpeed > 1 {
		opts.ParseSpeed = 1
	}
	return opts
}
