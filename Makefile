.PHONY: help build generate serve run clean

# Default target
help:
	@echo "Available targets:"
	@echo "  make build      - Build the static site generator binary"
	@echo "  make generate   - Generate the static site in public/"
	@echo "  make serve      - Serve the generated site locally on port 8000"
	@echo "  make run        - Build the generator and generate the site"
	@echo "  make clean      - Remove generated files (ssg binary and public/ directory)"
	@echo "  make help       - Show this help message"

# Build the static site generator binary
build:
	@echo "Building static site generator..."
	go build -o ssg

# Generate the static site
generate: build
	@echo "Generating static site..."
	./ssg

# Serve the generated site locally
serve:
	@echo "Serving site at http://localhost:8000"
	@cd public && python3 -m http.server 8000

# Build and generate in one step
run: generate
	@echo "Site generated successfully in public/"

# Clean generated files
clean:
	@echo "Cleaning generated files..."
	rm -f ssg
	rm -rf public/
