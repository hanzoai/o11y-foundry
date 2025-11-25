package signoz

// Schema version must follow semantic versioning
#SchemaVersion: =~"^v[0-9]+$"

// Platform enumeration - add more platforms as needed
#Platform: "docker" | "linux" | "kubernetes" | "aws" | "gcp" | "azure" | "windows"

// Deployment type enumeration
#DeploymentType: "standard" | "highly-available"

// Environment variable key-value pair
#EnvVar: {
	key:   string
	value: string
}

// Component definition
#Component: {
	name:     string
	replicas: int & >=1
	version:  string & =~"^[0-9]+\\.[0-9]+(\\.[0-9]+)?(-.*)?$"
	env?: [...#EnvVar]
}

// Platform-specific requirements
#Requirements: {
	docker: ["docker", "docker-compose"]
	linux: ["systemd", "curl", "tar"]
	kubernetes: ["kubectl", "helm"]
	aws: ["aws-cli", "kubectl", "helm"]
	gcp: ["gcloud", "kubectl", "helm"]
	azure: ["az-cli", "kubectl", "helm"]
	windows: ["powershell", "chocolatey"]
}

// Main configuration schema
#Config: {
	schemaVersion: #SchemaVersion
	platform:      #Platform
	type:          #DeploymentType

	// Components involved in the deployment
	components: [...#Component]

	// Requirements based on platform
	requirements: #Requirements[platform]
}

// Validate that the config conforms to the schema
#Config
