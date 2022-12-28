# OpenFaaS Templates

<!-- toc -->

* [Available Templates](#available-templates)
* [Creating a function](#creating-a-function)

<!-- Regenerate with "pre-commit run -a markdown-toc" -->

<!-- tocstop -->

Templates for [OpenFaaS](https://www.openfaas.com)

## Available Templates

 - [mongoose-crud](./template/mongoose-crud)
 - [golang-crud](./template/golang-crud)

## Creating a function

> This requires the [OpenFaaS CLI](https://github.com/openfaas/faas-cli)

```shell
faas-cli template pull https://github.com/mrsimonemms/openfaas-templates

faas-cli new <function-name> --lang <template-name>
```

For more detailed information, please see the [OpenFaaS documentation](https://docs.openfaas.com/cli/templates).
