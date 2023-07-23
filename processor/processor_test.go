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
		baseDir     = "./modifiers/internal/testdata/"
		resultsDir  = "/results"
		binariesDir = "/Binaries"
		topLevel    = "/APB Reloaded"
		logger      = logx.New(texth.New(os.Stderr))
	)

	for _, testcase := range []struct {
		name          string
		targetTestDir string
		targetDir     string
		targetFile    string
		cfg           *config.Config
		err           error
	}{
		{
			name:          "FPS/Original/Complete",
			targetTestDir: "fps/complete_orig",
			targetDir:     "/Engine/Config",
			targetFile:    "/BaseEngine.ini",
			cfg: &config.Config{
				Path: baseDir + "fps/complete_orig" + topLevel + binariesDir,
				FrameRate: &config.FrameRateConfig{
					MinRate:      60,
					MaxRate:      300,
					SmoothedRate: 300,
				},
			},
		},
		{
			name:          "Original/Short",
			targetTestDir: "fps/short_orig",
			targetDir:     "/Engine/Config",
			targetFile:    "/BaseEngine.ini",
			cfg: &config.Config{
				Path: baseDir + "fps/short_orig" + topLevel + binariesDir,
				FrameRate: &config.FrameRateConfig{
					MinRate:      60,
					MaxRate:      300,
					SmoothedRate: 300,
				},
			},
		},
		{
			name:          "Fail/InvalidPath",
			targetTestDir: "fps/short_orig",
			targetDir:     "/Engine/Config",
			targetFile:    "/BaseEngine.ini",
			cfg: &config.Config{
				Path: baseDir + "fps/short_fake",
				FrameRate: &config.FrameRateConfig{
					MinRate:      60,
					MaxRate:      300,
					SmoothedRate: 300,
				},
			},
			err: processor.ErrInvalidPath,
		},
		{
			name:          "Sprint/Original/Complete",
			targetTestDir: "sprint/complete_orig",
			targetDir:     "/APBGame/Config",
			targetFile:    "/DefaultInput.ini",
			cfg: &config.Config{
				Path: baseDir + "sprint/complete_orig" + topLevel + binariesDir,
				Input: &config.InputConfig{
					SprintLock: true,
				},
			},
		},
		{
			name:          "Crouch/Original/Complete",
			targetTestDir: "crouch/complete_orig",
			targetDir:     "/APBGame/Config",
			targetFile:    "/DefaultInput.ini",
			cfg: &config.Config{
				Path: baseDir + "crouch/complete_orig" + topLevel + binariesDir,
				Input: &config.InputConfig{
					CrouchHold: true,
				},
			},
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			fs := &testFS{
				basePath: baseDir + testcase.targetTestDir + topLevel,
				filePath: testcase.targetDir + testcase.targetFile,
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

			wants, err := os.ReadFile(baseDir + testcase.targetTestDir + resultsDir + testcase.targetFile)
			require.NoError(t, err)
			require.Equal(t, wants, data)
		})
	}
}
