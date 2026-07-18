package cicd

import (
	"fmt"
	"strings"
)

type AzurePipelinesGenerator struct{}

func (g *AzurePipelinesGenerator) Name() string {
	return "Azure Pipelines"
}

func (g *AzurePipelinesGenerator) Generate(config *PipelineConfig) (string, error) {
	var sb strings.Builder

	sb.WriteString("trigger:\n")
	if config.Trigger.OnPush {
		sb.WriteString("  branches:\n")
		sb.WriteString("    include:\n")
		sb.WriteString("      - main\n")
		sb.WriteString("      - develop\n")
	}
	if config.Trigger.OnPR {
		sb.WriteString("pr:\n")
		sb.WriteString("  branches:\n")
		sb.WriteString("    include:\n")
		sb.WriteString("      - main\n")
	}
	if config.Trigger.OnRelease {
		sb.WriteString("trigger:\n")
		sb.WriteString("  tags:\n")
		sb.WriteString("    include:\n")
		sb.WriteString("      - v*\n")
	}
	if config.Trigger.Schedule != "" {
		sb.WriteString("schedules:\n")
		sb.WriteString(fmt.Sprintf("  - cron: '%s'\n", config.Trigger.Schedule))
		sb.WriteString("    displayName: 'Scheduled Build'\n")
		sb.WriteString("    branches:\n")
		sb.WriteString("      include:\n")
		sb.WriteString("        - main\n")
	}
	sb.WriteString("\n")

	sb.WriteString("pool:\n")
	sb.WriteString("  vmImage: 'ubuntu-latest'\n\n")

	if len(config.Secrets) > 0 {
		sb.WriteString("variables:\n")
		for _, secret := range config.Secrets {
			sb.WriteString(fmt.Sprintf("  %s: $(%s)\n", strings.ToUpper(secret), secret))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("stages:\n")

	// Build stage
	sb.WriteString("  - stage: Build\n")
	sb.WriteString("    jobs:\n")
	sb.WriteString("      - job: BuildJob\n")
	sb.WriteString("        steps:\n")
	for _, lang := range config.Languages {
		switch lang {
		case "go":
			sb.WriteString("          - task: GoTool@0\n")
			sb.WriteString("            inputs:\n")
			sb.WriteString("              version: '1.22'\n")
			sb.WriteString("          - script: go build ./...\n")
			sb.WriteString("            displayName: 'Build'\n")
		case "node", "typescript":
			sb.WriteString("          - task: NodeTool@0\n")
			sb.WriteString("            inputs:\n")
			sb.WriteString("              versionSpec: '20.x'\n")
			sb.WriteString("          - script: npm ci && npm run build\n")
			sb.WriteString("            displayName: 'Build'\n")
		case "python":
			sb.WriteString("          - task: UsePythonVersion@0\n")
			sb.WriteString("            inputs:\n")
			sb.WriteString("              versionSpec: '3.12'\n")
			sb.WriteString("          - script: pip install -r requirements.txt\n")
			sb.WriteString("            displayName: 'Install'\n")
		case "java":
			sb.WriteString("          - task: JavaToolInstaller@0\n")
			sb.WriteString("            inputs:\n")
			sb.WriteString("              versionSpec: '21'\n")
			sb.WriteString("              jdkArchitectureOption: 'x64'\n")
			sb.WriteString("              jdkSourceOption: 'PreInstalled'\n")
			sb.WriteString("          - script: mvn clean compile\n")
			sb.WriteString("            displayName: 'Build'\n")
		case "rust":
			sb.WriteString("          - script: |\n")
			sb.WriteString("              curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y\n")
			sb.WriteString("              source $HOME/.cargo/env\n")
			sb.WriteString("              cargo build --release\n")
			sb.WriteString("            displayName: 'Build'\n")
		}
	}

	// Test stage
	sb.WriteString("  - stage: Test\n")
	sb.WriteString("    dependsOn: Build\n")
	sb.WriteString("    jobs:\n")
	sb.WriteString("      - job: TestJob\n")
	sb.WriteString("        steps:\n")
	for _, lang := range config.Languages {
		switch lang {
		case "go":
			sb.WriteString("          - task: GoTool@0\n")
			sb.WriteString("            inputs:\n")
			sb.WriteString("              version: '1.22'\n")
			sb.WriteString("          - script: go test ./...\n")
			sb.WriteString("            displayName: 'Test'\n")
		case "node", "typescript":
			sb.WriteString("          - task: NodeTool@0\n")
			sb.WriteString("            inputs:\n")
			sb.WriteString("              versionSpec: '20.x'\n")
			sb.WriteString("          - script: npm ci && npm test\n")
			sb.WriteString("            displayName: 'Test'\n")
		case "python":
			sb.WriteString("          - task: UsePythonVersion@0\n")
			sb.WriteString("            inputs:\n")
			sb.WriteString("              versionSpec: '3.12'\n")
			sb.WriteString("          - script: pytest\n")
			sb.WriteString("            displayName: 'Test'\n")
		case "java":
			sb.WriteString("          - task: JavaToolInstaller@0\n")
			sb.WriteString("            inputs:\n")
			sb.WriteString("              versionSpec: '21'\n")
			sb.WriteString("              jdkArchitectureOption: 'x64'\n")
			sb.WriteString("              jdkSourceOption: 'PreInstalled'\n")
			sb.WriteString("          - script: mvn test\n")
			sb.WriteString("            displayName: 'Test'\n")
		case "rust":
			sb.WriteString("          - script: |\n")
			sb.WriteString("              source $HOME/.cargo/env\n")
			sb.WriteString("              cargo test\n")
			sb.WriteString("            displayName: 'Test'\n")
		}
	}

	// Deploy stage
	sb.WriteString("  - stage: Deploy\n")
	sb.WriteString("    dependsOn: Test\n")
	sb.WriteString("    condition: and(succeeded(), eq(variables['Build.SourceBranch'], 'refs/heads/main'))\n")
	sb.WriteString("    jobs:\n")
	sb.WriteString("      - deployment: DeployJob\n")
	sb.WriteString("        environment: 'production'\n")
	sb.WriteString("        strategy:\n")
	sb.WriteString("          runOnce:\n")
	sb.WriteString("            deploy:\n")
	sb.WriteString("              steps:\n")
	sb.WriteString("                - script: echo 'Deploying...'\n")
	sb.WriteString("                  displayName: 'Deploy'\n")
	for _, step := range config.Steps {
		sb.WriteString(fmt.Sprintf("                - script: %s\n", step.Command))
		sb.WriteString(fmt.Sprintf("                  displayName: '%s'\n", step.Name))
	}

	return sb.String(), nil
}
