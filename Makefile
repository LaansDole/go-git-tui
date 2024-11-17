.PHONY: execute
execute:
	@echo "Searching for all shell script files and changing their permissions to executable..."
	@find . -type f -name "*.sh" -exec chmod +x {} \; -exec echo "Made executable: {}" \;
	@echo "All shell script files have been processed."
