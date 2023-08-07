# Workflows

| Workflow           | Description                                                    |
|--------------------|----------------------------------------------------------------|
| ci-build.yaml      | Build, lint, test, codegen, build-ui, analyze, e2e-test        |
| image-reuse.yaml   | Build, push, and Sign container images                         |
| image.yaml         | Build container image for PR's & publish for push events       |
| init-release.yaml  | Build manifests and version then create a PR for release branch|
| release.yaml       | Build images, cli-binaries, provenances, and post actions      |
| update-snyk.yaml   | Scheduled snyk reports                                         |

Note: This actions are heavily borrowed from ArgoCD repository.