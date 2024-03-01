# elasticsearch-golang-docker-compose-example
An example docker-compose setup for connecting to Elasticsearch with Go

## Usage

1. Install Docker
2. Clone the repository
3. Run `docker-compose up`

## Expected Output

```
app-1            | 2024/03/01 21:34:01 Document indexed successfully: This is the first document.
app-1            | 2024/03/01 21:34:01 Document indexed successfully: This is the second document.
app-1            | 2024/03/01 21:34:01 Document indexed successfully: This is the third document.
app-1            | Matched document: This is the first document.
app-1            | Matched document: This is the second document.
app-1            | Matched document: This is the third document.
app-1 exited with code 0
```
