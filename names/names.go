// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package names

import (
	"regexp"
	"strings"

	re2 "github.com/dlclark/regexp2" // for negative lookahead support
	"github.com/iancoleman/strcase"

	"github.com/aws-controllers-k8s/pkg/strutil"
)

var (
	nonAlphaNumRegexp *regexp.Regexp = regexp.MustCompile("[^a-zA-Z0-9]+")
)

type initialismTranslator struct {
	// CamelCased initialism, e.g. Tls
	camel string
	// Uppercase representation of the initialism
	upper string
	// Lowercase representation of the initialism
	lower string
	// Regular expression matching the initialism within a subject string.
	// Usually nil, unless the camel-cased initialism is a series of characters
	// that is commonly confused with a longer form of the initialism (e.g. for
	// "Id", we don't want to match "Identifier")
	re *re2.Regexp
}

var (
	// NOTE(jaypipes): these are ordered. Some things need to be processed
	// before others. For example, we need to process "Dbi" before "Db"
	initialisms = []initialismTranslator{
		// Special... even though IDS is a valid initialism, in AWS APIs, the
		// camel-cased "Ids" refers to a set of Identifiers, so the correct
		// uppercase representation is "IDs"
		{"Ids", "IDs", "ids", re2.MustCompile("(?![U|u])Ids", re2.None)},
		// Need to prevent "Identifier" from becoming "IDentifier", and "Idle"
		// from becoming "IDle" and "IdempotencyToken" from becoming
		// "IDempotencyToken"
		{"Id", "ID", "id", re2.MustCompile("Id(?!entifier|le|entity|empotency)", re2.None)},
		// Need to prevent "DbInstance" from becoming "dbinstance" when lower
		// prefix-converted (should be dbInstance). Amazingly, even within just
		// the RDS API, there are fields named "DbiResourceId",
		// "DBInstanceIdentifier" and "DbInstanceIdentifier" (note the
		// capitalization differences). This transformer handles this
		// problematic scenario and matches only the "Dbi" case-sensitive
		// expression and converts it to "DBI" or "dbi" depending on whether
		// the initialism appears at the start of the name
		{"Dbi", "DBI", "dbi", re2.MustCompile("Dbi", re2.None)},
		{"Db", "DB", "db", re2.MustCompile("Db(?!i)", re2.None)},
		{"Db", "DB", "db", re2.MustCompile("DB", re2.None)},
		// Prevent "CACertificateIdentifier" from becoming
		// "cACertificateIdentifier when lower prefix-converted (should be
		// "caCertificateIdentifier")
		{"CACert", "CACert", "caCert", re2.MustCompile("CACert", re2.None)},
		// Prevent "MD5OfBody" from becoming "MD5OfBody" when lower
		// prefix-converted (should be "md5OfBody")
		{"MD5Of", "MD5Of", "md5Of", re2.MustCompile("M[dD]5Of", re2.None)},
		// Prevent IPC from becoming IPc (ECS Task definition field)
		{"Ipc", "IPC", "ipc", re2.MustCompile("Ipc", re2.None)},
		// Prevent IPAddress from becoming iPAddress
		{"IPAddress", "IPAddress", "ip_address", nil},
		// Prevent IPv4 from becoming iPv4
		{"IPv4", "IPv4", "ipv4", re2.MustCompile("I[Pp]v4", re2.None)},
		{"IPv6", "IPv6", "ipv6", re2.MustCompile("I[Pp]v6", re2.None)},
		// Prevent "MultipartUpload" from becoming "MultIPartUpload"
		// and "IPAM" from becoming "IPam"
		{"Ip", "IP", "ip", re2.MustCompile("Ip(?!art|am)", re2.None)},
		{"IPSet", "IPSet", "ip_set", nil},
		// Model fields containing AMI will always capitalize the 'A' hence we don't
		// have to look for words starting with a lowercase 'A'
		{"Amis", "AMIs", "amis", re2.MustCompile("Amis", re2.None)},
		{"Ami", "AMI", "ami", re2.MustCompile("Ami", re2.None)},
		// Easy find-and-replacements...
		{"Acl", "ACL", "acl", nil},
		{"Acm", "ACM", "acm", nil},
		{"AIML", "AIML", "aiml", nil},
		{"Acp", "ACP", "acp", nil},
		{"Api", "API", "api", nil},
		{"Arn", "ARN", "arn", nil},
		{"Asn", "ASN", "asn", nil},
		// eventbridge has a NetworkConfiguration.awsvpcConfiguration field for
		// configuration of ECS tasks in "awsvpc" mode. aws-sdk-go transforms
		// this to AwsvpcConfiguration in order to export the field name in
		// Golang.
		// (See https://github.com/aws/aws-sdk-go/blob/5707eba1610d563b9c563dbc862587649bcb9811/service/eventbridge/api.go#L13088)
		// We need to prevent AwsvpcConfiguration from becoming
		// AWSvpcConfiguration
		{"Awsvpc", "AWSVPC", "awsVPC", nil},
		{"Aws", "AWS", "aws", nil},
		{"Az", "AZ", "az", nil},
		{"Bgp", "BGP", "bgp", nil},
		{"Cors", "CORS", "cors", nil},
		{"Cidr", "CIDR", "cidr", nil},
		{"Cname", "CNAME", "cname", nil},
		{"Cpu", "CPU", "cpu", nil},
		{"Crl", "CRL", "crl", nil},
		{"Cps", "CPS", "cps", nil},
		{"Csr", "CSR", "csr", nil},
		{"Dhcp", "DHCP", "dhcp", nil},
		{"Dns", "DNS", "dns", nil},
		{"Dpd", "DPD", "dpd", nil},
		{"Ebs", "EBS", "ebs", nil},
		{"Ec2", "EC2", "ec2", nil},
		// Prevent "Secret" from becoming "s_ecr_et"
		// Prevent "Decrease" from becoming "d_ecr_ease"
		{"Ecr", "ECR", "ecr", re2.MustCompile("(?!S|s|D|d)[Ee]cr(?!et|ease)", re2.None)},
		{"Ecs", "ECS", "ecs", nil},
		// Prevent "Edit" from becoming "EDIt"
		{"Edi", "EDI", "edi", re2.MustCompile("Edi(?!t)", re2.None)},
		{"Efs", "EFS", "efs", nil},
		{"Eks", "EKS", "eks", nil},
		// Prevent "Enable" and "Enabling" from becoming "ENAble"
		{"Ena", "ENA", "ena", re2.MustCompile("Ena(?!bl)", re2.None)},
		{"Ecmp", "ECMP", "ecmp", nil},
		{"Fifo", "FIFO", "fifo", nil},
		{"Fpga", "FPGA", "fpga", nil},
		{"Gid", "GID", "gid", nil},
		{"Gpu", "GPU", "gpu", nil},
		{"Grpc", "GRPC", "grpc", nil},
		{"Html", "HTML", "html", nil},
		// Prevent HTTPSPort from becoming httpSPort
		{"Http", "HTTP", "http", re2.MustCompile("(HTTP((?!S[A-Z]))|Http(?!s))", re2.None)},
		{"Https", "HTTPS", "https", nil},
		{"Iam", "IAM", "iam", nil},
		{"Icmp", "ICMP", "icmp", nil},
		// Prevent "IOPS" from becoming "IOps"
		{"Io", "IO", "io", re2.MustCompile("Io(?!ps)", re2.None)},
		{"Iops", "IOPS", "iops", nil},
		{"Ipam", "IPAM", "ipam", nil},
		{"Ja3", "JA3", "ja3", nil},
		{"Json", "JSON", "json", nil},
		{"Jwt", "JWT", "jwt", nil},
		{"Kms", "KMS", "kms", nil},
		{"Ldap", "LDAP", "ldap", nil},
		{"Mfa", "MFA", "mfa", nil},
		{"Mibps", "MiBps", "miBps", re2.MustCompile("Mibps", re2.None)},
		// Prevent "Native" from becoming "NATive"
		{"Nat", "NAT", "nat", re2.MustCompile("Nat(?!i)", re2.None)},
		// Prevent Oid from becoming oID and OIDC from becoming OIDc
		{"Oid", "OID", "oid", re2.MustCompile("Oid(?!c)", re2.None)},
		{"OID", "OID", "oid", re2.MustCompile("OID(?!C)", re2.None)},
		{"Oidc", "OIDC", "oidc", nil},
		{"Ocsp", "OCSP", "ocsp", nil},
		{"Pca", "PCA", "pca", nil},
		{"Pid", "PID", "pid", nil},
		// Capitalize the 'd' following RAM in certain cases
		{"Ramdisk", "RAMDisk", "ramDisk", re2.MustCompile("Ramdisk", re2.None)},
		// Model fields starting with 'Ram' refer to RAM
		{"Ram", "RAM", "ram", re2.MustCompile("Ram", re2.None)},
		{"Rfc", "RFC", "rfc", nil},
		{"Sasl", "SASL", "sasl", nil},
		{"Scram", "SCRAM", "scram", nil},
		{"Sdk", "SDK", "sdk", nil},
		{"Sha256", "SHA256", "sha256", nil},
		{"Sns", "SNS", "sns", nil},
		{"Sqli", "SQLI", "sqli", nil},
		{"Sql", "SQL", "sql", nil},
		{"Sqs", "SQS", "sqs", nil},
		{"Sriov", "SRIOV", "sriov", nil},
		{"Sse", "SSE", "sse", nil},
		{"Ssl", "SSL", "ssl", nil},
		{"Tcp", "TCP", "tcp", nil},
		{"Tde", "TDE", "tde", nil},
		{"Tpm", "TPM", "tpm", nil},
		{"Tls", "TLS", "tls", nil},
		{"Ttl", "TTL", "ttl", re2.MustCompile("(?!Thro)((?i)ttl)(?!ing|e)", re2.None)},
		{"Udp", "UDP", "udp", nil},
		// Need to prevent "security" from becoming "SecURIty"
		{"Uri", "URI", "uri", re2.MustCompile("(?!sec)uri(?!ty)|(Uri)", re2.None)},
		{"Url", "URL", "url", nil},
		{"Uuid", "UUID", "uuid", nil},
		{"Uids", "UIDs", "uids", re2.MustCompile("Uids", re2.None)},
		{"Uid", "UID", "uid", re2.MustCompile("Uid", re2.None)},
		// Need to prevent "Uid" or "Uuid" from becoming "UId" or "UUId"
		{"Ui", "UI", "ui", re2.MustCompile("U(I|i)(?!D|d)", re2.None)},
		{"Vlan", "VLAN", "vlan", nil},
		{"Vpc", "VPC", "vpc", nil},
		{"Vpn", "VPN", "vpn", nil},
		{"Vgw", "VGW", "vgw", nil},
		{"Waf", "WAF", "waf", nil},
		{"Xml", "XML", "xml", nil},
		{"Xss", "XSS", "xss", nil},
		{"Yaml", "YAML", "yaml", nil},
	}
)

var goKeywords = []string{
	"break",
	"case",
	"chan",
	"const",
	"continue",
	"default",
	"defer",
	"else",
	"fallthrough",
	"for",
	"func",
	"go",
	"goto",
	"if",
	"import",
	"interface",
	"map",
	"package",
	"range",
	"return",
	"select",
	"struct",
	"switch",
	"type",
	"var",
}

// Names contains variations of a name
type Names struct {
	Original      string
	Camel         string
	CamelLower    string
	Lower         string
	Snake         string
	SnakeStripped string
}

// New returns a Names containing variations of a supplied name
func New(original string) Names {
	return Names{
		Original:   original,
		Camel:      goName(original, false, false),
		CamelLower: goName(original, true, false),
		Lower:      strings.ToLower(original),
		Snake:      goName(original, false, true),
		SnakeStripped: nonAlphaNumRegexp.ReplaceAllString(
			goName(original, false, true), "",
		),
	}
}

func goName(original string, lowerFirst bool, snake bool) (result string) {
	result = original
	if !lowerFirst {
		result = strcase.ToCamel(result)
	}
	result, err := normalizeInitialisms(result, lowerFirst, snake)
	if err != nil {
		panic(err)
	}
	if lowerFirst {
		result, err = normalizeInitialisms(strcase.ToLowerCamel(result), lowerFirst, snake)
		if err != nil {
			panic(err)
		}
	}
	if snake {
		result = strcase.ToSnake(result)
	}
	if strutil.InStrings(result, goKeywords) {
		result = result + "_"
	}
	return
}

// normalizeInitialisms takes a subject string and adapts the string according
// to the Go best practice naming convention for initialisms.
//
// Examples:
//
//	original   | lowerFirst | output
//
// ------------+ ---------- + -------------------------
// Identifier  | true       | Identifier
// Identifier  | false      | Identifier
// Id          | true       | id
// Id          | false      | ID
// SSEKMSKeyId | true       | sseKMSKeyID
// SSEKMSKeyId | false      | SSEKMSKeyID
// RoleArn     | true       | roleARN
// RoleArn     | false      | RoleARN
//
// See: https://github.com/golang/go/wiki/CodeReviewComments#initialisms
func normalizeInitialisms(original string, lowerFirst bool, snake bool) (result string, err error) {
	result = original
	for _, initTrx := range initialisms {
		if initTrx.re == nil {
			if snake {
				// If we need to snakecase, we need to look for the uppercase
				// or lowercase initialism and replace with the lowercase
				// initialism plus an underscore. For example, if original ==
				// SSEKMSId and we pass snake == true, we want to return
				// sse_kms_key_id
				toReplace := "_" + initTrx.lower + "_"
				result = strings.Replace(result, initTrx.lower, toReplace, -1)
				result = strings.Replace(result, initTrx.upper, toReplace, -1)
				continue
			}
			if lowerFirst && strings.Index(result, initTrx.upper) == 0 {
				// if we need to lowercase initialisms, check to see if the
				// initialism's capitalized form starts the string, and if so,
				// lowercase it. For example, if we get original == SSEKMSKeyId
				// and we pass lower == true, we want to return sseKMSKeyID
				result = strings.Replace(result, initTrx.upper, initTrx.lower, 1)
			}
			// Replace CamelCased initialisms with the uppercase representation
			// of the initialism EXCEPT when the CamelCased initialism appears
			// at the start of the original string and we've passed a true
			// lower parameter, in which case we lowercase just the first
			// occurrence of the CamelCased initialism
			pos := strings.Index(result, initTrx.camel)
			switch pos {
			case -1:
				continue
			case 0:
				if lowerFirst {
					toReplace := initTrx.lower
					result = strings.Replace(result, initTrx.camel, toReplace, 1)
				}
				toReplace := initTrx.upper
				if snake {
					toReplace = "_" + toReplace + "_"
				}
				result = strings.Replace(result, initTrx.camel, toReplace, -1)
			default:
				toReplace := initTrx.upper
				if snake {
					toReplace = "_" + toReplace + "_"
				}
				result = strings.Replace(result, initTrx.camel, toReplace, -1)
			}
		} else {
			match, err := initTrx.re.FindStringMatch(result)
			if err != nil {
				return "", err
			}
			if match == nil {
				continue
			}
			startFrom := match.Group.Capture.Index
			if lowerFirst {
				if startFrom == 0 {
					// The matched string appears at the start of the string --
					// e.g. IdFirstElementId. In this case, if we've asked to lower
					// the output, we need to lower only the first occurrence of
					// the matched expression, not all of it -- e.g.
					// idFirstElementID
					toReplace := initTrx.lower
					result, err = initTrx.re.Replace(result, toReplace, 0, 1)
					if err != nil {
						return "", err
					}
					match, err = initTrx.re.FindNextMatch(match)
					if err != nil {
						return "", nil
					}
					if match == nil {
						continue
					}
					startFrom = match.Group.Capture.Index
				}
			}
			toReplace := initTrx.upper
			if snake {
				toReplace = "_" + initTrx.lower + "_"
			}
			result, err = initTrx.re.Replace(result, toReplace, startFrom, -1)
			if err != nil {
				return "", err
			}
		}
	}
	if snake {
		result = strings.Replace(result, "__", "_", -1)
		result = strings.Trim(result, "_")
	}
	return result, nil
}
