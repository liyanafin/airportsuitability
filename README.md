# AirportSuitability API

AirportSuitability is a Go project designed to evaluate the suitability of airports for specific aircraft, considering factors such as runways, weather, and NOTAMs. This project is developed using Test-Driven Development (TDD) principles.

## Project Foundation

- Modular Go project structure
- Initial domain model for Airport and Runway
- First TDD test for Airport struct
- Ready for incremental, test-driven feature development

## Directory Structure

```
internal/
  domain/
    airport.go         # Airport and Runway domain models
    airport_test.go    # TDD: Basic Airport struct tests
main.go               # Entry point (empty for test-only foundation)
```

## Getting Started

### Prerequisites
- Go 1.18 or later

### Running Tests

To run the initial test suite:

```sh
go test ./internal/domain
```

## Next Steps
- Expand domain models (aircraft, weather, NOTAM, suitability)
- Add repositories, services, and API layers
- Continue TDD for each new feature

## License
MIT License
