# gcpsecretmanager2env

A cli tool that parses env file, fetch secrets from gcp secret manager, and substitutes them.

## Example

./.env

```env
AN_ENV_FROM_SECRET_MANAGER=projects/<PROJECT_ID>/secrets/<ENV_NAME>/versions/1
```

command

```
gcpsecretmanager2env -output ./.output.env ./.env
```

./.output.env

```
AN_ENV_FROM_SECRET_MANAGER=<VALUE_FROM_SECRET_MANAGEr>
```

## Usage

```
Usage: gcpsecretmanager2env [OPTIONS] <input-file>
Note: <input-file> is a required positional argument.
  -credential string
        gcp credential file. it will see GOOGLE_APPLICATION_CREDENTIALS when it's not set
  -help
        show help
  -output string
        output file
```