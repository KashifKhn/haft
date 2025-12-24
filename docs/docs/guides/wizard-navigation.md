---
sidebar_position: 1
title: Wizard Navigation
description: Master the interactive TUI wizard
---

# Wizard Navigation

The Haft wizard provides a rich terminal interface for configuring your Spring Boot project. This guide covers all keyboard shortcuts and navigation patterns.

## General Navigation

| Key | Action |
|-----|--------|
| `Enter` | Confirm selection / Move to next step |
| `Esc` | Go back to previous step |
| `Ctrl+C` | Cancel and exit |

## Text Input Steps

Used for: Project Name, Group ID, Artifact ID, Description, Package Name

| Key | Action |
|-----|--------|
| `←` `→` | Move cursor |
| `Backspace` | Delete character |
| `Enter` | Confirm and continue |
| `Esc` | Go back |

### Auto-Generated Values

Some fields are auto-generated based on previous inputs:

- **Package Name** = Group ID + Artifact ID (sanitized)
  - Example: `com.example` + `my-app` → `com.example.myapp`

You can edit the auto-generated value before confirming.

## Single Select Steps

Used for: Java Version, Spring Boot Version, Build Tool, Packaging, Config Format, Git Init

| Key | Action |
|-----|--------|
| `↑` `↓` | Navigate options |
| `Enter` | Select and continue |
| `Esc` | Go back |

## Multi-Select Steps

Used for selecting multiple options (if applicable).

| Key | Action |
|-----|--------|
| `↑` `↓` | Navigate options |
| `Space` | Toggle selection |
| `Enter` | Confirm selections |
| `Esc` | Go back |

## Dependency Picker

The dependency picker has the most features:

### Navigation

| Key | Action |
|-----|--------|
| `↑` `↓` | Navigate dependencies |
| `Space` | Toggle dependency selection |
| `Enter` | Confirm and continue |
| `Esc` | Go back |

### Category Filtering

| Key | Action |
|-----|--------|
| `Tab` | Next category |
| `Shift+Tab` | Previous category |
| `0` | Show all categories |
| `1-9` | Jump to specific category |

### Category Numbers

| Key | Category |
|-----|----------|
| `0` | All |
| `1` | Developer Tools |
| `2` | Web |
| `3` | SQL |
| `4` | NoSQL |
| `5` | Security |
| `6` | Messaging |
| `7` | Cloud |
| `8` | Observability |
| `9` | Testing |

### Search

| Key | Action |
|-----|--------|
| `/` | Enter search mode |
| `Esc` | Exit search mode |
| Type | Filter dependencies |

Search matches against:
- Dependency name
- Dependency description
- Dependency ID

## Tips

### Quick Project Setup

1. Type project name → `Enter`
2. Accept default Group ID → `Enter`
3. Accept default Artifact ID → `Enter`
4. Skip description → `Enter`
5. Accept package name → `Enter`
6. Select Java 21 → `Enter`
7. Select latest Spring Boot → `Enter`
8. Select Maven → `Enter`
9. Select JAR → `Enter`
10. Select YAML → `Enter`
11. Pick dependencies → `Enter`
12. Select Git init → `Enter`

### Efficient Dependency Selection

1. Press `2` to jump to Web category
2. Select `spring-boot-starter-web`
3. Press `3` to jump to SQL
4. Select `spring-boot-starter-data-jpa`
5. Press `1` to jump to Developer Tools
6. Select `lombok`
7. Press `Enter` to confirm

### Going Back

Made a mistake? Press `Esc` to go back to any previous step. Your selections are preserved.

## Accessibility

- All prompts include help text at the bottom
- Selected items are highlighted
- Category bar shows current position
- Search results update in real-time
