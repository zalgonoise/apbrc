package processor_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zalgonoise/logx"
	"github.com/zalgonoise/logx/handlers/texth"

	"github.com/zalgonoise/apbrc/config"
	"github.com/zalgonoise/apbrc/processor"
)

type testFS struct {
	basePath string
	filePath string

	origData []byte
}

func (fs *testFS) Rollback() error {
	path := fs.basePath + fs.filePath

	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}

	if _, err = f.Write(fs.origData); err != nil {
		return err
	}

	if err = f.Close(); err != nil {
		return err
	}

	return nil
}

func TestProcessor_Run(t *testing.T) {
	var (
		baseDir       = "./modifiers/internal/testdata/"
		resultsDir    = "/results"
		binariesDir   = "/Binaries"
		topLevel      = "/APB Reloaded"
		defaultTarget = "/Engine/Config"
		filename      = "/BaseEngine.ini"
		logger        = logx.New(texth.New(os.Stderr))
	)

	for _, testcase := range []struct {
		name      string
		targetDir string
		cfg       *config.Config
		err       error
	}{
		{
			name:      "Original/Complete",
			targetDir: "complete_orig",
			cfg: &config.Config{
				Path: "./modifiers/internal/testdata/complete_orig" + topLevel + binariesDir,
				FrameRate: &config.FrameRateConfig{
					MinRate:      60,
					MaxRate:      300,
					SmoothedRate: 300,
				},
			},
		},
		{
			name:      "Original/Short",
			targetDir: "short_orig",
			cfg: &config.Config{
				Path: "./modifiers/internal/testdata/short_orig" + topLevel + binariesDir,
				FrameRate: &config.FrameRateConfig{
					MinRate:      60,
					MaxRate:      300,
					SmoothedRate: 300,
				},
			},
		},
		{
			name:      "Fail/InvalidPath",
			targetDir: "short_orig",
			cfg: &config.Config{
				Path: "./modifiers/internal/testdata/test_fake",
				FrameRate: &config.FrameRateConfig{
					MinRate:      60,
					MaxRate:      300,
					SmoothedRate: 300,
				},
			},
			err: processor.ErrInvalidPath,
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			fs := &testFS{
				basePath: baseDir + testcase.targetDir + topLevel,
				filePath: defaultTarget + filename,
			}

			data, err := os.ReadFile(fs.basePath + fs.filePath)
			require.NoError(t, err)

			fs.origData = data
			defer fs.Rollback()

			proc := processor.New(testcase.cfg, logger)

			if err = proc.Run(); err != nil {
				require.ErrorIs(t, err, testcase.err)

				return
			}

			data, err = os.ReadFile(fs.basePath + fs.filePath)
			require.NoError(t, err)

			wants, err := os.ReadFile(baseDir + testcase.targetDir + resultsDir + filename)
			require.NoError(t, err)
			require.Equal(t, wants, data)
		})
	}
}
