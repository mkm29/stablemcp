# Changelog

All notable changes to the StableMCP project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.1] - 2025-03-29

### Added
- Comprehensive configuration system with proper defaults
- OpenAI integration with configurable API key, model, and base URL
- Structured logging with configurable level and format
- Telemetry support for metrics and distributed tracing
- TLS configuration options for secure server connections
- Download path option for storing generated images

### Changed
- Refactored configuration loading to use deterministic file paths
- Updated configuration to use `.stablemcp.yaml` or `.stablemcp.json` files
- Improved configurability with nested configuration options
- Enhanced default configuration values for better out-of-box experience

### Fixed
- Fixed configuration loading to properly handle default values
- Addressed potential issues with configuration file search paths

## [0.1.0] - 2025-03-28

### Added
- Initial project structure and scaffolding
- Basic server command implementation
- Configuration loading from multiple locations
- Preliminary Model Context Protocol support

[0.1.1]: https://github.com/modelcontextprotocol/stablemcp/compare/v0.1.0...v0.1.1
[0.1.0]: https://github.com/modelcontextprotocol/stablemcp/releases/tag/v0.1.0