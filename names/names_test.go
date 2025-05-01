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

package names_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/aws-controllers-k8s/pkg/names"
)

func TestNames(t *testing.T) {
	assert := assert.New(t)

	testCases := []struct {
		original            string
		expectCamel         string
		expectCamelLower    string
		expectSnake         string
		expectSnakeStripped string
	}{
		{"AIML", "AIML", "aiml", "aiml", "aiml"},
		{"Ami", "AMI", "ami", "ami", "ami"},
		{"Acm", "ACM", "acm", "acm", "acm"},
		{"AmiLaunchIndex", "AMILaunchIndex", "amiLaunchIndex", "ami_launch_index", "amilaunchindex"},
		{"Amis", "AMIs", "amis", "amis", "amis"},
		{"AmiType", "AMIType", "amiType", "ami_type", "amitype"},
		{"AwsVpcConfiguration", "AWSVPCConfiguration", "awsVPCConfiguration", "aws_vpc_configuration", "awsvpcconfiguration"},
		{"AWSVpcConfiguration", "AWSVPCConfiguration", "awsVPCConfiguration", "aws_vpc_configuration", "awsvpcconfiguration"},
		{"AWSVPCConfiguration", "AWSVPCConfiguration", "awsVPCConfiguration", "aws_vpc_configuration", "awsvpcconfiguration"},
		// eventbridge has a NetworkConfiguration.AwsvpcConfiguration field for
		// configuration of ECS tasks in "awsvpc" mode
		{"AwsvpcConfiguration", "AWSVPCConfiguration", "awsVPCConfiguration", "aws_vpc_configuration", "awsvpcconfiguration"},
		{"CacheSecurityGroup", "CacheSecurityGroup", "cacheSecurityGroup", "cache_security_group", "cachesecuritygroup"},
		{"Camila", "Camila", "camila", "camila", "camila"},
		{"AuthorizerResultTtlInSeconds", "AuthorizerResultTTLInSeconds", "authorizerResultTTLInSeconds", "authorizer_result_ttl_in_seconds", "authorizerresultttlinseconds"},
		{"DbInstanceId", "DBInstanceID", "dbInstanceID", "db_instance_id", "dbinstanceid"},
		{"DBInstanceId", "DBInstanceID", "dbInstanceID", "db_instance_id", "dbinstanceid"},
		{"DBInstanceID", "DBInstanceID", "dbInstanceID", "db_instance_id", "dbinstanceid"},
		{"DBInstanceIdentifier", "DBInstanceIdentifier", "dbInstanceIdentifier", "db_instance_identifier", "dbinstanceidentifier"},
		{"DbiResourceId", "DBIResourceID", "dbiResourceID", "dbi_resource_id", "dbiresourceid"},
		{"DpdTimeoutAction", "DPDTimeoutAction", "dpdTimeoutAction", "dpd_timeout_action", "dpdtimeoutaction"},
		{"Dynamic", "Dynamic", "dynamic", "dynamic", "dynamic"},
		{"Ecmp", "ECMP", "ecmp", "ecmp", "ecmp"},
		{"EdiPartyName", "EDIPartyName", "ediPartyName", "edi_party_name", "edipartyname"},
		{"Editable", "Editable", "editable", "editable", "editable"},
		{"Ena", "ENA", "ena", "ena", "ena"},
		{"Ecr", "ECR", "ecr", "ecr", "ecr"},
		{"Ecs", "ECS", "ecs", "ecs", "ecs"},
		{"Examine", "Examine", "examine", "examine", "examine"},
		{"Family", "Family", "family", "family", "family"},
		{"FifoTopic", "FIFOTopic", "fifoTopic", "fifo_topic", "fifotopic"},
		{"HttpsPort", "HTTPSPort", "httpsPort", "https_port", "httpsport"},
		{"HTTPSPort", "HTTPSPort", "httpsPort", "https_port", "httpsport"},
		{"Id", "ID", "id", "id", "id"},
		{"ID", "ID", "id", "id", "id"},
		{"Idle", "Idle", "idle", "idle", "idle"},
		{"IdempotencyToken", "IdempotencyToken", "idempotencyToken", "idempotency_token", "idempotencytoken"},
		{"Identifier", "Identifier", "identifier", "identifier", "identifier"},
		{"IoPerformance", "IOPerformance", "ioPerformance", "io_performance", "ioperformance"},
		{"Iops", "IOPS", "iops", "iops", "iops"},
		{"IPAddressType", "IPAddressType", "ipAddressType", "ip_address_type", "ipaddresstype"},
		{"IPSetType", "IPSetType", "ipSetType", "ip_set_type", "ipsettype"},
		{"Ip", "IP", "ip", "ip", "ip"},
		// The ipv_4/ipv_6 is a special case mainly caused by github.com/iancoleman/strcase
		// which does not handle this case correctly. See https://github.com/iancoleman/strcase/issues/22
		// This is OK for now as snake case is only used for package names and not for field names.
		{"IPv4", "IPv4", "ipv4", "ipv_4", "ipv4"},
		{"Ipv4", "IPv4", "ipv4", "ipv_4", "ipv4"},
		{"IPv6", "IPv6", "ipv6", "ipv_6", "ipv6"},
		{"Ipv6", "IPv6", "ipv6", "ipv_6", "ipv6"},
		{"Ipam", "IPAM", "ipam", "ipam", "ipam"},
		{"Ipc", "IPC", "ipc", "ipc", "ipc"},
		{"Ja3", "JA3", "ja3", "ja_3", "ja3"},
		{"Frame", "Frame", "frame", "frame", "frame"},
		{"KeyId", "KeyID", "keyID", "key_id", "keyid"},
		{"KeyID", "KeyID", "keyID", "key_id", "keyid"},
		{"KeyIdentifier", "KeyIdentifier", "keyIdentifier", "key_identifier", "keyidentifier"},
		{"LdapServerMetadata", "LDAPServerMetadata", "ldapServerMetadata", "ldap_server_metadata", "ldapservermetadata"},
		{"MaxIdleConnectionsPercent", "MaxIdleConnectionsPercent", "maxIdleConnectionsPercent", "max_idle_connections_percent", "maxidleconnectionspercent"},
		{"MultipartUpload", "MultipartUpload", "multipartUpload", "multipart_upload", "multipartupload"},
		{"Nat", "NAT", "nat", "nat", "nat"},
		{"NatGateway", "NATGateway", "natGateway", "nat_gateway", "natgateway"},
		{"NativeAuditFieldsIncluded", "NativeAuditFieldsIncluded", "nativeAuditFieldsIncluded", "native_audit_fields_included", "nativeauditfieldsincluded"},
		{"NumberOfAmiToKeep", "NumberOfAMIToKeep", "numberOfAMIToKeep", "number_of_ami_to_keep", "numberofamitokeep"},
		{"OID", "OID", "oid", "oid", "oid"},
		{"Oid", "OID", "oid", "oid", "oid"},
		{"Oidc", "OIDC", "oidc", "oidc", "oidc"},
		{"OIDC", "OIDC", "oidc", "oidc", "oidc"},
		{"Package", "Package", "package_", "package_", "package"},
		{"Param", "Param", "param", "param", "param"},
		{"Pid", "PID", "pid", "pid", "pid"},
		{"Pca", "PCA", "pca", "pca", "pca"},
		{"Uid", "UID", "uid", "uid", "uid"},
		{"Uids", "UIDs", "uids", "uids", "uids"},
		{"Gids", "GIDs", "gids", "gids", "gids"},
		{"Grpc", "GRPC", "grpc", "grpc", "grpc"},
		{"Ram", "RAM", "ram", "ram", "ram"},
		{"RamdiskId", "RAMDiskID", "ramDiskID", "ram_disk_id", "ramdiskid"},
		{"RamDiskId", "RAMDiskID", "ramDiskID", "ram_disk_id", "ramdiskid"},
		{"RepositoryUriTest", "RepositoryURITest", "repositoryURITest", "repository_uri_test", "repositoryuritest"},
		{"RequestedAmiVersion", "RequestedAMIVersion", "requestedAMIVersion", "requested_ami_version", "requestedamiversion"},
		{"SaslScram512Auth", "SASLSCRAM512Auth", "saslSCRAM512Auth", "sasl_scram_512_auth", "saslscram512auth"},
		{"Secret", "Secret", "secret", "secret", "secret"},
		{"Secrets", "Secrets", "secrets", "secrets", "secrets"},
		{"Sns", "SNS", "sns", "sns", "sns"},
		{"Sqli", "SQLI", "sqli", "sqli", "sqli"},
		{"Sql", "SQL", "sql", "sql", "sql"},
		{"Sqs", "SQS", "sqs", "sqs", "sqs"},
		{"SriovNetSupport", "SRIOVNetSupport", "sriovNetSupport", "sriov_net_support", "sriovnetsupport"},
		{"SSEKMSKeyID", "SSEKMSKeyID", "sseKMSKeyID", "sse_kms_key_id", "ssekmskeyid"},
		{"Tpm", "TPM", "tpm", "tpm", "tpm"},
		{"TTL", "TTL", "ttl", "ttl", "ttl"},
		{"Throttle", "Throttle", "throttle", "throttle", "throttle"},
		{"Throttling", "Throttling", "throttling", "throttling", "throttling"},
		{"UUID", "UUID", "uuid", "uuid", "uuid"},
		{"Vlan", "VLAN", "vlan", "vlan", "vlan"},
		{"Xss", "XSS", "xss", "xss", "xss"},
		{"MiBps", "MiBps", "miBps", "mi_bps", "mibps"},
		{"LastDecreaseDateTime", "LastDecreaseDateTime", "lastDecreaseDateTime", "last_decrease_date_time", "lastdecreasedatetime"},
		{"NumberOfDecreasesToday", "NumberOfDecreasesToday", "numberOfDecreasesToday", "number_of_decreases_today", "numberofdecreasestoday"},
	}
	for _, tc := range testCases {
		n := names.New(tc.original)
		msg := fmt.Sprintf("for original %s expected camel name of %s but got %s", tc.original, tc.expectCamel, n.Camel)
		assert.Equal(tc.expectCamel, n.Camel, msg)
		msg = fmt.Sprintf("for original %s expected lowercase camel name of %s but got %s", tc.original, tc.expectCamelLower, n.CamelLower)
		assert.Equal(tc.expectCamelLower, n.CamelLower, msg)
		msg = fmt.Sprintf("for original %s expected snake name of %s but got %s", tc.original, tc.expectSnake, n.Snake)
		assert.Equal(tc.expectSnake, n.Snake, msg)
		msg = fmt.Sprintf("for original %s expected snake stripped name of %s but got %s", tc.original, tc.expectSnakeStripped, n.SnakeStripped)
		assert.Equal(tc.expectSnakeStripped, n.SnakeStripped, msg)
	}
}
