# Terraform provider for AuthFed

This Terraform provider is based on [terraform-provider-http](https://github.com/terraform-providers/terraform-provider-http).

## Getting started

```sh
go get
go build
cp -v terraform-provider-authfed ~/.terraform.d/plugins/terraform-provider-authfed_v0.1.0
```

## Usage

```tf
provider "authfed" {
  cert_file = "public.pem"
  key_file = "private.pem"
}

resource "authfed_http_object" "hello_world" {
  url = "https://helloworld.net/blob/hello_world.json"
  content = <<EOF
{
  "hello": "world"
}
EOF
}
```

## License

Eclipse Public License v1.0, see LICENSE file.
