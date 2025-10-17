# Update Templates Guide

## Prerequisites

- `jq` - JSON processor
- `Go` - For building and generating docs

## Update Process

### 1. Run Update Commands

```bash
make templates-delete
make templates-get
make templates-generate
```

### 2. Check Git Changes

```bash
git status
```

- **Untracked directories** → New services added
- **Modified files** → Updated configurations
- **Deleted directories** → Removed services

### 3. Handle Removed Services

**Important:** Services cannot be removed in minor/patch versions. They must be deprecated first.

If you see deleted directories:

1. Find the template ID by checking the old `internal/provider/templates.json`

2. Add deprecation entry in `internal/provider/templates.go` → `getDeprecatedResourcesConfig()`:

   ```go
   {
       TemplateId:         <id>,
       ResourceName:       "<resource_name>",
       DocumentationName:  "<Display Name>",
       DeprecationMessage: "This resource is no more supported by Elestio and will be removed in the next major version of the provider.",
   }
   ```

3. Commit both the deletion and the deprecation entry

### 4. Generate Documentation

```bash
make generate
```

### 5. Commit Changes

```bash
git add .
git commit -m "Update templates: add X, update Y, deprecate Z"
```

## Quick Reference

```bash
make templates-delete
make templates-get
make templates-generate
git status
# If services removed: add deprecation in templates.go
make generate
git add .
git commit -m "Update templates: ..."
```
