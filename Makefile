PROJECT_NAME := grafana-dashboards
WORKSPACE := ./tmp

.PHONY: build
build: clean
	@mkdir -p $(WORKSPACE)
	@zip $(WORKSPACE)/$(PROJECT_NAME).alfredworkflow info.plist bin/dashboards

.PHONY:  clean
clean:
	@rm -rf $(WORKSPACE)/*

.PHONY: install
install: build
	@open $(WORKSPACE)/$(PROJECT_NAME).alfredworkflow
