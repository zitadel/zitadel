name: Bug Report
description: Create a report to help us improve ZITADEL
title: "[Bug]: "
labels: ["type: bug", "state: triage"]
assignees:
body:
- type: markdown
  attributes:
  value: |
  Thanks for taking the time to fill out this bug report!
- type: checkboxes
  id: preflight
  attributes:
  label: Preflight Checklist
  options:
    - label:
      I could not find a solution in the documentation, the existing issues or discussions
      required: true
    - label:
      I have joined the [ZITADEL chat](https://zitadel.com/chat)
- type: dropdown
  id: environment
  attributes:
  label: Environment
  description: How do you use ZITADEL?
  options:
  - ZITADEL Cloud
  - Selfhosted
    required: true
- type: textarea
  id: description
  attributes:
  label: Describe the bug
  description: A clear and concise description of what the bug is.
  required: true
- type: textarea
  id: reproduce
  attributes:
  label: To reproduce
  description: Steps to reproduce the behaviour
  required: true
- type: textarea
  id: screenshots
  label: Screenshots
  placeholder: If applicable, add screenshots to help explain your problem.
- type: textarea
  id: expected
  label: Expected behavior
  placeholder: A clear and concise description of what you expected to happen.
- type: textarea
  id: version
  label: Version
  description: Which version of ZITADEL are you using.
- type: textarea
  id: os
  label: Operating System
  description: Please complete informations about your operating-system, device, browser, etc.
- type: textarea
  id: additional
  label: Additional Context
  description: Please add any other infos that could be useful.