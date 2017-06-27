
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


