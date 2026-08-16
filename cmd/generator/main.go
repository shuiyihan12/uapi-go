// Command generator scaffolds Go service packages from the manifest of
// Travelport service WSDLs, producing the types and call wrappers for each
// service under pkg/services.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// ServiceTemplate is the Go text/template used to scaffold a single service
// client (Service struct, constructor and Close method).
const ServiceTemplate = `package {{.PackageName}}

import (
	"context"
	"encoding/xml"
	"fmt"

	"github.com/shuiyihan12/uapi-go/pkg/client"
	"github.com/shuiyihan12/uapi-go/pkg/logging"
	"go.uber.org/zap"
)

// {{.ServiceName}}Service is the {{.Description}} service client.
type {{.ServiceName}}Service struct {
	client   *client.EnterpriseSOAPClient
	logger   logging.Logger
	endpoint string
}

// New{{.ServiceName}}Service creates a {{.Description}} service client.
func New{{.ServiceName}}Service(config client.SOAPConfig, logger logging.Logger) (*{{.ServiceName}}Service, error) {
	enterpriseConfig := client.EnterpriseConfig{
		SOAPConfig:  config,
		ServiceName: "{{.PackageName}}-service",
		Logger:      logger,
	}

	enterpriseClient, err := client.NewEnterpriseSOAPClient(enterpriseConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create {{.PackageName}} service client: %v", err)
	}

	return &{{.ServiceName}}Service{
		client:   enterpriseClient,
		logger:   logger,
		endpoint: config.Endpoint,
	}, nil
}

// Close closes the service.
func (s *{{.ServiceName}}Service) Close() error {
	return s.client.Close()
}
`

// ServiceInfo describes a UAPI service to scaffold, including its name,
// package name, description and WSDL file path.
type ServiceInfo struct {
	ServiceName string
	PackageName string
	Description string
	WSDLFile    string
}

// services is the built-in manifest listing every known UAPI service and its
// metadata.
var services = []ServiceInfo{
	{"Hotel", "hotel", "hotel", "wsdl/hotel_v55_0/Hotel.wsdl"},
	{"Universal", "universal", "universal record", "wsdl/universal_v55_0/UniversalRecord.wsdl"},
}

// main parses command-line flags and renders ServiceTemplate for the
// services manifest (or the single service named via -service).
func main() {
	var (
		outputDir = flag.String("output", "pkg/services", "Output directory for service packages")
		force     = flag.Bool("force", false, "Force overwrite existing files")
		service   = flag.String("service", "", "Generate specific service only")
	)
	flag.Parse()

	// Parse the template.
	tmpl, err := template.New("service").Parse(ServiceTemplate)
	if err != nil {
		log.Fatalf("Failed to parse template: %v", err)
	}

	// Filter the services.
	var targetServices []ServiceInfo
	if *service != "" {
		for _, svc := range services {
			if strings.EqualFold(svc.PackageName, *service) || strings.EqualFold(svc.ServiceName, *service) {
				targetServices = append(targetServices, svc)
				break
			}
		}
		if len(targetServices) == 0 {
			log.Fatalf("Service '%s' not found", *service)
		}
	} else {
		targetServices = services
	}

	// Generate the service packages.
	for _, svc := range targetServices {
		generateServicePackage(tmpl, *outputDir, svc, *force)
	}

	fmt.Printf("Generated %d service packages successfully!\n", len(targetServices))
}

// generateServicePackage renders the template for one service into
// service.go inside its package directory, skipping when the target file
// already exists and -force is not set.
func generateServicePackage(tmpl *template.Template, outputDir string, svc ServiceInfo, force bool) {
	// Create the service directory.
	serviceDir := filepath.Join(outputDir, svc.PackageName)
	if err := os.MkdirAll(serviceDir, 0755); err != nil {
		log.Printf("Failed to create directory %s: %v", serviceDir, err)
		return
	}

	// Target service file.
	serviceFile := filepath.Join(serviceDir, "service.go")

	// Skip if the file already exists.
	if !force {
		if _, err := os.Stat(serviceFile); err == nil {
			fmt.Printf("Service file %s already exists, skipping (use -force to overwrite)\n", serviceFile)
			return
		}
	}

	// Create the file.
	file, err := os.Create(serviceFile)
	if err != nil {
		log.Printf("Failed to create file %s: %v", serviceFile, err)
		return
	}
	defer file.Close()

	// Execute the template.
	if err := tmpl.Execute(file, svc); err != nil {
		log.Printf("Failed to execute template for %s: %v", svc.ServiceName, err)
		return
	}

	fmt.Printf("Generated service package: %s\n", serviceFile)
}
