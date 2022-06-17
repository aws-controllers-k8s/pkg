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
		{"Ami", "AMI", "ami", "ami", "ami"},
		{"AmiLaunchIndex", "AMILaunchIndex", "amiLaunchIndex", "ami_launch_index", "amilaunchindex"},
		{"Amis", "AMIs", "amis", "amis", "amis"},
		{"AmiType", "AMIType", "amiType", "ami_type", "amitype"},
		{"CacheSecurityGroup", "CacheSecurityGroup", "cacheSecurityGroup", "cache_security_group", "cachesecuritygroup"},
		{"Camila", "Camila", "camila", "camila", "camila"},
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
		{"Examine", "Examine", "examine", "examine", "examine"},
		{"Family", "Family", "family", "family", "family"},
		{"Id", "ID", "id", "id", "id"},
		{"ID", "ID", "id", "id", "id"},
		{"Identifier", "Identifier", "identifier", "identifier", "identifier"},
		{"IoPerformance", "IOPerformance", "ioPerformance", "io_performance", "ioperformance"},
		{"Iops", "IOPS", "iops", "iops", "iops"},
		{"Ip", "IP", "ip", "ip", "ip"},
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
		{"Package", "Package", "package_", "package_", "package"},
		{"Param", "Param", "param", "param", "param"},
		{"Ram", "RAM", "ram", "ram", "ram"},
		{"RamdiskId", "RAMDiskID", "ramDiskID", "ram_disk_id", "ramdiskid"},
		{"RamDiskId", "RAMDiskID", "ramDiskID", "ram_disk_id", "ramdiskid"},
		{"RepositoryUriTest", "RepositoryURITest", "repositoryURITest", "repository_uri_test", "repositoryuritest"},
		{"RequestedAmiVersion", "RequestedAMIVersion", "requestedAMIVersion", "requested_ami_version", "requestedamiversion"},
		{"Sns", "SNS", "sns", "sns", "sns"},
		{"Sqs", "SQS", "sqs", "sqs", "sqs"},
		{"SriovNetSupport", "SRIOVNetSupport", "sriovNetSupport", "sriov_net_support", "sriovnetsupport"},
		{"SSEKMSKeyID", "SSEKMSKeyID", "sseKMSKeyID", "sse_kms_key_id", "ssekmskeyid"},
		{"UUID", "UUID", "uuid", "uuid", "uuid"},
		{"Vlan", "VLAN", "vlan", "vlan", "vlan"},
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
