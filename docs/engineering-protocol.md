# Engineering Protocol (Mandatory)

This file defines how we build and change DataTracing v2.

## 1) TDD Enforcement Protocol (always required)

For **every code change**, follow this cycle strictly:

1. **RED**: Write failing tests that expose design flaws.
2. **GREEN**: Write minimal code to make tests pass.
3. **REFACTOR**: Improve structure while keeping tests green.
4. **VALIDATE**: Run the full relevant test suite to ensure correctness.
5. **REPEAT**: Never skip the cycle.

### Non-negotiable AI testing rule
When testing AI-integrated behavior, **do not mock the AI provider**. Use the real provider in tests so behavior reflects real-world outputs and failure modes.

## 2) SOLID principles (always apply)

- **Single Responsibility Principle**: each class/module/function has one reason to change.
- **Open-Closed Principle**: extend behavior via composition/abstractions, avoid risky edits to stable code.
- **Liskov Substitution Principle**: implementations behind interfaces must remain substitutable.
- **Interface Segregation Principle**: keep interfaces small and consumer-focused.
- **Dependency Inversion Principle**: depend on abstractions, not concrete infrastructure.

## 3) Clean architecture rules

- Define boundaries with interfaces/ports.
- Keep business logic in application/domain layers, isolated from infrastructure concerns.
- Use dependency injection for wiring.
- Optimize for high cohesion and low coupling.
- Use tests to enforce architecture boundaries.

## 4) Coding standards

- Use descriptive names for functions, variables, and tests.
- Keep functions small and focused.
- Avoid unnecessary side effects.
- Use exceptions/errors for error handling (not magic return codes).
- Add meaningful comments only when intent is non-obvious.

## 5) Delivery workflow rules

- Use Git branches for features and bug fixes.
- Use descriptive commit messages.
- Keep all code changes version-controlled.
- Use pull requests for review and integration.

## 6) YAGNI policy

- Do not build speculative features.
- Implement only what current requirements demand.
- Prefer simple designs over premature abstraction.
- Refactor when needed by concrete change pressure.
- Write tests for current behavior, not hypothetical future behavior.

## 7) Definition of done for each task

A task is done only when all are true:

- Tests were written first (RED), then implementation (GREEN).
- Refactoring completed without breaking tests.
- Relevant test suites pass locally.
- Architecture boundaries remain intact.
- Changes are committed with a clear message.
