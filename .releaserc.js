module.exports = {
    branches: ["master"],
    plugins: [
        "@semantic-release/commit-analyzer",
        "@semantic-release/release-notes-generator",
        ["@semantic-release/github", {
            "assets": [
                {"path": ".artifacts/documentation-darwin-amd64/documentation-darwin-amd64", "label": "Darwin x86_64"},
                {"path": ".artifacts/documentation-linux-amd64/documentation-linux-amd64", "label": "Linux x86_64"},
                {
                    "path": ".artifacts/documentation-windows-amd64/documentation-windows-amd64",
                    "label": "Windows x86_64"
                },
            ]
        }],
    ]
};