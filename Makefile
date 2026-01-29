.PHONY: help build generate serve run clean new-post

# Default target
help:
	@echo "Available targets:"
	@echo "  make build      - Build the static site generator binary"
	@echo "  make generate   - Generate the static site in public/"
	@echo "  make serve      - Serve the generated site locally on port 8000"
	@echo "  make run        - Build the generator and generate the site"
	@echo "  make clean      - Remove generated files (ssg binary and public/ directory)"
	@echo "  make new-post SLUG=my-post - Create a new draft blog post"
	@echo "  make help       - Show this help message"

# Build the static site generator binary
build:
	@echo "Building static site generator..."
	go build -o ssg

# Generate the static site
generate: build
	@echo "Generating static site..."
	./ssg

# Serve the generated site locally (includes drafts for development)
serve: build
	@echo "Generating static site (including drafts)..."
	./ssg -include-drafts -skip-rss-if-exists
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

# Create a new draft blog post
new-post:
	@if [ -z "$(SLUG)" ]; then \
		echo "Error: SLUG is required. Usage: make new-post SLUG=my-post-slug"; \
		exit 1; \
	fi
	@if [ -f "content/blog/$(SLUG).md" ]; then \
		echo "Error: File content/blog/$(SLUG).md already exists"; \
		exit 1; \
	fi
	@echo "Creating new draft post: content/blog/$(SLUG).md"
	@printf -- "---\ntitle: $(SLUG)\nslug: $(SLUG)\ndate: $$(date +%Y-%m-%d)\ndraft: true\n---\n\nYour content here.\n" > content/blog/$(SLUG).md
	@echo "Created! Edit content/blog/$(SLUG).md to write your post."
