# Octopus Deploy Decryptor

This tool allows for decrypting projects exported from [Octopus Deploy](https://octopus.com).

It creates an `export.html` report file, which reconstructs the project and variable structure.
Sensitive/encrypted variable values are decrypted print as plain text in the report.

## Decrypting a Project Export

1. Build the project using `go build`
2. Run the tool within the directory of the unzipped project export:

```sh
octopus-deploy-decryptor --password 'password' --sensitive-only
```

Make sure to provide the correct password used during export.
The `sensitive-only` option will limit the shown values to sensitive variables.

More on how to export projects can be found here: https://octopus.com/docs/projects/export-import

## Decrypting values manually, using the project export Password

Use the following Go code to decrypt the data:

```go
key, err := KeyFromPassword(password)

v, err := DecryptString(key, encryptedValue)
```

## Decrypting values manually, using the Master Key

1. Retrieve the master key from Octopus Deploy: https://octopus.com/docs/security/data-encryption
2. Use the following Go code to decrypt the data:

```go
key, err := KeyFromMasterKey(masterKey)

v, err := DecryptString(key, encryptedValue)
```

More on how to get the variable values from the database can be found here: https://squirrelistic.com/blog/how_to_display_value_of_sensitive_variable_in_octopus_deploy

This blog article also shows how to use Powershell, if you prefer that.

### Contributors

- [Felix Kollmann](https://github.com/fkollmann)

## License

Published under the [MIT License](./LICENSE).

It is provided as is, in a "works ony my machine" fashion.
