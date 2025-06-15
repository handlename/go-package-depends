# Dependencies

## Layers

Upper layers cannot depend on lower layers.

1. Domain layer
  - Implementation of core entities
2. Application layer
  - Business logic using objects from the domain layer
3. Presentation layer
  - UI and presentation logic
4. Infra layer
  - Gateway to the real world

## Packages in layers

Upper packages cannot depend on lower packages.

1. Domain layer
  - domain/entity
  - domain/valueobject
    - domain/service
2. Application layer
  - app/service
    - app/usecase
3. Presentation layer
  - api
  - cli
4. Infra layer
  - infra/database
  - infra/cache
