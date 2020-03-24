package core_test

import (
	"strings"
	"testing"

	core "github.com/bitrise-io/workflow-editor-core"
	"github.com/c2fo/testify/require"
)

const (
	minimalValidSecrets    = "{}"
	minimalValidBitriseYML = "format_version: 1.3.0"
)

func TestValidateBitriseConfigAndSecret(t *testing.T) {
	validConfigs := []string{
		minimalValidBitriseYML,
		//
		`format_version: 1.3.0
app:
  envs:
  - KEY_ONE: value one
workflows:
  empty_wf:
    steps: []
`,
		//
	}
	validSecrets := []string{
		minimalValidSecrets,
		"",
		"#",
		" ",
		"\n",
		"{}",
		"envs:",
		"envs: ",
		`envs: []`,
		//
		`envs:
- SECRET_ONE: secret value one`,
		//
		`envs:
- SECRET_ONE: secret value one
`,
		//
	}
	validWithWarning :=
		`
format_version: 1.1.0
trigger_map:
- pattern: ci/quick
  workflow: _prepare_and_setup
workflows:
 _prepare_and_setup:
  description: desc
`

	t.Log("Valid combinations")
	{
		for _, aValidConfig := range validConfigs {
			for _, aValidSecret := range validSecrets {
				t.Log("Config: ", aValidConfig)
				t.Log("Secret: ", aValidSecret)
				warningItems, err := core.ValidateBitriseConfigAndSecret(aValidConfig, aValidSecret)
				require.NoError(t, err)
				require.Nil(t, warningItems)
			}
		}
	}

	t.Log("Valid config with warnings")
	{
		warnings, err := core.ValidateBitriseConfigAndSecret(validWithWarning,
			minimalValidSecrets)
		require.NoError(t, err)
		require.Equal(t, "workflow (_prepare_and_setup) defined in trigger item (pattern: ci/quick && is_pull_request_allowed: false -> workflow: _prepare_and_setup), but utility workflows can't be triggered directly", warnings.Config[0])
	}

	t.Log("Invalid configs - empty")
	{
		{
			_, err := core.ValidateBitriseConfigAndSecret(``, minimalValidSecrets)
			require.Error(t, err)
			require.True(t, strings.Contains(err.Error(), "Validation failed: Config validation error: "), err.Error())
		}

		{
			_, err := core.ValidateBitriseConfigAndSecret(`{}`, minimalValidSecrets)
			require.Error(t, err)
			require.True(t, strings.Contains(err.Error(), "Validation failed: Config validation error: Failed to get config (bitrise.yml) from base 64 data, err: Failed to parse bitrise config, error: missing format_version"), err.Error())
		}
	}

	t.Log("Invalid configs - 1")
	{
		_, err := core.ValidateBitriseConfigAndSecret(`format_version: 1.3.0
app:
  envs:
  - A
`, minimalValidSecrets)

		require.Error(t, err)
		require.True(t, strings.Contains(err.Error(), "Validation failed: Config validation error: Failed to get config (bitrise.yml) from base 64 data, err: Failed to parse bitrise config, error: yaml: unmarshal errors:\n  line 4: cannot unmarshal !!str `A` into models.EnvironmentItemModel"), err.Error())
	}

	t.Log("Invalid configs - missing format_version")
	{
		_, err := core.ValidateBitriseConfigAndSecret(`
app:
  envs:
  - KEY_ONE: value one
workflows:
  empty_wf:
    steps: []
`, minimalValidSecrets)

		require.Error(t, err)
		require.True(t, strings.Contains(err.Error(), "Validation failed: Config validation error: Failed to get config (bitrise.yml) from base 64 data, err: Failed to parse bitrise config, error: missing format_version"), err.Error())
	}

	t.Log("Invalid secrets - envs as empty hash")
	{
		_, err := core.ValidateBitriseConfigAndSecret(minimalValidBitriseYML, "envs: {}")
		require.Error(t, err)
		require.True(t, strings.Contains(err.Error(), "Validation failed: Secret validation error: Failed to get inventory from base 64 data, err: yaml: unmarshal errors:\n  line 1: cannot unmarshal !!map into []models.EnvironmentItemModel"), err.Error())
	}

	t.Log("Invalid secrets - envs as hash with value")
	{
		_, err := core.ValidateBitriseConfigAndSecret(minimalValidBitriseYML, `envs:
  KEY_ONE: value one`)
		require.Error(t, err)
		require.True(t, strings.Contains(err.Error(), "Validation failed: Secret validation error: Failed to get inventory from base 64 data, err: yaml: unmarshal errors:\n  line 2: cannot unmarshal !!map into []models.EnvironmentItemModel"), err.Error())
	}
}
