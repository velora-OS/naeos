package cloud

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type mockRunner struct {
	stdout []byte
	err    error
}

func (m *mockRunner) Run(name string, args []string, dir string) ([]byte, error) {
	return m.stdout, m.err
}

func TestTerraformRunnerCreation(t *testing.T) {
	dir := t.TempDir()
	tr := NewTerraformRunner(dir)
	if tr.WorkDir != dir {
		t.Errorf("expected workdir %s, got %s", dir, tr.WorkDir)
	}
	if tr.Runner == nil {
		t.Error("expected non-nil runner")
	}
}

func TestTerraformRunnerWithMock(t *testing.T) {
	dir := t.TempDir()
	mock := &mockRunner{stdout: []byte("{}")}
	tr := NewTerraformRunnerWithRunner(dir, mock)
	if tr.Runner != mock {
		t.Error("expected mock runner")
	}
}

func TestWriteHCL(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "nested", "tf")
	tr := NewTerraformRunner(dir)

	hcl := `resource "null_resource" "test" {}`
	if err := tr.writeHCL(hcl); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(dir, "main.tf"))
	if err != nil {
		t.Fatalf("failed to read written file: %v", err)
	}
	if string(data) != hcl {
		t.Errorf("HCL content mismatch:\ngot:  %s\nwant: %s", string(data), hcl)
	}
}

func TestInit(t *testing.T) {
	dir := t.TempDir()
	mock := &mockRunner{stdout: []byte("ok")}
	tr := NewTerraformRunnerWithRunner(dir, mock)

	if err := tr.Init(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestInitError(t *testing.T) {
	dir := t.TempDir()
	mock := &mockRunner{err: fmt.Errorf("command failed")}
	tr := NewTerraformRunnerWithRunner(dir, mock)

	if err := tr.Init(); err == nil {
		t.Error("expected error from Init")
	}
}

func TestPlanJSON(t *testing.T) {
	planJSON := `{"@level":"info","@message":"Plan: 2 to add, 0 to change, 1 to destroy.","changes":{"add":2,"change":0,"destroy":1}}`
	mock := &mockRunner{stdout: []byte(planJSON)}
	tr := NewTerraformRunnerWithRunner(t.TempDir(), mock)

	output, err := tr.Plan()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if output.Changes.Add != 2 {
		t.Errorf("expected 2 adds, got %d", output.Changes.Add)
	}
	if output.Changes.Destroy != 1 {
		t.Errorf("expected 1 destroy, got %d", output.Changes.Destroy)
	}
}

func TestPlanMultiLineJSON(t *testing.T) {
	planJSON := strings.Join([]string{
		`{"@level":"info","@message":"Terraform will perform the following actions:"}`,
		`{"@level":"info","@message":"Plan: 1 to add.","changes":{"add":1,"change":0,"destroy":0}}`,
	}, "\n")
	mock := &mockRunner{stdout: []byte(planJSON)}
	tr := NewTerraformRunnerWithRunner(t.TempDir(), mock)

	output, err := tr.Plan()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if output.Changes.Add != 1 {
		t.Errorf("expected 1 add, got %d", output.Changes.Add)
	}
}

func TestApply(t *testing.T) {
	dir := t.TempDir()
	mock := &mockRunner{stdout: []byte("ok")}
	tr := NewTerraformRunnerWithRunner(dir, mock)

	if err := tr.Apply(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestApplyError(t *testing.T) {
	dir := t.TempDir()
	mock := &mockRunner{err: fmt.Errorf("apply failed")}
	tr := NewTerraformRunnerWithRunner(dir, mock)

	if err := tr.Apply(); err == nil {
		t.Error("expected error from Apply")
	}
}

func TestDestroyAll(t *testing.T) {
	dir := t.TempDir()
	mock := &mockRunner{stdout: []byte("ok")}
	tr := NewTerraformRunnerWithRunner(dir, mock)

	if err := tr.DestroyAll(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestOutput(t *testing.T) {
	outputJSON := `{"bucket_id":{"value":"my-bucket","sensitive":false,"type":"string"}}`
	mock := &mockRunner{stdout: []byte(outputJSON)}
	tr := NewTerraformRunnerWithRunner(t.TempDir(), mock)

	output, err := tr.Output()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v, ok := output["bucket_id"]; !ok {
		t.Error("missing bucket_id in output")
	} else if v.Value != "my-bucket" {
		t.Errorf("expected value 'my-bucket', got %v", v.Value)
	}
}

func TestOutputError(t *testing.T) {
	mock := &mockRunner{err: fmt.Errorf("no outputs")}
	tr := NewTerraformRunnerWithRunner(t.TempDir(), mock)

	_, err := tr.Output()
	if err == nil {
		t.Error("expected error from Output")
	}
}

func TestDeploy(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "deploy-test")
	mock := &mockRunner{stdout: []byte("ok")}
	tr := NewTerraformRunnerWithRunner(dir, mock)

	hcl := `resource "null_resource" "test" {}`
	if err := tr.Deploy(hcl); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(dir, "main.tf"))
	if err != nil {
		t.Fatalf("failed to read main.tf: %v", err)
	}
	if string(data) != hcl {
		t.Errorf("HCL mismatch: got %s", string(data))
	}
}

func TestDeployWritesHCLBeforeInit(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "hcl-order")
	mock := &mockRunner{stdout: []byte("ok")}
	tr := NewTerraformRunnerWithRunner(dir, mock)

	hcl := `resource "null_resource" "test" {}`
	if err := tr.Deploy(hcl); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(dir, "main.tf"))
	if err != nil {
		t.Fatalf("main.tf not created: %v", err)
	}
	if string(data) != hcl {
		t.Errorf("HCL content mismatch: got %s", string(data))
	}
}

func TestTempWorkDir(t *testing.T) {
	dir, err := TempWorkDir("naeos-test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer os.RemoveAll(dir)

	if !strings.HasPrefix(filepath.Base(dir), "naeos-test-terraform-") {
		t.Errorf("unexpected dir name: %s", filepath.Base(dir))
	}

	info, err := os.Stat(dir)
	if err != nil {
		t.Fatalf("dir does not exist: %v", err)
	}
	if !info.IsDir() {
		t.Error("expected directory")
	}
}

func TestParsePlanJSONNoChanges(t *testing.T) {
	mock := &mockRunner{stdout: []byte(`{"@message":"no changes"}`)}
	tr := NewTerraformRunnerWithRunner(t.TempDir(), mock)

	output, err := tr.Plan()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if output.Changes.Add != 0 || output.Changes.Destroy != 0 {
		t.Errorf("expected zero changes, got add=%d destroy=%d", output.Changes.Add, output.Changes.Destroy)
	}
}

func TestDeployError(t *testing.T) {
	mock := &mockRunner{err: fmt.Errorf("init failed")}
	tr := NewTerraformRunnerWithRunner(t.TempDir(), mock)

	err := tr.Deploy(`resource "null_resource" "test" {}`)
	if err == nil {
		t.Error("expected error from Deploy")
	}
}
