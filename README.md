# Trem Ligeiro Login

This applications is responsible for login of the restaurant "Trem Ligeiro" from Company "Trem Bão Tecnologia".

## How to run

### Requirements

-   [Go 1.22](https://golang.org/doc/install)
-   [Make](https://www.gnu.org/software/make/)
-   [AWS CLI](https://aws.amazon.com/pt/cli/)
-   [AWS SAM](https://aws.amazon.com/pt/serverless/sam/)
-   [Terraform 1.11.2](https://developer.hashicorp.com/terraform/install?product_intent=terraform)

### Building the Application

Use this commands to build the application.

```bash
make pre-build
make build
```

### Run locally

Starts applications using:

```bash
make run
```

### Run locally with SAM

Use this commands to build the application with SAM.

```bash
make sam-build
```

Starts applications using:
```
make sam-run
```

### Cloud
Configure GitHub Actions enviroment variables, then, starts DeployLambda workflow.