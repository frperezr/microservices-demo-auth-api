# Auth Api

## Method

## Login

### Request

```
{
  email: string,
  password: string
}
```

### Response

```
{
    token: string,
    error: object {
        code: int,
        message: string
    }
}
```

## Signup

### Request

```
{
  data: object {
    email: string,
    name: string,
    last_name: string,
    password: string
  }
}
```

### Response

```
{
    token: string,
    error: object {
        code: int,
        message: string
    }
}
```
