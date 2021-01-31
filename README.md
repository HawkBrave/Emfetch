# Emfetch

## Query dummy data with Emmet-like syntax

## About

## Documentation

Start with typing `<url>/`.

### Create a simple query

Add `{["people":{"name":"John","surname":"Smith"}*100]}`.
The resulted response should be the following JSON:

```JSON
{
  "people": [
    { "name": "John", "surname": "Smith" },
    // .
    // .
    // .
    // 100 rows total
  ]
}
```
