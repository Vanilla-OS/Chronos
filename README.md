<div align="center">
    <img src="chronos.svg" height="128">
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

## API Reference

### Get Status

Get the status of the server.

- **URL**: `http://localhost:8080/`
- **Method**: GET
- **Response**:
```json
{
  "status": "ok"
}
```

### Get Articles

Get a list of articles, grouped by language.

- **URL**: `http://localhost:8080/articles`
- **Method**: GET
- **Response**:
```json
{
  "title": "Chronos",
  "SupportedLang": ["en", "it"],
  "groupedArticles": {
    "en": {
      "articles_repo/articles/en/test.md": {
        "Title": "Test Article",
        "Description": "This is a test article written in English.",
        "PublicationDate": "2023-06-10",
        "Authors": ["mirkobrombin"],
        "Body": "..."
      }
    },
    "it": {
      "articles_repo/articles/it/test.md": {
        "Title": "Articolo di Test",
        "Description": "Questo è un articolo di test scritto in italiano.",
        "PublicationDate": "2023-06-10",
        "Authors": ["mirkobrombin"],
        "Body": "..."
      }
    }
  }
}
```

### Get Supported Languages

Get a list of supported languages.

- **URL**: `http://localhost:8080/langs`
- **Method**: GET
- **Response**:
```json
{
  "SupportedLang": ["en", "it"]
}
```

### Get Article by Language and Slug

Get a specific article by providing its language and slug.

- **URL**: `http://localhost:8080/articles/en/test`
  (or `http://localhost:8080/articles/test` based on browser language)
- **Method**: GET
- **Response**:
```json
{
  "Title": "Test Article",
  "Description": "This is a test article written in English.",
  "PublicationDate": "2023-06-10",
  "Authors": ["mirkobrombin"],
  "Body": "..."
}
```

### Search Articles

Search articles based on a query string.

- **URL**: `http://localhost:8080/search?q=test`
- **Method**: GET
- **Response**:
```json
{
  "query": "test",
  "results": [
    {
      "Title": "Test Article",
      "Description": "This is a test article written in English.",
      "PublicationDate": "2023-06-10",
      "Authors": ["mirkobrombin"],
      "Body": "..."
    }
  ]
}
```
