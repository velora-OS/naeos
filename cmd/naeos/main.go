package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/NAEOS-foundation/naeos/pkg/pipeline"
)

var (
	version    = "dev"
	cliVerbose bool
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string) error {
	root := newRootCommand()
	root.SetArgs(args)
	return root.Execute()
}

func newRootCommand() *cobra.Command {
	root := &cobra.Command{
		Use:           "naeos",
		Short:         "NAEOS CLI",
		Long:          "NAEOS is a declarative engineering runtime for specification-driven project delivery.",
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if cliVerbose {
				fmt.Fprintln(os.Stderr, "[naeos] verbose mode enabled")
			}
		},
	}

	root.PersistentFlags().BoolVar(&cliVerbose, "verbose", false, "enable verbose logging")

	root.AddCommand(newInitCommand())
	root.AddCommand(newRunCommand())
	root.AddCommand(newValidateCommand())
	root.AddCommand(newInspectCommand())
	root.AddCommand(newDoctorCommand())
	root.AddCommand(newRepairCommand())
	root.AddCommand(newScaffoldCommand())
	root.AddCommand(newExportCommand())
	root.AddCommand(newPreviewCommand())
	root.AddCommand(newKernelCommand())
	root.AddCommand(newVersionCommand())
	return root
}

func newVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show NAEOS version",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := fmt.Fprintf(cmd.OutOrStdout(), "naeos %s\n", version)
			return err
		},
	}
}

func newInitCommand() *cobra.Command {
	var output string

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Generate a default NAEOS config file",
		RunE: func(cmd *cobra.Command, args []string) error {
			content := strings.Join([]string{
				"pipeline:",
				"  name: naeos-dev",
				"  mode: development",
				"  verbose: true",
				"  output_dir: ./out",
			}, "\n") + "\n"

			if err := os.WriteFile(output, []byte(content), 0o600); err != nil {
				return fmt.Errorf("write config: %w", err)
			}

			_, err := fmt.Fprintf(cmd.OutOrStdout(), "created %s\n", output)
			return err
		},
	}

	cmd.Flags().StringVarP(&output, "output", "o", "config.example.yaml", "path for the generated config file")
	return cmd
}

func newRunCommand() *cobra.Command {
	var configPath, input, inputFile, outputFormat, outputFile string
	var languages []string

	cmd := &cobra.Command{
		Use:   "run",
		Short: "Execute the NAEOS pipeline",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if configPath == "" {
				return fmt.Errorf("missing required --config")
			}

			inputValue, err := loadInput(input, inputFile)
			if err != nil {
				return err
			}

			cfg, err := pipeline.ConfigFromFile(configPath)
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}
			if cliVerbose {
				cfg.Verbose = true
			}
			if len(languages) > 0 {
				cfg.Languages = languages
			}

			p, err := pipeline.New(cfg)
			if err != nil {
				return fmt.Errorf("failed to construct pipeline: %w", err)
			}

			result, err := p.Run(inputValue)
			if err != nil {
				return fmt.Errorf("pipeline run failed: %w", err)
			}

			payload := map[string]any{
				"pipeline":   cfg.Name,
				"mode":       cfg.Mode,
				"verbose":    cfg.Verbose,
				"output_dir": cfg.OutputDir,
				"artifacts":  len(result.Artifacts),
				"tasks":      len(result.Tasks),
			}

			if len(languages) > 0 {
				payload["languages"] = languages
			}

			rendered, err := renderOutput(payload, outputFormat, func() []byte {
				return []byte(fmt.Sprintf("pipeline=%s mode=%s verbose=%t output_dir=%s\nartifacts=%d tasks=%d\n", result.NEIR.Project, cfg.Mode, cfg.Verbose, cfg.OutputDir, len(result.Artifacts), len(result.Tasks)))
			})
			if err != nil {
				return err
			}

			return writeOrPrint(cmd, rendered, outputFile)
		},
	}

	cmd.Flags().StringVar(&configPath, "config", "", "path to JSON or YAML config file")
	cmd.Flags().StringVar(&input, "input", "", "specification input to process")
	cmd.Flags().StringVar(&inputFile, "input-file", "", "path to a specification file")
	cmd.Flags().StringVar(&outputFormat, "output", "text", "output format: text, json, or yaml")
	cmd.Flags().StringVar(&outputFile, "output-file", "", "optional file path to write the formatted output")
	cmd.Flags().StringArrayVar(&languages, "language", nil, "target language for code generation (go, typescript, python, java, rust)")
	return cmd
}

func newValidateCommand() *cobra.Command {
	var configPath, input, inputFile, outputFormat, outputFile string
	var languages []string

	cmd := &cobra.Command{
		Use:     "validate",
		Aliases: []string{"v"},
		Short:   "Validate a specification using the NAEOS pipeline",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if configPath == "" {
				return fmt.Errorf("missing required --config")
			}

			inputValue, err := loadInput(input, inputFile)
			if err != nil {
				return err
			}

			cfg, err := pipeline.ConfigFromFile(configPath)
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}
			if cliVerbose {
				cfg.Verbose = true
			}
			if len(languages) > 0 {
				cfg.Languages = languages
			}

			p, err := pipeline.New(cfg)
			if err != nil {
				return fmt.Errorf("failed to construct pipeline: %w", err)
			}

			result, err := p.Validate(inputValue)
			if err != nil {
				return fmt.Errorf("pipeline validate failed: %w", err)
			}

			payload := map[string]any{
				"pipeline":   cfg.Name,
				"mode":       cfg.Mode,
				"verbose":    cfg.Verbose,
				"output_dir": cfg.OutputDir,
				"status":     "valid",
				"project":    result.NEIR.Project,
				"source_len": len(result.Source),
			}

			if len(languages) > 0 {
				payload["languages"] = languages
			}

			rendered, err := renderOutput(payload, outputFormat, func() []byte {
				return []byte(fmt.Sprintf("config=%s mode=%s verbose=%t output_dir=%s\nstatus=valid project=%v source_len=%d\n",
					cfg.Name, cfg.Mode, cfg.Verbose, cfg.OutputDir, result.NEIR.Project, len(result.Source)))
			})
			if err != nil {
				return err
			}

			return writeOrPrint(cmd, rendered, outputFile)
		},
	}

	cmd.Flags().StringVar(&configPath, "config", "", "path to JSON or YAML config file")
	cmd.Flags().StringVar(&input, "input", "", "specification input to process")
	cmd.Flags().StringVar(&inputFile, "input-file", "", "path to a specification file")
	cmd.Flags().StringVar(&outputFormat, "output", "text", "output format: text, json, or yaml")
	cmd.Flags().StringVar(&outputFile, "output-file", "", "optional file path to write the formatted output")
	cmd.Flags().StringArrayVar(&languages, "language", nil, "target language for code generation (go, typescript, python, java, rust)")
	return cmd
}

func newInspectCommand() *cobra.Command {
	var configPath, input, inputFile, outputFormat, outputFile string

	cmd := &cobra.Command{
		Use:   "inspect",
		Short: "Inspect the NAEOS pipeline result",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if configPath == "" {
				return fmt.Errorf("missing required --config")
			}

			inputValue, err := loadInput(input, inputFile)
			if err != nil {
				return err
			}

			cfg, err := pipeline.ConfigFromFile(configPath)
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			specInput, err := resolveInput(inputValue)
			if err != nil {
				return fmt.Errorf("resolve input: %w", err)
			}

			p, err := pipeline.New(cfg)
			if err != nil {
				return fmt.Errorf("failed to construct pipeline: %w", err)
			}

			result, err := p.Run(specInput)
			if err != nil {
				return fmt.Errorf("inspect run failed: %w", err)
			}

			projectName := ""
			if result.NEIR != nil && result.NEIR.Project != nil {
				projectName = result.NEIR.Project.Name
			}

			payload := map[string]any{
				"pipeline":     cfg.Name,
				"mode":         cfg.Mode,
				"verbose":      cfg.Verbose,
				"output_dir":   cfg.OutputDir,
				"project":      projectName,
				"input":        specInput,
				"artifacts":    len(result.Artifacts),
				"tasks":        len(result.Tasks),
				"source_words": len(strings.Fields(result.Source)),
			}

			rendered, err := renderOutput(payload, outputFormat, func() []byte {
				return []byte(fmt.Sprintf("pipeline=%s mode=%s verbose=%t output_dir=%s\nproject=%s artifacts=%d tasks=%d input=%q\n", cfg.Name, cfg.Mode, cfg.Verbose, cfg.OutputDir, projectName, len(result.Artifacts), len(result.Tasks), specInput))
			})
			if err != nil {
				return err
			}

			return writeOrPrint(cmd, rendered, outputFile)
		},
	}

	cmd.Flags().StringVar(&configPath, "config", "", "path to JSON or YAML config file")
	cmd.Flags().StringVar(&input, "input", "", "specification input or file path to process")
	cmd.Flags().StringVar(&inputFile, "input-file", "", "path to a specification file")
	cmd.Flags().StringVar(&outputFormat, "output", "text", "output format: text, json, or yaml")
	cmd.Flags().StringVar(&outputFile, "output-file", "", "optional file path to write the formatted output")
	return cmd
}

func newDoctorCommand() *cobra.Command {
	var configPath string

	cmd := &cobra.Command{
		Use:   "doctor",
		Short: "Run diagnostics on the NAEOS configuration",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if configPath == "" {
				return fmt.Errorf("missing required --config")
			}

			cfg, err := pipeline.ConfigFromFile(configPath)
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			status := map[string]any{
				"config":            configPath,
				"pipeline_name":     cfg.Name,
				"mode":              cfg.Mode,
				"output_dir":        cfg.OutputDir,
				"output_dir_exists": false,
			}
			if cfg.OutputDir != "" {
				_, err := os.Stat(cfg.OutputDir)
				status["output_dir_exists"] = err == nil
			}

			data, err := json.MarshalIndent(status, "", "  ")
			if err != nil {
				return fmt.Errorf("encode doctor output: %w", err)
			}
			_, err = fmt.Fprintln(cmd.OutOrStdout(), string(data))
			return err
		},
	}

	cmd.Flags().StringVar(&configPath, "config", "", "path to JSON or YAML config file")
	return cmd
}

func newRepairCommand() *cobra.Command {
	var configPath string

	cmd := &cobra.Command{
		Use:   "repair",
		Short: "Repair the NAEOS output directory",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if configPath == "" {
				return fmt.Errorf("missing required --config")
			}

			cfg, err := pipeline.ConfigFromFile(configPath)
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}
			if cfg.OutputDir == "" {
				return fmt.Errorf("output_dir is not configured")
			}
			if err := os.MkdirAll(cfg.OutputDir, 0o755); err != nil {
				return fmt.Errorf("create output dir: %w", err)
			}

			readmePath := filepath.Join(cfg.OutputDir, "README.md")
			if _, err := os.Stat(readmePath); err != nil {
				if err := os.WriteFile(readmePath, []byte("# Repaired output\n"), 0o600); err != nil {
					return fmt.Errorf("write repair readme: %w", err)
				}
			}

			_, err = fmt.Fprintf(cmd.OutOrStdout(), "repaired %s\n", cfg.OutputDir)
			return err
		},
	}

	cmd.Flags().StringVar(&configPath, "config", "", "path to JSON or YAML config file")
	return cmd
}

func newScaffoldCommand() *cobra.Command {
	var name, output string
	var languages []string

	cmd := &cobra.Command{
		Use:   "scaffold",
		Short: "Generate a starter project scaffold",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if name == "" {
				return fmt.Errorf("missing required --name")
			}
			if output == "" {
				output = name
			}

			if err := os.MkdirAll(output, 0o755); err != nil {
				return fmt.Errorf("create scaffold dir: %w", err)
			}

			readme := fmt.Sprintf("# %s\n\nGenerated by NAEOS scaffold.\n\n## Quick start\n\n- Review spec.yaml\n- Run `go mod tidy`\n- Run `go run ./cmd/app`\n- Visit `http://localhost:8080/health`\n- Visit `http://localhost:8080/api/v1`\n- Visit `http://localhost:8080/api/v1/resources`\n", name)
			spec := fmt.Sprintf("project: %s\nmodules:\n  - name: core\n    path: ./internal/core\nservices:\n  - name: api\n    kind: http\n    port: 8080\n", name)
			makefile := fmt.Sprintf(".PHONY: help\n\nhelp:\n\t@echo 'Available targets:'\n\t@echo '  make help'\n\t@echo '  make scaffold'\n\nscaffold:\n\t@echo 'Scaffold complete for %s'\n", name)
			gitignore := "# Build artifacts\n/bin/\n/out/\n*.log\n"
			dockerfile := "FROM golang:1.22-alpine\nWORKDIR /app\nCOPY . .\nRUN go build ./cmd/app\nCMD [\"/app/app\"]\n"
			ciWorkflow := "name: ci\n\non: [push, pull_request]\n\njobs:\n  build:\n    runs-on: ubuntu-latest\n    steps:\n      - uses: actions/checkout@v4\n      - uses: actions/setup-go@v5\n        with:\n          go-version: '1.22'\n      - run: go test ./...\n"
			goMod := "module github.com/example/" + name + "\n\ngo 1.22\n"
			mainGo := "package main\n\nimport (\n\t\"fmt\"\n\t\"log\"\n\t\"net/http\"\n\n\t\"github.com/example/" + name + "/internal/core\"\n\tcoreconfig \"github.com/example/" + name + "/internal/core/config\"\n\tcorehttp \"github.com/example/" + name + "/internal/core/http\"\n\tcoremiddleware \"github.com/example/" + name + "/internal/core/middleware\"\n)\n\nfunc main() {\n\tcfg := coreconfig.Load(\"config.yaml\")\n\thandler := core.NewHandler(nil)\n\t_ = handler\n\tmux := http.NewServeMux()\n\tmux.HandleFunc(\"/\", func(w http.ResponseWriter, r *http.Request) {\n\t\t_, _ = fmt.Fprintf(w, \"hello from scaffold on port %d\", cfg.Port)\n\t})\n\tmux.HandleFunc(\"/health\", func(w http.ResponseWriter, r *http.Request) {\n\t\t_, _ = fmt.Fprintln(w, \"ok\")\n\t})\n\tmux.HandleFunc(\"/api/v1\", func(w http.ResponseWriter, r *http.Request) {\n\t\t_, _ = fmt.Fprintln(w, \"api v1 ready\")\n\t})\n\tmux.HandleFunc(\"/api/v1/resources\", func(w http.ResponseWriter, r *http.Request) {\n\t\t_, _ = fmt.Fprintln(w, \"resources endpoint\")\n\t})\n\t_ = corehttp.Handler{}\n\twrapped := coremiddleware.LoggingMiddleware{}.Wrap(mux)\n\tlog.Printf(\"listening on :%d\", cfg.Port)\n\tif err := http.ListenAndServe(fmt.Sprintf(\":%d\", cfg.Port), wrapped); err != nil {\n\t\tlog.Fatal(err)\n\t}\n}\n"

			coreModuleFiles := map[string]string{
				"internal/core/README.md":             "# core\n\nStarter module generated by NAEOS.\n",
				"internal/core/package.go":            "package core\n\n// core module placeholder.\n",
				"internal/core/config.yaml":           "name: core\nmodule: core\n",
				"internal/core/handler.go":            "package core\n\n// Handler is a small starter implementation for the core module.\ntype Handler struct {\n\tservice Service\n}\n\nfunc NewHandler(service Service) *Handler {\n\treturn &Handler{service: service}\n}\n",
				"internal/core/repository.go":         "package core\n\n// Repository interface describes the persistence boundary for the core module.\ntype Repository interface {\n\tList() []string\n}\n",
				"internal/core/service.go":            "package core\n\n// Service interface describes the application behavior for the core module.\ntype Service interface {\n\tHandle() string\n}\n",
				"internal/core/domain/model.go":       "package domain\n\n// Model is a sample domain object for the core module.\ntype Model struct {\n\tName string\n}\n",
				"internal/core/http/handler.go":       "package http\n\nimport (\n\t\"fmt\"\n\t\"net/http\"\n)\n\n// Handler is a starter HTTP handler for the core module.\ntype Handler struct{}\n\nfunc (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {\n\t_, _ = fmt.Fprintln(w, \"handler for core\")\n}\n",
				"internal/core/http/router.go":        "package http\n\n// Router is a starter router stub for the core module.\ntype Router struct{}\n",
				"internal/core/middleware/logging.go": "package middleware\n\nimport \"net/http\"\n\n// LoggingMiddleware wraps the next handler with a simple logging stub.\ntype LoggingMiddleware struct{}\n\nfunc (m LoggingMiddleware) Wrap(next http.Handler) http.Handler {\n\treturn http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {\n\t\tnext.ServeHTTP(w, r)\n\t})\n}\n",
				"internal/core/config/config.go":      "package config\n\n// Config is a starter configuration container for the core module.\ntype Config struct {\n\tPort int\n}\n",
				"internal/core/config/load.go":        "package config\n\nimport (\n\t\"encoding/json\"\n\t\"fmt\"\n\t\"os\"\n\n\t\"gopkg.in/yaml.v3\"\n)\n\n// Load returns a starter configuration for the core module.\nfunc Load(path string) Config {\n\tif path == \"\" {\n\t\tpath = \"config.yaml\"\n\t}\n\tdata, err := os.ReadFile(path)\n\tif err != nil {\n\t\treturn Config{Port: 8080}\n\t}\n\tvar cfg Config\n\tswitch ext := fmt.Sprintf(\"%s\", path); ext {\n\tcase \"config.json\":\n\t\t_ = json.Unmarshal(data, &cfg)\n\tdefault:\n\t\t_ = yaml.Unmarshal(data, &cfg)\n\t}\n\tif cfg.Port == 0 {\n\t\tcfg.Port = 8080\n\t}\n\treturn cfg\n}\n",
			}

			files := map[string]string{
				"README.md":                readme,
				"spec.yaml":                spec,
				"Makefile":                 makefile,
				".gitignore":               gitignore,
				"Dockerfile":               dockerfile,
				".github/workflows/ci.yml": ciWorkflow,
				"go.mod":                   goMod,
				"cmd/app/main.go":          mainGo,
				"config.yaml":              "port: 8080\n",
				"config.json":              "{\"port\": 8080}\n",
			}

			for fileName, content := range files {
				if err := writeFileInDir(output, fileName, content); err != nil {
					return err
				}
			}
			for fileName, content := range coreModuleFiles {
				if err := writeFileInDir(output, fileName, content); err != nil {
					return err
				}
			}

			_, err := fmt.Fprintf(cmd.OutOrStdout(), "scaffolded %s\n", output)
			return err
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "project name for the scaffold")
	cmd.Flags().StringVar(&output, "output", "", "directory where scaffold files will be created")
	cmd.Flags().StringArrayVar(&languages, "language", nil, "target language for code generation (go, typescript, python, java, rust)")
	return cmd
}

func newExportCommand() *cobra.Command {
	var configPath, input, outputDir string
	var languages []string
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export generated artifacts to a directory",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if configPath == "" {
				return fmt.Errorf("missing required --config")
			}
			if input == "" {
				return fmt.Errorf("missing required --input")
			}

			cfg, err := pipeline.ConfigFromFile(configPath)
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}
			if outputDir == "" {
				outputDir = cfg.OutputDir
			}
			if outputDir == "" {
				return fmt.Errorf("missing output directory")
			}
			if len(languages) > 0 {
				cfg.Languages = languages
			}

			specInput, err := resolveInput(input)
			if err != nil {
				return fmt.Errorf("resolve input: %w", err)
			}

			p, err := pipeline.New(cfg)
			if err != nil {
				return fmt.Errorf("failed to construct pipeline: %w", err)
			}

			result, err := p.Run(specInput)
			if err != nil {
				return fmt.Errorf("export run failed: %w", err)
			}

			if dryRun {
				for _, artifact := range result.Artifacts {
					fmt.Fprintf(cmd.OutOrStdout(), "  %s (%d bytes)\n", artifact.Path, len(artifact.Content))
				}
				_, err = fmt.Fprintf(cmd.OutOrStdout(), "\ndry-run: %d artifacts would be generated\n", len(result.Artifacts))
				return err
			}

			if err := os.MkdirAll(outputDir, 0o755); err != nil {
				return fmt.Errorf("create export dir: %w", err)
			}
			for _, artifact := range result.Artifacts {
				artifactPath := filepath.Join(outputDir, artifact.Path)
				if err := os.MkdirAll(filepath.Dir(artifactPath), 0o755); err != nil {
					return fmt.Errorf("create artifact dir: %w", err)
				}
				if err := os.WriteFile(artifactPath, artifact.Content, 0o600); err != nil {
					return fmt.Errorf("write artifact %s: %w", artifact.Path, err)
				}
			}

			_, err = fmt.Fprintf(cmd.OutOrStdout(), "exported %d artifacts to %s\n", len(result.Artifacts), outputDir)
			if len(languages) > 0 {
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "languages: %s\n", strings.Join(languages, ", "))
			}
			return err
		},
	}

	cmd.Flags().StringVar(&configPath, "config", "", "path to JSON or YAML config file")
	cmd.Flags().StringVar(&input, "input", "", "specification input or file path to process")
	cmd.Flags().StringVar(&outputDir, "output-dir", "", "directory to write exported artifacts")
	cmd.Flags().StringArrayVar(&languages, "language", nil, "target language for code generation (go, typescript, python, java, rust)")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "preview artifacts without writing to disk")
	return cmd
}

func newPreviewCommand() *cobra.Command {
	var configPath, input string

	cmd := &cobra.Command{
		Use:   "preview",
		Short: "Preview generated artifacts without writing them",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if configPath == "" {
				return fmt.Errorf("missing required --config")
			}
			if input == "" {
				return fmt.Errorf("missing required --input")
			}

			cfg, err := pipeline.ConfigFromFile(configPath)
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			specInput, err := resolveInput(input)
			if err != nil {
				return fmt.Errorf("resolve input: %w", err)
			}

			p, err := pipeline.New(cfg)
			if err != nil {
				return fmt.Errorf("failed to construct pipeline: %w", err)
			}

			result, err := p.Run(specInput)
			if err != nil {
				return fmt.Errorf("preview run failed: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "preview for %s\n", cfg.Name)
			for _, artifact := range result.Artifacts {
				fmt.Fprintf(cmd.OutOrStdout(), "- %s\n", artifact.Path)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&configPath, "config", "", "path to JSON or YAML config file")
	cmd.Flags().StringVar(&input, "input", "", "specification input or file path to process")
	return cmd
}

func newKernelCommand() *cobra.Command {
	var configPath, outputFormat, topic, payload string

	cmd := &cobra.Command{
		Use:   "kernel",
		Short: "Inspect the NAEOS kernel and service registry",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "services",
		Short: "List registered kernel services",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if configPath == "" {
				return fmt.Errorf("missing required --config")
			}

			cfg, err := pipeline.ConfigFromFile(configPath)
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			p, err := pipeline.New(cfg)
			if err != nil {
				return fmt.Errorf("failed to construct pipeline: %w", err)
			}

			services := p.RegisteredKernelServices()
			rendered, err := renderOutput(services, outputFormat, func() []byte {
				return []byte(strings.Join(services, "\n") + "\n")
			})
			if err != nil {
				return err
			}

			_, err = cmd.OutOrStdout().Write(rendered)
			return err
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "metrics",
		Short: "Show kernel telemetry metrics",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if configPath == "" {
				return fmt.Errorf("missing required --config")
			}

			cfg, err := pipeline.ConfigFromFile(configPath)
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			p, err := pipeline.New(cfg)
			if err != nil {
				return fmt.Errorf("failed to construct pipeline: %w", err)
			}

			metrics := p.KernelMetrics()
			rendered, err := renderOutput(metrics, outputFormat, func() []byte {
				return []byte(fmt.Sprintf("events=%d\nlast_event=%s\n", metrics.Events, metrics.LastEvent.Name))
			})
			if err != nil {
				return err
			}

			_, err = cmd.OutOrStdout().Write(rendered)
			return err
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "events",
		Short: "List active kernel event topics",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if configPath == "" {
				return fmt.Errorf("missing required --config")
			}

			cfg, err := pipeline.ConfigFromFile(configPath)
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			p, err := pipeline.New(cfg)
			if err != nil {
				return fmt.Errorf("failed to construct pipeline: %w", err)
			}

			topics := p.KernelTopics()
			rendered, err := renderOutput(topics, outputFormat, func() []byte {
				return []byte(strings.Join(topics, "\n") + "\n")
			})
			if err != nil {
				return err
			}

			_, err = cmd.OutOrStdout().Write(rendered)
			return err
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "publish",
		Short: "Publish an event to the kernel event bus",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if configPath == "" {
				return fmt.Errorf("missing required --config")
			}
			if topic == "" {
				return fmt.Errorf("missing required --topic")
			}

			cfg, err := pipeline.ConfigFromFile(configPath)
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			p, err := pipeline.New(cfg)
			if err != nil {
				return fmt.Errorf("failed to construct pipeline: %w", err)
			}

			if err := p.Publish(topic, payload); err != nil {
				return err
			}
			_, err = fmt.Fprintf(cmd.OutOrStdout(), "published topic=%s payload=%v\n", topic, payload)
			return err
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "subscribe",
		Short: "Subscribe to a kernel event topic and optionally publish a sample payload",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if configPath == "" {
				return fmt.Errorf("missing required --config")
			}
			if topic == "" {
				return fmt.Errorf("missing required --topic")
			}

			cfg, err := pipeline.ConfigFromFile(configPath)
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			p, err := pipeline.New(cfg)
			if err != nil {
				return fmt.Errorf("failed to construct pipeline: %w", err)
			}

			var received any
			if err := p.Subscribe(topic, func(v any) {
				received = v
			}); err != nil {
				return err
			}

			if payload != "" {
				if err := p.Publish(topic, payload); err != nil {
					return err
				}
			}

			rendered, err := renderOutput(map[string]any{"topic": topic, "received": received}, outputFormat, func() []byte {
				return []byte(fmt.Sprintf("topic=%s received=%v\n", topic, received))
			})
			if err != nil {
				return err
			}

			_, err = cmd.OutOrStdout().Write(rendered)
			return err
		},
	})

	cmd.PersistentFlags().StringVar(&configPath, "config", "", "path to JSON or YAML config file")
	cmd.PersistentFlags().StringVar(&outputFormat, "output", "text", "output format: text, json, or yaml")
	cmd.PersistentFlags().StringVar(&topic, "topic", "", "kernel event topic")
	cmd.PersistentFlags().StringVar(&payload, "payload", "", "event payload to publish")
	return cmd
}

func loadInput(input, inputFile string) (string, error) {
	if input == "" && inputFile == "" {
		return "", fmt.Errorf("missing required --input or --input-file")
	}
	if input != "" && inputFile != "" {
		return "", fmt.Errorf("cannot use both --input and --input-file")
	}
	if inputFile != "" {
		data, err := os.ReadFile(inputFile)
		if err != nil {
			return "", fmt.Errorf("read input file: %w", err)
		}
		return string(data), nil
	}
	return input, nil
}

func resolveInput(input string) (string, error) {
	if input == "" {
		return "", nil
	}

	info, err := os.Stat(input)
	if err != nil {
		if os.IsNotExist(err) {
			return input, nil
		}
		return "", err
	}
	if info.IsDir() {
		return input, nil
	}

	content, err := os.ReadFile(input)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func renderOutput(data any, format string, defaultFn func() []byte) ([]byte, error) {
	switch strings.ToLower(format) {
	case "json":
		result, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("encode json output: %w", err)
		}
		return append(result, '\n'), nil
	case "yaml":
		result, err := yaml.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("encode yaml output: %w", err)
		}
		return result, nil
	default:
		return defaultFn(), nil
	}
}

func writeOrPrint(cmd *cobra.Command, data []byte, filePath string) error {
	if filePath != "" {
		return os.WriteFile(filePath, data, 0o600)
	}
	_, err := cmd.OutOrStdout().Write(data)
	return err
}

func writeFileInDir(dir, fileName, content string) error {
	path := filepath.Join(dir, fileName)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create scaffold dir for %s: %w", fileName, err)
	}
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		return fmt.Errorf("write scaffold %s: %w", fileName, err)
	}
	return nil
}
