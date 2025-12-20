Template: When creating a release:
1.  **Version:** specific Semantic Versioning (vX.Y.Z).
2.  **Changelog:** Update `CHANGELOG.md` moving `[Unreleased]` items to the new version header. Add date.
3.  **Tag:** Create a git tag `git tag -a vX.Y.Z -m "Release vX.Y.Z"`.
4.  **Push:** `git push origin vX.Y.Z`.
5.  **Verify:** Check CI pipeline for release artifacts.
