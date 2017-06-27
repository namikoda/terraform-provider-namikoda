#!/bin/bash

APIKEY="aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"

for executable in 'jq' 'cat'; do
  if ! type -P "$executable" >/dev/null; then
    echo "This script requires the '$executable' executable on the PATH"
    exit 1
  fi
done

cat <<EOF
resource "aws_waf_ipset" "newrelic-synthetics-ipset" {
  name = "newrelic-synthetics"
EOF

declare -a ipv4s
readarray -t ipv4s <<< "$(curl -sg -H "X-Namikoda-Key: $APIKEY" 'https://api.namikoda.com/v1/public/ipsfor/newrelic-synthetics' | jq -r '.ipv4s[]')"

for ipv4 in "${ipv4s[@]}"; do
  if [ -n "$ipv4" ]; then 
    cat <<EOF
  ip_set_descriptors {
    type  = "IPV4"
    value = "${ipv4}"
  }
EOF
  fi
done

declare -a ipv6s
readarray -t ipv6s <<< "$(curl -sg -H "X-Namikoda-Key: $APIKEY" 'https://api.namikoda.com/v1/public/ipsfor/newrelic-synthetics' | jq -r '.ipv6s[]')"

echo "ipv6s has ${#ipv6s[@]} elements"
for ipv6 in "${ipv6s[@]}"; do
  if [ -n "$ipv6" ]; then 
    cat <<EOF
  ip_set_descriptors {
    type  = "IPV6"
    value = "${ipv6}"
  }
EOF
  fi
done

cat <<EOF
resource "aws_waf_byte_match_set" "newrelic-synthetics-headermatch" {
  name = "newrelic-headermatch"

  byte_match_tuples {
    text_transformation   = "NONE"
    target_string         = "New Relic"
    positional_constraint = "CONTAINS"

    field_to_match {
      type = "HEADER"
      data = "X-Abuse-Info"
    }
  }
}

resource "aws_waf_rule" "newrelic-synthetics" {
  depends_on  = ["aws_waf_ipset.newrelic-synthetics-ipset"]
  name        = "newrelicSynthetics"
  metric_name = "newrelicSynthetics"

  predicates {
    data_id = "\${aws_waf_ipset.newrelic-synthetics-ipset.id}"
    negated = false
    type    = "IPMatch"
  }
  predicates {
    data_id = "\${aws_waf_byte_match_set.newrelic-synthetics-headermatch.id}"
    negated = false
    type    = "ByteMatch"
  }
}

resource "aws_waf_web_acl" "waf_acl" {
  name        = "tfWebACL"
  metric_name = "tfWebACL"

  default_action {
    type = "BLOCK"
  }

  rules {
    action {
      type = "ALLOW"
    }

    priority = 1
    rule_id  = "\${aws_waf_rule.newrelic-synthetics.id}"
  }
}
EOF

