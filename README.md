# Namikoda Terraform provider

#### Table of Contents
1. [Description](#description)
1. [Installing](#installing)
1. [Using](#using)
    * [Example: Allowing Github web hooks through an AWS security group](#example-allowing-github-web-hooks-through-an-aws-security-group)
1. [Custom IP ranges](#custom-ip-ranges)
    * [Example: Referencing a custom IP range](#example-referencing-a-custom-ip-range)
1. [Developing](#developing)
1. [Troubleshooting](#troubleshooting)
    * [401 Unauthorized](#http-request-error-response-code-401)
    * [404 Not Found](#http-request-error-response-code-404)
1. [Object Model](#object-model)
    * [Fields](#fields)
    * [Parameters](#parameters)

## Description

### IP range management in the cloud, made easy

Namikoda provides an easy interface for getting configuration values consistent across your physical and cloud infrastructure. We focus on providing IP range data and the tooling to act on it. Confidently use human-readable tags instead of lists of IPs in your configuration.

This module is an interface between Hashicorp's Terraform and the Namikoda API.  It is based on the excellent work of the [build in http provider](https://github.com/terraform-providers/terraform-provider-http).

## Installing
1. If you haven't yet, generate an API key at https://manage.namikoda.com.  See the [registration documentation](https://docs.namikoda.com/registration/index.html) for the step-by-step process.
1. Download or build the `terraform-provider-namikoda` binary and put it somewhere
1. edit `~/.terraformrc` and have it include the following

```
providers {
  ipsfor = "/path/to/terraform-provider-namikoda"
}
```

## Using
  The simplest use case gets IP address ranges from the `ipsfor` data component.
###  Example: Allowing Github web hooks through an AWS security group
```
data "ipsfor" "github" {
  apikey = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
  id = "github"
}

resource "aws_security_group" "github-ingress" {
  name        = "github-ingress"
  description = "Allow all inbound traffic from Github webhooks"

  ingress {
    from_port        = 0
    to_port          = 0
    protocol         = "-1"
    cidr_blocks      = [ "${data.ipsfor.github.ipv4s}" ]
    ipv6_cidr_blocks = [ "${data.ipsfor.github.ipv6s}" ]
  }
}
```

In this example, you create a data object using the `ipsfor` provider called `github`.  You'd put in your Namikoda API key from [the management portal](https://manage.namikoda.com), and tell the provider to get ranges with the Namikoda ID of `github`.  Then, craft an AWS security group called `github-ingress` with the results.

Keep in mind that there may be external limitations at play.  For instance, AWS [strongly limits the number of rules that a security group can have](http://docs.aws.amazon.com/AmazonVPC/latest/UserGuide/VPC_Appendix_Limits.html#vpc-limits-security-groups).


##  Custom IP ranges

If you have custom ranges managed in Namikoda, you can access them directly through the data provider.  

###  Example: Referencing a custom IP range
```
data "ipsfor" "office-ips" {
  owner = "my-namikoda-owner-name"
  apikey = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
  id = "office-ips"
}
```

Once you have the data in terraform, you can use it just like the above.


# Developing

Developing against this module works basically like any golang project.  Fork the repo, run `go get`, and start editing.  We do love PRs. :-)

The module is based on the excellent work of the [built in http provider](https://github.com/terraform-providers/terraform-provider-http) from the upstream terraform repo.  The README over there has some additional suggestions and requirements you should take a look at.

# Troubleshooting

This module passes through HTTP error codes to the user.  Here are some frequent examples

##  `HTTP request error. Response code: 401`
"401 Unauthorized" comes back when Namikoda can't verify that the string in the `apikey` field is a valid and current Namikoda API key.  Log in to https://manage.namikoda.com and make sure you've got the right apikey.  It will be in the format "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee".

## `HTTP request error. Response code: 404`
"404 Not Found" comes back when Namikoda does not know about an id for a given request.  If you're using public IP ranges, make sure that the ID is spelled correctly and matches an ID on the https://namikoda.com/ips.html page.  If you're using custom IP ranges, make sure the ID matches one you've put into the system at https://manage.namikoda.com and also that the `owner` matches your particular owner ID.

# Object Model

## Fields

The `ipsfor` object has the following fields:

### ipv4s
An array of strings in CIDR format that represent the IPv4 addresses under the particular id.

### ipv6s
An array of strings in CIDR format that represent the IPv6 addresses under the particular id.

### lastUpdate
A string in [ISO 8901 date and time](https://en.wikipedia.org/wiki/ISO_8601#Combined_date_and_time_representations) format that represents the last time the IP ranges were updated.

### name
A human-readable string describing who the id refers to.

### id
A short, unique string identifying the source of the IP information

### value
An array of strings in CIDR format that concatenate both the IPv4 and IPv6 arrays.

## Parameters

The `ipsfor` object takes the following parameters

### apikey
Required.  A string containing the Namikoda API key you can get from https://manage.namikoda.com

### id
Required.  A string containing the ID of the entity you're requesting IP information about.  Lists are available at https://manage.namikoda.com as well as https://namikoda.com/ips.html

### owner
Optional.  A string containing the unique identifier of your organization, used for pulling custom IP values.
