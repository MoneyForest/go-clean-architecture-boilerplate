# go-clean-architecture-boilerplate

## Architecture Overview

This boilerplate is built on the Clean Architecture.
It puts the business logic at the center and protects it from the details of external implementations, resulting in testable and maintainable code.

![clean-architecture-circles](docs/clean-arch-circles.png)

### Layer Structure

```mermaid
%%{init: {"flowchart": {"htmlLabels": true}} }%%
graph LR
    subgraph External["External Interfaces"]
        HTTP["HTTP"]
        DB["MySQL"]
        CACHE["Redis"]
        SQS["AWS SQS"]
    end

    subgraph Infrastructure["Infrastructure Layer"]
        subgraph Controller["Controller"]
            H["Handler"]
        end
        subgraph Gateway["Gateway"]
            MR["MySQL Repository"]
            RR["Redis Repository"]
            SR["SQS Repository"]
        end
    end

    subgraph Usecase["Usecase Layer"]
        I["Interactor"]
    end

    subgraph Domain["Domain Layer"]
        DM["Domain Model"]
        DS["Domain Service"]
        RI["Repository Interface"]
    end

    %% External → Infrastructure
    HTTP --> H
    MR --> DB
    RR --> CACHE
    SR --> SQS

    %% Controller → Usecase
    H --> I

    %% Usecase → Domain
    I --> DM
    I --> DS
    I --> RI
    DS --> DM
    DS --> RI

    %% Repository Implementations
    Gateway --> RI

    classDef domain fill:#e8372c,stroke:#333,stroke-width:2px,color:#fff
    classDef usecase fill:#f58634,stroke:#333,stroke-width:2px,color:#fff
    classDef infra fill:#58aa46,stroke:#333,stroke-width:2px,color:#fff
    classDef external fill:#2c4a9b,stroke:#333,stroke-width:2px,color:#fff
```

### Dependency Rules

1. Dependencies always point inwards
   - Infrastructure → Usecase → Domain
   - Outer layers know about inner layers
   - Inner layers don't know about outer layers

2. Dependency Inversion Principle
   - Repository interfaces are defined in the domain layer
   - Implementations are provided in the infrastructure layer

### Layer Mapping to Directory

1. **Domain Layer** (`internal/domain`) - Innermost layer
   - Defines pure business logic independent of external implementations
   - `domain/model`: Core business entities and rules
   - `domain/repository`: Interfaces defining operations on domain models

2. **Usecase Layer** (`internal/usecase`) - Middle layer
   - Implements application-specific use cases
   - Ensuring the integrity of the data store, such as DB transaction management, is also the responsibility of this layer.
   - `usecase/interactor`: Implements business logic that operates on domain models
   - `usecase/port`: Defines input and output ports

3. **Infrastructure Layer** (`internal/infrastructure`) - Outermost layer
   - Concrete implementations of interfaces with external systems, frameworks, and databases
   - `infrastructure/cmd`: Implements routing from entry points to controllers based on commands
   - `infrastructure/controller`: Implements HTTP servers, gRPC servers, commands, and subscribers
   - `infrastructure/gateway`: Connects to external systems (MySQL, Redis)

## Directory Structure Overview

```sh
├── internal
│   ├── domain
│   │   ├── model
│   │   └── repository
│   │
│   ├── usecase
│   │   ├── interactor
│   │   └── port
│   │       ├── input
│   │       └── output
│   │
│   ├── infrastructure
│   │   ├── cmd
│   │   │   ├── http
│   │   │   ├── subscriber
│   │   │   └── task
│   │   │
│   │   ├── controller
│   │   │   ├── http
│   │   │   ├── subscriber
│   │   │   └── task
│   │   │
│   │   ├── gateway
│   │   │   ├── mysql
│   │   │   └── redis
│   │   │
│   │   └── environment
│   │
│   └── dependency
│
└── pkg
```

## Technology Stack

| Category | Component | Details |
|----------|-----------|----------|
| **Backend** | Architecture | Clean Architecture |
| | Language | Go |
| | Framework | • Chi (HTTP router and dispatcher)<br>• Cobra (Command-line framework) |
| | Documentation | OpenAPI (API documentation) |
| **Database & Caching** | PDatabase | MySQL 8.0 |
| | Cache | Redis |
| | ORM | database/sql |
| **Infrastructure & Cloud** | Container | Docker |
| | Message Queue | AWS SQS |
| | CI/CD | GitHub Actions |
| **Tools** | Documentation Generator | Swag (Generate OpenAPI documentation) |
