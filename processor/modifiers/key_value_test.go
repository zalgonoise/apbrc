package modifiers_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zalgonoise/logx"
	"github.com/zalgonoise/logx/handlers/texth"

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
		defaultTarget     = "/Engine/Config"
		filename          = "/BaseEngine.ini"
		fpsModifierFormat = "%s=%d\r\n"
		logger            = logx.New(texth.New(os.Stderr))
	)

	for _, testcase := range []struct {
		name      string
		targetDir string
		modifier  modifiers.KeyValueModifier[int]
	}{
		{
			name:      "Original/Complete",
			targetDir: "complete_orig",
			modifier:  modifiers.NewFPSModifier(60, 300, 300, logger),
		},
		{
			name:      "Original/Short",
			targetDir: "short_orig",
			modifier:  modifiers.NewFPSModifier(60, 300, 300, logger),
		},
		{
			name:      "Fake",
			targetDir: "test_fake",
			modifier: modifiers.NewKeyValueModifier(
				defaultTarget+filename, logger,
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

			path := baseDir + testcase.targetDir + topLevel
			err = testcase.modifier.Apply(path)
			require.NoError(t, err)

			data, err = os.ReadFile(fs.basePath + fs.filePath)
			require.NoError(t, err)

			wants, err := os.ReadFile(baseDir + testcase.targetDir + resultsDir + filename)
			require.NoError(t, err)
			require.Equal(t, wants, data)
		})
	}
}
