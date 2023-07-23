package modifiers_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zalgonoise/logx"
	"github.com/zalgonoise/logx/handlers/texth"

	"github.com/zalgonoise/apbrc/processor"
	"github.com/zalgonoise/apbrc/processor/modifiers"
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
		fpsModifierFormat = "%s=%d\r\n"
		logger            = logx.New(texth.New(os.Stderr))
	)

	for _, testcase := range []struct {
		name          string
		targetTestDir string
		targetDir     string
		targetFile    string
		modifier      processor.Modifier
	}{
		{
			name:          "FPSModifier/Original/Complete",
			targetTestDir: "fps/complete_orig",
			targetDir:     "/Engine/Config",
			targetFile:    "/BaseEngine.ini",
			modifier:      modifiers.NewFPSModifier(60, 300, 300, logger),
		},
		{
			name:          "FPSModifier/Original/Short",
			targetTestDir: "fps/short_orig",
			targetDir:     "/Engine/Config",
			targetFile:    "/BaseEngine.ini",
			modifier:      modifiers.NewFPSModifier(60, 300, 300, logger),
		},
		{
			name:          "FPSModifier/Fake",
			targetTestDir: "fps/short_fake",
			targetDir:     "/Engine/Config",
			targetFile:    "/BaseEngine.ini",
			modifier: modifiers.NewKeyValueModifier(
				"/Engine/Config/BaseEngine.ini", logger,
				modifiers.KeyValue[int]{
					Key:    "SomeAttributeKey",
					Value:  600,
					Format: fpsModifierFormat,
				},
				modifiers.KeyValue[int]{
					Key:    "OtherAttributeKey",
					Value:  700,
					Format: fpsModifierFormat,
				},
				modifiers.KeyValue[int]{
					Key:    "LastAttributeKey",
					Value:  800,
					Format: fpsModifierFormat,
				},
			),
		},
		{
			name:          "SprintLock/Original/Complete",
			targetTestDir: "sprint/complete_orig",
			targetDir:     "/APBGame/Config",
			targetFile:    "/DefaultInput.ini",
			modifier:      modifiers.NewSprintLockModifier(true, logger),
		},
		{
			name:          "CrouchLock/Original/Complete",
			targetTestDir: "crouch/complete_orig",
			targetDir:     "/APBGame/Config",
			targetFile:    "/DefaultInput.ini",
			modifier:      modifiers.NewCrouchLockModifier(true, logger),
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

			path := baseDir + testcase.targetTestDir + topLevel
			err = testcase.modifier.Apply(path)
			require.NoError(t, err)

			data, err = os.ReadFile(fs.basePath + fs.filePath)
			require.NoError(t, err)

			wants, err := os.ReadFile(baseDir + testcase.targetTestDir + resultsDir + testcase.targetFile)
			require.NoError(t, err)
			require.Equal(t, wants, data)
		})
	}
}
