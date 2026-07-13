package cloud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	naeoserrors "github.com/NAEOS-foundation/naeos/internal/errors"
)

type CommandRunner interface {
	Run(name string, args []string, dir string) ([]byte, error)
}

type ExecCommandRunner struct{}

func (r *ExecCommandRunner) Run(name string, args []string, dir string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("%s %s: %w\nstderr: %s", name, strings.Join(args, " "), err, stderr.String())
	}
	return stdout.Bytes(), nil
}

type TerraformRunner struct {
	WorkDir string
	Runner  CommandRunner
}

func NewTerraformRunner(workDir string) *TerraformRunner {
	return &TerraformRunner{
		WorkDir: workDir,
		Runner:  &ExecCommandRunner{},
	}
}

func NewTerraformRunnerWithRunner(workDir string, runner CommandRunner) *TerraformRunner {
	return &TerraformRunner{
		WorkDir: workDir,
		Runner:  runner,
	}
}

func (t *TerraformRunner) writeHCL(hcl string) error {
	if err := os.MkdirAll(t.WorkDir, 0o755); err != nil {
		return naeoserrors.Wrap(naeoserrors.ErrCloud, "failed to create terraform working directory", err)
	}
	path := filepath.Join(t.WorkDir, "main.tf")
	if err := os.WriteFile(path, []byte(hcl), 0o644); err != nil {
		return naeoserrors.Wrap(naeoserrors.ErrCloud, "failed to write HCL file", err)
	}
	return nil
}

func (t *TerraformRunner) Init() error {
	_, err := t.Runner.Run("terraform", []string{"init", "-input=false"}, t.WorkDir)
	if err != nil {
		return naeoserrors.Wrap(naeoserrors.ErrCloud, "terraform init failed", err)
	}
	return nil
}

type PlanOutput struct {
	Changes struct {
		Add    int `json:"add"`
		Change int `json:"change"`
		Destroy int `json:"destroy"`
	} `json:"changes"`
}

func (t *TerraformRunner) Plan() (*PlanOutput, error) {
	raw, err := t.Runner.Run("terraform", []string{"plan", "-input=false", "-json"}, t.WorkDir)
	if err != nil {
		return nil, naeoserrors.Wrap(naeoserrors.ErrCloud, "terraform plan failed", err)
	}
	return parsePlanJSON(raw)
}

func parsePlanOutput(raw []byte) (*PlanOutput, error) {
	return parsePlanJSON(raw)
}

func parsePlanJSON(raw []byte) (*PlanOutput, error) {
	lines := strings.Split(strings.TrimSpace(string(raw)), "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}
		var event struct {
			Type        string      `json:"@message"`
			Changes     *struct {
				Add    int `json:"add"`
				Change int `json:"change"`
				Destroy int `json:"destroy"`
			} `json:"changes"`
		}
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			continue
		}
		if event.Changes != nil {
			return &PlanOutput{
				Changes: struct {
					Add    int `json:"add"`
					Change int `json:"change"`
					Destroy int `json:"destroy"`
				}{
					Add:     event.Changes.Add,
					Change:  event.Changes.Change,
					Destroy: event.Changes.Destroy,
				},
			}, nil
		}
	}
	return &PlanOutput{}, nil
}

func (t *TerraformRunner) Apply() error {
	_, err := t.Runner.Run("terraform", []string{"apply", "-auto-approve", "-input=false"}, t.WorkDir)
	if err != nil {
		return naeoserrors.Wrap(naeoserrors.ErrCloud, "terraform apply failed", err)
	}
	return nil
}

func (t *TerraformRunner) ApplyDestroy() error {
	_, err := t.Runner.Run("terraform", []string{"destroy", "-auto-approve", "-input=false"}, t.WorkDir)
	if err != nil {
		return naeoserrors.Wrap(naeoserrors.ErrCloud, "terraform destroy failed", err)
	}
	return nil
}

type OutputValue struct {
	Value     interface{} `json:"value"`
	Sensitive bool        `json:"sensitive"`
	Type      interface{} `json:"type"`
}

func (t *TerraformRunner) Output() (map[string]OutputValue, error) {
	raw, err := t.Runner.Run("terraform", []string{"output", "-json"}, t.WorkDir)
	if err != nil {
		return nil, naeoserrors.Wrap(naeoserrors.ErrCloud, "terraform output failed", err)
	}
	var outputs map[string]OutputValue
	if err := json.Unmarshal(raw, &outputs); err != nil {
		return nil, naeoserrors.Wrap(naeoserrors.ErrCloud, "failed to parse terraform output", err)
	}
	return outputs, nil
}

func (t *TerraformRunner) Deploy(hcl string) error {
	if err := t.writeHCL(hcl); err != nil {
		return err
	}
	if err := t.Init(); err != nil {
		return err
	}
	if err := t.Apply(); err != nil {
		return err
	}
	return nil
}

func (t *TerraformRunner) DestroyAll() error {
	if err := t.Init(); err != nil {
		return err
	}
	if err := t.ApplyDestroy(); err != nil {
		return err
	}
	return nil
}

func TempWorkDir(prefix string) (string, error) {
	dir, err := os.MkdirTemp("", prefix+"-terraform-*")
	if err != nil {
		return "", naeoserrors.Wrap(naeoserrors.ErrCloud, "failed to create temp directory", err)
	}
	return dir, nil
}
