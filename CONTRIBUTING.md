# Contributing

Thanks for helping make Pii better!

If you'd like to add new exported APIs, please [open an issue](https://github.com/alifengineer/pii/issues/new) describing your proposal. Discussing API changes ahead of time makes pull request review much smoother.

## Setup

1. Fork and clone the repository.

   ```bash
   gh repo fork --clone alifengineer/pii
   ```

2. Install dependencies:

   ```bash
   go mod download
   ```

3. Verify that tests pass:

   ```bash
   make test
   make lint
   ```

## Making Changes

1. Create a feature branch:

   ```bash
   git checkout -b my-feature
   ```

2. Make your changes and verify tests still pass:

   ```bash
   make test
   make lint
   ```

3. Push and open a pull request:

   ```bash
   git push origin my-feature
   gh pr create
   ```

## Guidelines

- Add tests for new functionality
- Maintain backward compatibility
- Write clear commit messages
- Run `make lint` before submitting
