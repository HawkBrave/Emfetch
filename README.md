# Emfetch

> Query dummy data without any effort

## About

## Documentation

Start with typing `<url>/`.

### Create a simple query

Add `{["people":{"name":"John","surname":"Smith"}*100]}`.
The resulted response should be the following JSON:

```
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
