<div align="center">
    <img src="chronos.svg" height="64">
    <h1>Chronos</h1>
    <p>A simple, fast and lightweight documentation server written in Go.</p>
    <hr />
</div>


## Features
- API based response
- Markdown support
- Search
- Localized
- Git based articles storage

## Configuration
Edit the `config/chronos.json` file to configure the server.

## Run (dev)
```bash
go run main.go
```

## Build
```bash
go build -o chronos main.go
```

## Run (prod)
```bash
./chronos
```
## Write articles
Articles are stored in the `articles` folder. The folder structure is the following:
```
articles
├── en
│   └── article1.md
└── it
    └── article1.md
```

You can use a Git repository to store your articles. Just set the `gitrepo` property in the `config/chronos.json` file and ensure that the `articles` folder is present in the repository.
