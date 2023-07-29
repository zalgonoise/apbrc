package modifiers_test

import (
	"context"
	"os"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zalgonoise/apbrc/config"
	"github.com/zalgonoise/apbrc/processor/modifiers"
	"github.com/zalgonoise/apbrc/processor/modifiers/engine"
	"github.com/zalgonoise/apbrc/processor/modifiers/input"
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

func TestKeyValueModifier_Apply(t *testing.T) {
	var (
		baseDir           = "./internal/testdata/"
		resultsDir        = "/results"
		topLevel          = "/APB Reloaded"
		fpsModifierFormat = "%s=%d"
		lf                = "\n"
	)

	if runtime.GOOS == "windows" {
		lf = "\r\n"
	}

	for _, testcase := range []struct {
		name          string
		targetTestDir string
		targetDir     string
		targetFile    string
		modifier      modifiers.Applier
	}{
		{
			name:          "FPSModifier/Original/Complete",
			targetTestDir: "fps/complete_orig",
			targetDir:     "/Engine/Config",
			targetFile:    "/BaseEngine.ini",
			modifier: engine.FrameRate(config.FrameRateConfig{
				Cap: 300,
				Min: 60,
				Max: 300,
			}),
		},
		{
			name:          "FPSModifier/Original/Short",
			targetTestDir: "fps/short_orig",
			targetDir:     "/Engine/Config",
			targetFile:    "/BaseEngine.ini",
			modifier: engine.FrameRate(config.FrameRateConfig{
				Cap: 300,
				Min: 60,
				Max: 300,
			}),
		},
		{
			name:          "FPSModifier/Fake",
			targetTestDir: "fps/short_fake",
			targetDir:     "/Engine/Config",
			targetFile:    "/BaseEngine.ini",
			modifier: modifiers.NewModifier(
				"/Engine/Config/BaseEngine.ini",
				modifiers.KeyValue[int]{
					Key:    "SomeAttributeKey",
					Data:   600,
					Format: fpsModifierFormat + lf,
				},
				modifiers.KeyValue[int]{
					Key:    "OtherAttributeKey",
					Data:   700,
					Format: fpsModifierFormat + lf,
				},
				modifiers.KeyValue[int]{
					Key:    "LastAttributeKey",
					Data:   800,
					Format: fpsModifierFormat + lf,
				},
			),
		},
		{
			name:          "SprintLock/Original/Complete",
			targetTestDir: "sprint/complete_orig",
			targetDir:     "/APBGame/Config",
			targetFile:    "/DefaultInput.ini",
			modifier: input.Input(config.InputConfig{
				SprintLock: true,
			}),
		},
		{
			name:          "CrouchLock/Original/Complete",
			targetTestDir: "crouch/complete_orig",
			targetDir:     "/APBGame/Config",
			targetFile:    "/DefaultInput.ini",
			modifier: input.Input(config.InputConfig{
				CrouchHold: true,
			}),
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			ctx := context.Background()
			fs := &testFS{
				basePath: baseDir + testcase.targetTestDir + topLevel,
				filePath: testcase.targetDir + testcase.targetFile,
			}

			data, err := os.ReadFile(fs.basePath + fs.filePath)
			require.NoError(t, err)

			fs.origData = data
			defer fs.Rollback()

			path := baseDir + testcase.targetTestDir + topLevel
			err = testcase.modifier.Apply(ctx, path)
			require.NoError(t, err)

			data, err = os.ReadFile(fs.basePath + fs.filePath)
			require.NoError(t, err)

			wants, err := os.ReadFile(baseDir + testcase.targetTestDir + resultsDir + testcase.targetFile)
			require.NoError(t, err)
			require.Equal(t, wants, data)
		})
	}
}
