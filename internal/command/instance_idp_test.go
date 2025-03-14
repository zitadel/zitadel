package command

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	openid "github.com/zitadel/oidc/v3/pkg/oidc"
	"go.uber.org/mock/gomock"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/id"
	id_mock "github.com/zitadel/zitadel/internal/id/mock"
	"github.com/zitadel/zitadel/internal/repository/idp"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var (
	validSAMLMetadata = []byte(`<?xml version="1.0" encoding="UTF-8"?>
<EntityDescriptor xmlns="urn:oasis:names:tc:SAML:2.0:metadata" entityID="http://localhost:8080/saml/v2/metadata" ID="_8b02ecf6-aea4-4eda-96c6-190551f05b07">
  <Signature xmlns="http://www.w3.org/2000/09/xmldsig#">
    <SignedInfo xmlns="http://www.w3.org/2000/09/xmldsig#">
      <CanonicalizationMethod xmlns="http://www.w3.org/2000/09/xmldsig#" Algorithm="http://www.w3.org/2001/10/xml-exc-c14n#"></CanonicalizationMethod>
      <SignatureMethod xmlns="http://www.w3.org/2000/09/xmldsig#" Algorithm="http://www.w3.org/2001/04/xmldsig-more#rsa-sha256"></SignatureMethod>
      <Reference xmlns="http://www.w3.org/2000/09/xmldsig#" URI="#_8b02ecf6-aea4-4eda-96c6-190551f05b07">
        <Transforms xmlns="http://www.w3.org/2000/09/xmldsig#">
          <Transform xmlns="http://www.w3.org/2000/09/xmldsig#" Algorithm="http://www.w3.org/2000/09/xmldsig#enveloped-signature"></Transform>
          <Transform xmlns="http://www.w3.org/2000/09/xmldsig#" Algorithm="http://www.w3.org/2001/10/xml-exc-c14n#"></Transform>
        </Transforms>
        <DigestMethod xmlns="http://www.w3.org/2000/09/xmldsig#" Algorithm="http://www.w3.org/2001/04/xmlenc#sha256"></DigestMethod>
        <DigestValue xmlns="http://www.w3.org/2000/09/xmldsig#">Tyw4csdpNNq0E7wi5FXWdVNkdPNg+cM6kK21VB2+iF0=</DigestValue>
      </Reference>
    </SignedInfo>
    <SignatureValue xmlns="http://www.w3.org/2000/09/xmldsig#">hWQSYmnBJENy/okk2qRDuHaZiyqpDsdV6BF9/T/LNjUh/8z4dV2NEZvkNhFEyQ+bqdj+NmRWvKqpg1dtgNJxQc32+IsLQvXNYyhMCtyG570/jaTOtm8daV4NKJyTV7SdwM6yfXgubz5YCRTyV13W2gBIFYppIRImIv5NDcjz+lEmWhnrkw8G2wRSFUY7VvkDn9rgsTzw/Pnsw6hlzpjGDYPMPx3ux3kjFVevdhFGNo+VC7t9ozruuGyH3yue9Re6FZoqa4oyWaPSOwei0ZH6UNqkX93Eo5Y49QKwaO8Rm+kWsOhdTqebVmCc+SpWbbrZbQj4nSLgWGlvCkZSivmH7ezr4Ol1ZkRetQ92UQ7xJS7E0y6uXAGvdgpDnyqHCOFfhTS6yqltHtc3m7JZex327xkv6e69uAEOSiv++sifVUIE0h/5u3hZLvwmTPrkoRVY4wgZ4ieb86QPvhw4UPeYapOhCBk5RfjoEFIeYwPUw5rtOlpTyeBJiKMpH1+mDAoa+8HQytZoMrnnY1s612vINtY7jU5igMwIk6MitQpRGibnBVBHRc2A6aE+XS333ganFK9hX6TzNkpHUb66NINDZ8Rgb1thn3MABArGlomtM5/enrAixWExZp70TSElor7SBdBW57H7OZCYUCobZuPRDLsCO6LLKeVrbdygWeRqr/o=</SignatureValue>
    <KeyInfo xmlns="http://www.w3.org/2000/09/xmldsig#">
	  <X509Data xmlns="http://www.w3.org/2000/09/xmldsig#">
		<X509Certificate xmlns="http://www.w3.org/2000/09/xmldsig#">MIIFIjCCAwqgAwIBAgICA7YwDQYJKoZIhvcNAQELBQAwLDEQMA4GA1UEChMHWklUQURFTDEYMBYGA1UEAxMPWklUQURFTCBTQU1MIENBMB4XDTI0MTEyNzEwMjc0NFoXDTI1MTEyNzE2Mjc0NFowMjEQMA4GA1UEChMHWklUQURFTDEeMBwGA1UEAxMVWklUQURFTCBTQU1MIG1ldGFkYXRhMIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEApEpYT7EjbRBp0Hw7PGCiSgUoJtwd2nwZOhGy5WZVWvraAtHzW5ih2B6UwEShjwCmRJZeKYEN9JKJbpAy2EdL/l2rm/pArVNvSQu6sN4izz5p2rd9NfHAO3/EcvYdrelWLQj8WQx6LVM282Z4wbclp8Jz1y8Ow43352hGfFVc1x8gauoNl5MAy4kdbvs8UqihqcRmEyIOWl6UwTApb+XIRSRz0Yop99Fv9ALJwfUppsx+d4j9rlRDvrQJMJz7GC/19L9INTbY0HsVEiTltdAWHwREwrpwxNJQt42p3W/zpf1mjwXd3qNNDZAr1t2POPP4SXd598kabBZ3EMWGGxFw+NYYajyjG5EFOZw09FFJn2jIcovejvigfdqem5DGPECvHefqcqHkBPGukI3RaotXpAYyAGfnV7slVytSW484IX3KloAJLICbETbFGGsGQzIDw8rUqWyaOCOttw2fVNDyRFUMHrGe1PhJ9qA1If+KCWYD0iJqF03rIEhdrvNSdQNYkRa0DdtpacQLpzQtqsUioODqX0W3uzLceJEXLBbU0ZEk8mWZM/auwMo3ycPNXDVwrb6AkUKar+sqSumUuixw7da3KF1/mynh6M2Eo4NRB16oUiyN0EYrit/RRJjsTdH+71cj0V+8KqO88cBpmm+lO6x4RM5xpOf/EwwQHivxgRkCAwEAAaNIMEYwDgYDVR0PAQH/BAQDAgWgMBMGA1UdJQQMMAoGCCsGAQUFBwMCMB8GA1UdIwQYMBaAFIzl7uckcPWldirXeOFL3rH6K8FLMA0GCSqGSIb3DQEBCwUAA4ICAQBz+7R99uX1Us9T4BB2RK3RD9K8Q5foNmxJ8GbxpOQFL8IG1DE3FqBssciJkOsKY+1+Y6eow2TgmD9MxfCY444C8k8YDDjxIcs+4dEaWMUxA6NoEy378ciy0U1E6rpYLxWYTxXmsELyODWwTrRNIiWfbBD2m0w9HYbK6QvX6IYQqYoTOJJ3WJKsMCeQ8XhQsJYNINZEq8RsERY/aikOlTWN7ax4Mkr3bfnz1euXGClExCOM6ej4m2I33i4nyYBvvRkRRZRQCfkAQ+5WFVZoVXrQHNe/Oifit7tfLaDuybcjgkzzY3o0YbczzbdV69fVoj53VpR3QQOB+PCF/VJPUMtUFPEC05yH76g24KVBiM/Ws8GaERW1AxgupHSmvTY3GSiwDXQ2NzgDxUHfRHo8rxenJdEcPlGM0DstbUONDSFGLwvGDiidUVtqj1UB4yGL26bgtmwf61G4qsTn9PJMWdRmCeeOf7fmloRxTA0EEey3bulBBHim466tWHUhgOP+g1X0iE7CnwL8aJ//CCiQOAv1O6x5RLyxrmVTehPLr1T8qvnBmxpmuYU0kfbYpO3tMVe7VLabBx0cYh7izClZKHhgEj1w4aE9tIk7nqVAwvVocT3io8RrcKixlnBrFd7RYIuF3+RsYC/kYEgnZYKAig5u2TySgGmJ7nIS24FYW68WDg==</X509Certificate>
      </X509Data>
	</KeyInfo>
  </Signature>
  <IDPSSODescriptor xmlns="urn:oasis:names:tc:SAML:2.0:metadata" WantAuthnRequestsSigned="1" ID="_fd70402c-8a31-4a9a-a4a7-da526524c609" validUntil="2024-12-02T16:54:55.656Z" protocolSupportEnumeration="urn:oasis:names:tc:SAML:2.0:protocol">
	<SingleSignOnService xmlns="urn:oasis:names:tc:SAML:2.0:metadata" Binding="urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect" Location="http://localhost:8080/saml/v2/SSO"></SingleSignOnService>
	<SingleSignOnService xmlns="urn:oasis:names:tc:SAML:2.0:metadata" Binding="urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST" Location="http://localhost:8080/saml/v2/SSO"></SingleSignOnService>
	<AttributeProfile>urn:oasis:names:tc:SAML:2.0:profiles:attribute:basic</AttributeProfile>
	<Attribute xmlns="urn:oasis:names:tc:SAML:2.0:assertion" Name="Email" NameFormat="urn:oasis:names:tc:SAML:2.0:attrname-format:basic"><AttributeValue></AttributeValue></Attribute>
	<Attribute xmlns="urn:oasis:names:tc:SAML:2.0:assertion" Name="SurName" NameFormat="urn:oasis:names:tc:SAML:2.0:attrname-format:basic"><AttributeValue></AttributeValue></Attribute>
	<Attribute xmlns="urn:oasis:names:tc:SAML:2.0:assertion" Name="FirstName" NameFormat="urn:oasis:names:tc:SAML:2.0:attrname-format:basic"><AttributeValue></AttributeValue></Attribute>
	<Attribute xmlns="urn:oasis:names:tc:SAML:2.0:assertion" Name="FullName" NameFormat="urn:oasis:names:tc:SAML:2.0:attrname-format:basic"><AttributeValue></AttributeValue></Attribute>
	<Attribute xmlns="urn:oasis:names:tc:SAML:2.0:assertion" Name="UserName" NameFormat="urn:oasis:names:tc:SAML:2.0:attrname-format:basic"><AttributeValue></AttributeValue></Attribute>
	<Attribute xmlns="urn:oasis:names:tc:SAML:2.0:assertion" Name="UserID" NameFormat="urn:oasis:names:tc:SAML:2.0:attrname-format:basic"><AttributeValue></AttributeValue></Attribute>
	<SingleLogoutService xmlns="urn:oasis:names:tc:SAML:2.0:metadata" Binding="urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect" Location="http://localhost:8080/saml/v2/SLO"></SingleLogoutService>
	<SingleLogoutService xmlns="urn:oasis:names:tc:SAML:2.0:metadata" Binding="urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST" Location="http://localhost:8080/saml/v2/SLO"></SingleLogoutService>
	<NameIDFormat>urn:oasis:names:tc:SAML:2.0:nameid-format:persistent</NameIDFormat>
	<KeyDescriptor xmlns="urn:oasis:names:tc:SAML:2.0:metadata" use="signing">
      <KeyInfo xmlns="http://www.w3.org/2000/09/xmldsig#">
		<KeyName>http://localhost:8080/saml/v2/metadata IDP signing</KeyName>
		<X509Data xmlns="http://www.w3.org/2000/09/xmldsig#">
		  <X509Certificate xmlns="http://www.w3.org/2000/09/xmldsig#">MIIFIjCCAwqgAwIBAgICA7QwDQYJKoZIhvcNAQELBQAwLDEQMA4GA1UEChMHWklUQURFTDEYMBYGA1UEAxMPWklUQURFTCBTQU1MIENBMB4XDTI0MTEyNzEwMjUwMloXDTI1MTEyNzE2MjUwMlowMjEQMA4GA1UEChMHWklUQURFTDEeMBwGA1UEAxMVWklUQURFTCBTQU1MIHJlc3BvbnNlMIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEA2lUgaI6AS/9xvM9DNSWK6Ho64LpK8UIioM26QfvAfeQ/I2pgX6SwWxEbd7qv+PkJzaFTjrXSlwOmWsJYma+UsdyFClaGFRyCgY8SWxPceandC8a+hQIDS/irLd9XF33RWp0b/09HjQl+n0HZ4teUFDUd2U1mUf3XCpn0+Ho316bmi6xSW6zaMy5RsbUl01hgWj2fgapAsGAHSBphwCE3Dz/9I/UfHWQw1k2/UTgjc9uIujcza6WgOxfsKluXYIOxwNKTfmzzOJMUwXz6GRgB2jhQI29MuKOZOITA7pXq5kZKf0lSRU8zKFTMJaK4zAHQ6f877Drr8XdAHemuXGZ2JdH/Dbdwarzy3YBMCWsAYlpeEvaVAdiSpyR7fAZktNuHd39Zg00Vlj2wdc44Vk5yVssW7pv5qnVZ7JTrXX2uBYFecLAXmplQ2ph1VdSXZLEDGgjiNA2T/fBj7G4/VjsuCBZFm1I0KCJp3HWEJx5dwwhSVc5wOJEzl7fMuPYMKWH/RM6P/7LnO1ulpdmiKPa4gHzdg3hDZn42NKcVt3UYf0phtxpWMrZp/DUEeizhckrC4ed6cfGtS3CUtJEqoycrCROJ5Hy+ONHl5Aqxt+JoPU+t/XATuctfPxQVcDr0itHzo2cjh/AVTU+IC7C0oQHSS9CC8Fp58UqbtYwFtSAd7ecCAwEAAaNIMEYwDgYDVR0PAQH/BAQDAgWgMBMGA1UdJQQMMAoGCCsGAQUFBwMCMB8GA1UdIwQYMBaAFIzl7uckcPWldirXeOFL3rH6K8FLMA0GCSqGSIb3DQEBCwUAA4ICAQAp+IGZScVIbRdCq5HPjlYBPOY7UbL8ZXnlMW/HLELV9GndnULuFhnuQTIdA5dquCsk8RI1fKsScEV1rqWvHZeSo5nVbvUaPJctoD/4GACqE6F8axs1AgSOvpJMyuycjSzSh6gDM1z37Fdqc/2IRqgi7SKdDsfJpi8XW8LtErpp4kyE1rEXopsXG2fe1UH25bZpXraUqYvp61rwVUCazAtV/U7ARG5AnT0mPqzUriIPrfL+v/+2ntV/BSc8/uCqYnHbwpIwjPURCaxo1Pmm6EEkm+V/Ss4ieNwwkD2bLLLST1LoVMim7Ebfy53PEKpsznKsGlVSu0YYKUsStWQVpwhKQw0bQLCJHdpvZtZSDgS9RbSMZz+aY/fpoNx6wDvmMgtdrb3pVXZ8vPKdq9YDrGfFqP60QdZ3CuSHXCM/zX4742GgImJ4KYAcTuF1+BkGf5JLAJOUZBkfCQ/kBT5wr8+EotLxASOC6717whLBYMEG6N8osEk+LDqoJRTLqkzirJsyOHWChKK47yGkdS3HBIZfo91QrJwKpfATYziBjEnqipkTu+6jFylBIkxKTPye4b3vgcodZP8LSNVXAsMGTPNPJxzPWQ37ba4zMnYZ5iUerlaox/SNsn68DT6RajIb1A1JDq+HNFc3hQP2bzk2y5pCax8zo5swjdklnm4clfB2Lw==</X509Certificate>
		</X509Data>
      </KeyInfo>
	</KeyDescriptor>
  </IDPSSODescriptor>
  <AttributeAuthorityDescriptor xmlns="urn:oasis:names:tc:SAML:2.0:metadata" ID="_b3fed381-af56-4160-abf5-5ffd1e21cf61" validUntil="2024-12-02T16:54:55.656Z" protocolSupportEnumeration="urn:oasis:names:tc:SAML:2.0:protocol">
	<AttributeService xmlns="urn:oasis:names:tc:SAML:2.0:metadata" Binding="urn:oasis:names:tc:SAML:2.0:bindings:SOAP" Location="http://localhost:8080/saml/v2/attribute"></AttributeService>
	<NameIDFormat>urn:oasis:names:tc:SAML:2.0:nameid-format:persistent</NameIDFormat>
	<AttributeProfile>urn:oasis:names:tc:SAML:2.0:profiles:attribute:basic</AttributeProfile>
	<Attribute xmlns="urn:oasis:names:tc:SAML:2.0:assertion" Name="Email" NameFormat="urn:oasis:names:tc:SAML:2.0:attrname-format:basic"><AttributeValue></AttributeValue></Attribute>
	<Attribute xmlns="urn:oasis:names:tc:SAML:2.0:assertion" Name="SurName" NameFormat="urn:oasis:names:tc:SAML:2.0:attrname-format:basic"><AttributeValue></AttributeValue></Attribute>
	<Attribute xmlns="urn:oasis:names:tc:SAML:2.0:assertion" Name="FirstName" NameFormat="urn:oasis:names:tc:SAML:2.0:attrname-format:basic"><AttributeValue></AttributeValue></Attribute>
	<Attribute xmlns="urn:oasis:names:tc:SAML:2.0:assertion" Name="FullName" NameFormat="urn:oasis:names:tc:SAML:2.0:attrname-format:basic"><AttributeValue></AttributeValue></Attribute>
	<Attribute xmlns="urn:oasis:names:tc:SAML:2.0:assertion" Name="UserName" NameFormat="urn:oasis:names:tc:SAML:2.0:attrname-format:basic"><AttributeValue></AttributeValue></Attribute>
	<Attribute xmlns="urn:oasis:names:tc:SAML:2.0:assertion" Name="UserID" NameFormat="urn:oasis:names:tc:SAML:2.0:attrname-format:basic"><AttributeValue></AttributeValue></Attribute>
	<KeyDescriptor xmlns="urn:oasis:names:tc:SAML:2.0:metadata" use="signing">
	  <KeyInfo xmlns="http://www.w3.org/2000/09/xmldsig#">
		<KeyName>http://localhost:8080/saml/v2/metadata IDP signing</KeyName>
  		<X509Data xmlns="http://www.w3.org/2000/09/xmldsig#">
		  <X509Certificate xmlns="http://www.w3.org/2000/09/xmldsig#">MIIFIjCCAwqgAwIBAgICA7QwDQYJKoZIhvcNAQELBQAwLDEQMA4GA1UEChMHWklUQURFTDEYMBYGA1UEAxMPWklUQURFTCBTQU1MIENBMB4XDTI0MTEyNzEwMjUwMloXDTI1MTEyNzE2MjUwMlowMjEQMA4GA1UEChMHWklUQURFTDEeMBwGA1UEAxMVWklUQURFTCBTQU1MIHJlc3BvbnNlMIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEA2lUgaI6AS/9xvM9DNSWK6Ho64LpK8UIioM26QfvAfeQ/I2pgX6SwWxEbd7qv+PkJzaFTjrXSlwOmWsJYma+UsdyFClaGFRyCgY8SWxPceandC8a+hQIDS/irLd9XF33RWp0b/09HjQl+n0HZ4teUFDUd2U1mUf3XCpn0+Ho316bmi6xSW6zaMy5RsbUl01hgWj2fgapAsGAHSBphwCE3Dz/9I/UfHWQw1k2/UTgjc9uIujcza6WgOxfsKluXYIOxwNKTfmzzOJMUwXz6GRgB2jhQI29MuKOZOITA7pXq5kZKf0lSRU8zKFTMJaK4zAHQ6f877Drr8XdAHemuXGZ2JdH/Dbdwarzy3YBMCWsAYlpeEvaVAdiSpyR7fAZktNuHd39Zg00Vlj2wdc44Vk5yVssW7pv5qnVZ7JTrXX2uBYFecLAXmplQ2ph1VdSXZLEDGgjiNA2T/fBj7G4/VjsuCBZFm1I0KCJp3HWEJx5dwwhSVc5wOJEzl7fMuPYMKWH/RM6P/7LnO1ulpdmiKPa4gHzdg3hDZn42NKcVt3UYf0phtxpWMrZp/DUEeizhckrC4ed6cfGtS3CUtJEqoycrCROJ5Hy+ONHl5Aqxt+JoPU+t/XATuctfPxQVcDr0itHzo2cjh/AVTU+IC7C0oQHSS9CC8Fp58UqbtYwFtSAd7ecCAwEAAaNIMEYwDgYDVR0PAQH/BAQDAgWgMBMGA1UdJQQMMAoGCCsGAQUFBwMCMB8GA1UdIwQYMBaAFIzl7uckcPWldirXeOFL3rH6K8FLMA0GCSqGSIb3DQEBCwUAA4ICAQAp+IGZScVIbRdCq5HPjlYBPOY7UbL8ZXnlMW/HLELV9GndnULuFhnuQTIdA5dquCsk8RI1fKsScEV1rqWvHZeSo5nVbvUaPJctoD/4GACqE6F8axs1AgSOvpJMyuycjSzSh6gDM1z37Fdqc/2IRqgi7SKdDsfJpi8XW8LtErpp4kyE1rEXopsXG2fe1UH25bZpXraUqYvp61rwVUCazAtV/U7ARG5AnT0mPqzUriIPrfL+v/+2ntV/BSc8/uCqYnHbwpIwjPURCaxo1Pmm6EEkm+V/Ss4ieNwwkD2bLLLST1LoVMim7Ebfy53PEKpsznKsGlVSu0YYKUsStWQVpwhKQw0bQLCJHdpvZtZSDgS9RbSMZz+aY/fpoNx6wDvmMgtdrb3pVXZ8vPKdq9YDrGfFqP60QdZ3CuSHXCM/zX4742GgImJ4KYAcTuF1+BkGf5JLAJOUZBkfCQ/kBT5wr8+EotLxASOC6717whLBYMEG6N8osEk+LDqoJRTLqkzirJsyOHWChKK47yGkdS3HBIZfo91QrJwKpfATYziBjEnqipkTu+6jFylBIkxKTPye4b3vgcodZP8LSNVXAsMGTPNPJxzPWQ37ba4zMnYZ5iUerlaox/SNsn68DT6RajIb1A1JDq+HNFc3hQP2bzk2y5pCax8zo5swjdklnm4clfB2Lw==</X509Certificate>
		</X509Data>
	  </KeyInfo>
	</KeyDescriptor>
  </AttributeAuthorityDescriptor>
</EntityDescriptor>`)
	validLDAPRootCA = []byte(`-----BEGIN CERTIFICATE-----
MIIDITCCAgmgAwIBAgIUKjAUmxsHO44X+/TKBNciPgNl1GEwDQYJKoZIhvcNAQEL
BQAwIDEeMBwGA1UEAwwVbXlzZXJ2aWNlLmV4YW1wbGUuY29tMB4XDTI0MTIxOTEz
Mzc1MVoXDTI1MTIxOTEzMzc1MVowIDEeMBwGA1UEAwwVbXlzZXJ2aWNlLmV4YW1w
bGUuY29tMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA0QYuJsayILRI
hVT7G1DlitVSXnt1iw3gEXJZfe81Egz06fUbvXF6Yo1LJmwYpqe/rm+hf4FNUb8e
2O+LH2FieA9FkVe4P2gKOzw87A/KxvpV8stgNgl4LlqRCokbc1AzeE/NiLr5TcTD
RXm3DUcYxXxinprtDu2jftFysaOZmNAukvE/iL6qS3X6ggVEDDM7tY9n5FV2eJ4E
p0ImKfypi2aZYROxOK+v5x9ryFRMl4y07lMDvmtcV45uXYmfGNCgG9PNf91Kk/mh
JxEQbxycJwFoSi9XWljR8ahPdO11LXG7Dsj/RVbY8k2LdKNstl6Ae3aCpbe9u2Pj
vxYs1bVJuQIDAQABo1MwUTAdBgNVHQ4EFgQU+mRVN5HYJWgnpopReaLhf2cMcoYw
HwYDVR0jBBgwFoAU+mRVN5HYJWgnpopReaLhf2cMcoYwDwYDVR0TAQH/BAUwAwEB
/zANBgkqhkiG9w0BAQsFAAOCAQEABJpHVuc9tGhD04infRVlofvqXIUizTlOrjZX
vozW9pIhSWEHX8o+sJP8AMZLnrsdq+bm0HE0HvgYrw7Lb8pd4FpR46TkFHjeukoj
izqfgckjIBl2nwPGlynbKA0/U/rTCSxVt7XiAn+lgYUGIpOzNdk06/hRMitrMNB7
t2C97NseVC4b1ZgyFrozsefCfUmD8IJF0+XJ4Wzmsh0jRrI8koCtVmPYnKn6vw1b
cZprg/97CWHYrsavd406wOB60CMtYl83Q16ucOF1dretDFqJC5kY+aFLvuqfag2+
kIaoPV1MnGsxveQyyHdOsEatS5XOv/1OWcmnvePDPxcvb9jCcw==
-----END CERTIFICATE-----
`)
)

func TestCommandSide_AddInstanceGenericOAuthIDP(t *testing.T) {
	type fields struct {
		eventstore   func(*testing.T) *eventstore.Eventstore
		idGenerator  id.Generator
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx      context.Context
		provider GenericOAuthProvider
	}
	type res struct {
		id   string
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"invalid name",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				provider: GenericOAuthProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-D32ef", ""))
				},
			},
		},
		{
			"invalid clientID",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: GenericOAuthProvider{
					Name: "name",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-Dbgzf", ""))
				},
			},
		},
		{
			"invalid clientSecret",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: GenericOAuthProvider{
					Name:     "name",
					ClientID: "clientID",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-DF4ga", ""))
				},
			},
		},
		{
			"invalid auth endpoint",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: GenericOAuthProvider{
					Name:         "name",
					ClientID:     "clientID",
					ClientSecret: "clientSecret",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-B23bs", ""))
				},
			},
		},
		{
			"invalid token endpoint",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: GenericOAuthProvider{
					Name:                  "name",
					ClientID:              "clientID",
					ClientSecret:          "clientSecret",
					AuthorizationEndpoint: "auth",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-D2gj8", ""))
				},
			},
		},
		{
			"invalid user endpoint",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: GenericOAuthProvider{
					Name:                  "name",
					ClientID:              "clientID",
					ClientSecret:          "clientSecret",
					AuthorizationEndpoint: "auth",
					TokenEndpoint:         "token",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-Fb8jk", ""))
				},
			},
		},
		{
			"invalid id attribute",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: GenericOAuthProvider{
					Name:                  "name",
					ClientID:              "clientID",
					ClientSecret:          "clientSecret",
					AuthorizationEndpoint: "auth",
					TokenEndpoint:         "token",
					UserEndpoint:          "user",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-sdf3f", ""))
				},
			},
		},
		{
			name: "ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectPush(
						instance.NewOAuthIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
							"id1",
							"name",
							"clientID",
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("clientSecret"),
							},
							"auth",
							"token",
							"user",
							"idAttribute",
							nil,
							true,
							idp.Options{},
						),
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: GenericOAuthProvider{
					Name:                  "name",
					ClientID:              "clientID",
					ClientSecret:          "clientSecret",
					AuthorizationEndpoint: "auth",
					TokenEndpoint:         "token",
					UserEndpoint:          "user",
					IDAttribute:           "idAttribute",
					UsePKCE:               true,
				},
			},
			res: res{
				id:   "id1",
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
		{
			name: "ok all set",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectPush(
						instance.NewOAuthIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
							"id1",
							"name",
							"clientID",
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("clientSecret"),
							},
							"auth",
							"token",
							"user",
							"idAttribute",
							[]string{"user"},
							true,
							idp.Options{
								IsCreationAllowed: true,
								IsLinkingAllowed:  true,
								IsAutoCreation:    true,
								IsAutoUpdate:      true,
							},
						),
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: GenericOAuthProvider{
					Name:                  "name",
					ClientID:              "clientID",
					ClientSecret:          "clientSecret",
					AuthorizationEndpoint: "auth",
					TokenEndpoint:         "token",
					UserEndpoint:          "user",
					Scopes:                []string{"user"},
					IDAttribute:           "idAttribute",
					UsePKCE:               true,
					IDPOptions: idp.Options{
						IsCreationAllowed: true,
						IsLinkingAllowed:  true,
						IsAutoCreation:    true,
						IsAutoUpdate:      true,
					},
				},
			},
			res: res{
				id:   "id1",
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore(t),
				idGenerator:         tt.fields.idGenerator,
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			id, got, err := c.AddInstanceGenericOAuthProvider(tt.args.ctx, tt.args.provider)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.id, id)
				assertObjectDetails(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_UpdateInstanceGenericOAuthIDP(t *testing.T) {
	type fields struct {
		eventstore   func(*testing.T) *eventstore.Eventstore
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx      context.Context
		id       string
		provider GenericOAuthProvider
	}
	type res struct {
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"invalid id",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				provider: GenericOAuthProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-SAffg", ""))
				},
			},
		},
		{
			"invalid name",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				id:       "id1",
				provider: GenericOAuthProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-Sf3gh", ""))
				},
			},
		},
		{
			"invalid clientID",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: GenericOAuthProvider{
					Name: "name",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-SHJ3ui", ""))
				},
			},
		},
		{
			"invalid auth endpoint",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: GenericOAuthProvider{
					Name:     "name",
					ClientID: "clientID",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-SVrgh", ""))
				},
			},
		},
		{
			"invalid token endpoint",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: GenericOAuthProvider{
					Name:                  "name",
					ClientID:              "clientID",
					AuthorizationEndpoint: "auth",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-DJKeio", ""))
				},
			},
		},
		{
			"invalid user endpoint",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: GenericOAuthProvider{
					Name:                  "name",
					ClientID:              "clientID",
					AuthorizationEndpoint: "auth",
					TokenEndpoint:         "token",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-ILSJi", ""))
				},
			},
		},
		{
			"invalid id attribute",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: GenericOAuthProvider{
					Name:                  "name",
					ClientID:              "clientID",
					AuthorizationEndpoint: "auth",
					TokenEndpoint:         "token",
					UserEndpoint:          "user",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-JKD3h", ""))
				},
			},
		},
		{
			name: "not found",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: GenericOAuthProvider{
					Name:                  "name",
					ClientID:              "clientID",
					AuthorizationEndpoint: "auth",
					TokenEndpoint:         "token",
					UserEndpoint:          "user",
					IDAttribute:           "idAttribute",
				},
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "no changes",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							instance.NewOAuthIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
								"id1",
								"name",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								"auth",
								"token",
								"user",
								"idAttribute",
								nil,
								true,
								idp.Options{},
							)),
					),
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: GenericOAuthProvider{
					Name:                  "name",
					ClientID:              "clientID",
					AuthorizationEndpoint: "auth",
					TokenEndpoint:         "token",
					UserEndpoint:          "user",
					IDAttribute:           "idAttribute",
					UsePKCE:               true,
				},
			},
			res: res{
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
		{
			name: "change ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							instance.NewOAuthIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
								"id1",
								"name",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								"auth",
								"token",
								"user",
								"idAttribute",
								nil,
								false,
								idp.Options{},
							)),
					),
					expectPush(
						func() eventstore.Command {
							t := true
							event, _ := instance.NewOAuthIDPChangedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
								"id1",
								[]idp.OAuthIDPChanges{
									idp.ChangeOAuthName("new name"),
									idp.ChangeOAuthClientID("clientID2"),
									idp.ChangeOAuthClientSecret(&crypto.CryptoValue{
										CryptoType: crypto.TypeEncryption,
										Algorithm:  "enc",
										KeyID:      "id",
										Crypted:    []byte("newSecret"),
									}),
									idp.ChangeOAuthAuthorizationEndpoint("new auth"),
									idp.ChangeOAuthTokenEndpoint("new token"),
									idp.ChangeOAuthUserEndpoint("new user"),
									idp.ChangeOAuthScopes([]string{"openid", "profile"}),
									idp.ChangeOAuthIDAttribute("newAttribute"),
									idp.ChangeOAuthUsePKCE(true),
									idp.ChangeOAuthOptions(idp.OptionChanges{
										IsCreationAllowed: &t,
										IsLinkingAllowed:  &t,
										IsAutoCreation:    &t,
										IsAutoUpdate:      &t,
									}),
								},
							)
							return event
						}(),
					),
				),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: GenericOAuthProvider{
					Name:                  "new name",
					ClientID:              "clientID2",
					ClientSecret:          "newSecret",
					AuthorizationEndpoint: "new auth",
					TokenEndpoint:         "new token",
					UserEndpoint:          "new user",
					Scopes:                []string{"openid", "profile"},
					IDAttribute:           "newAttribute",
					UsePKCE:               true,
					IDPOptions: idp.Options{
						IsCreationAllowed: true,
						IsLinkingAllowed:  true,
						IsAutoCreation:    true,
						IsAutoUpdate:      true,
					},
				},
			},
			res: res{
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore(t),
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			got, err := c.UpdateInstanceGenericOAuthProvider(tt.args.ctx, tt.args.id, tt.args.provider)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assertObjectDetails(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_AddInstanceGenericOIDCIDP(t *testing.T) {
	type fields struct {
		eventstore   func(*testing.T) *eventstore.Eventstore
		idGenerator  id.Generator
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx      context.Context
		provider GenericOIDCProvider
	}
	type res struct {
		id   string
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"invalid name",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				provider: GenericOIDCProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-Sgtj5", ""))
				},
			},
		},
		{
			"invalid issuer",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: GenericOIDCProvider{
					Name: "name",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-Hz6zj", ""))
				},
			},
		},
		{
			"invalid clientID",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: GenericOIDCProvider{
					Name:   "name",
					Issuer: "issuer",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-fb5jm", ""))
				},
			},
		},
		{
			"invalid clientSecret",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: GenericOIDCProvider{
					Name:     "name",
					Issuer:   "issuer",
					ClientID: "clientID",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-Sfdf4", ""))
				},
			},
		},
		{
			name: "ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectPush(
						instance.NewOIDCIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
							"id1",
							"name",
							"issuer",
							"clientID",
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("clientSecret"),
							},
							nil,
							false,
							true,
							idp.Options{},
						),
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: GenericOIDCProvider{
					Name:         "name",
					Issuer:       "issuer",
					ClientID:     "clientID",
					ClientSecret: "clientSecret",
					UsePKCE:      true,
				},
			},
			res: res{
				id:   "id1",
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
		{
			name: "ok all set",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectPush(
						instance.NewOIDCIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
							"id1",
							"name",
							"issuer",
							"clientID",
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("clientSecret"),
							},
							[]string{openid.ScopeOpenID},
							true,
							true,
							idp.Options{
								IsCreationAllowed: true,
								IsLinkingAllowed:  true,
								IsAutoCreation:    true,
								IsAutoUpdate:      true,
							},
						),
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: GenericOIDCProvider{
					Name:             "name",
					Issuer:           "issuer",
					ClientID:         "clientID",
					ClientSecret:     "clientSecret",
					Scopes:           []string{openid.ScopeOpenID},
					IsIDTokenMapping: true,
					UsePKCE:          true,
					IDPOptions: idp.Options{
						IsCreationAllowed: true,
						IsLinkingAllowed:  true,
						IsAutoCreation:    true,
						IsAutoUpdate:      true,
					},
				},
			},
			res: res{
				id:   "id1",
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore(t),
				idGenerator:         tt.fields.idGenerator,
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			id, got, err := c.AddInstanceGenericOIDCProvider(tt.args.ctx, tt.args.provider)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.id, id)
				assertObjectDetails(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_UpdateInstanceGenericOIDCIDP(t *testing.T) {
	type fields struct {
		eventstore   func(*testing.T) *eventstore.Eventstore
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx      context.Context
		id       string
		provider GenericOIDCProvider
	}
	type res struct {
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"invalid id",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				provider: GenericOIDCProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-SAfd3", ""))
				},
			},
		},
		{
			"invalid name",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				id:       "id1",
				provider: GenericOIDCProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-Dvf4f", ""))
				},
			},
		},
		{
			"invalid issuer",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: GenericOIDCProvider{
					Name: "name",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-BDfr3", ""))
				},
			},
		},
		{
			"invalid clientID",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: GenericOIDCProvider{
					Name:   "name",
					Issuer: "issuer",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-Db3bs", ""))
				},
			},
		},
		{
			name: "not found",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: GenericOIDCProvider{
					Name:     "name",
					Issuer:   "issuer",
					ClientID: "clientID",
				},
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "no changes",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							instance.NewOIDCIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
								"id1",
								"name",
								"issuer",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								nil,
								false,
								false,
								idp.Options{},
							)),
					),
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: GenericOIDCProvider{
					Name:     "name",
					Issuer:   "issuer",
					ClientID: "clientID",
				},
			},
			res: res{
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
		{
			name: "change ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							instance.NewOIDCIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
								"id1",
								"name",
								"issuer",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								nil,
								false,
								false,
								idp.Options{},
							)),
					),
					expectPush(
						func() eventstore.Command {
							t := true
							event, _ := instance.NewOIDCIDPChangedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
								"id1",
								[]idp.OIDCIDPChanges{
									idp.ChangeOIDCName("new name"),
									idp.ChangeOIDCIssuer("new issuer"),
									idp.ChangeOIDCClientID("clientID2"),
									idp.ChangeOIDCClientSecret(&crypto.CryptoValue{
										CryptoType: crypto.TypeEncryption,
										Algorithm:  "enc",
										KeyID:      "id",
										Crypted:    []byte("newSecret"),
									}),
									idp.ChangeOIDCScopes([]string{"openid", "profile"}),
									idp.ChangeOIDCIsIDTokenMapping(true),
									idp.ChangeOIDCUsePKCE(true),
									idp.ChangeOIDCOptions(idp.OptionChanges{
										IsCreationAllowed: &t,
										IsLinkingAllowed:  &t,
										IsAutoCreation:    &t,
										IsAutoUpdate:      &t,
									}),
								},
							)
							return event
						}(),
					),
				),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: GenericOIDCProvider{
					Name:             "new name",
					Issuer:           "new issuer",
					ClientID:         "clientID2",
					ClientSecret:     "newSecret",
					Scopes:           []string{"openid", "profile"},
					IsIDTokenMapping: true,
					UsePKCE:          true,
					IDPOptions: idp.Options{
						IsCreationAllowed: true,
						IsLinkingAllowed:  true,
						IsAutoCreation:    true,
						IsAutoUpdate:      true,
					},
				},
			},
			res: res{
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore(t),
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			got, err := c.UpdateInstanceGenericOIDCProvider(tt.args.ctx, tt.args.id, tt.args.provider)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assertObjectDetails(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_MigrateInstanceGenericOIDCToAzureADProvider(t *testing.T) {
	type fields struct {
		eventstore   func(*testing.T) *eventstore.Eventstore
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx      context.Context
		id       string
		provider AzureADProvider
	}
	type res struct {
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{

		{
			"invalid name",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				provider: AzureADProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-sdf3g", ""))
				},
			},
		},
		{
			"invalid client id",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: AzureADProvider{
					Name: "name",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-Fhbr2", ""))
				},
			},
		},
		{
			"invalid client secret",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: AzureADProvider{
					Name:     "name",
					ClientID: "clientID",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-Dzh3g", ""))
				},
			},
		},
		{
			name: "not found",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: AzureADProvider{
					Name:         "name",
					ClientID:     "clientID",
					ClientSecret: "clientSecret",
				},
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "migrate ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							instance.NewOIDCIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
								"id1",
								"name",
								"issuer",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								nil,
								false,
								false,
								idp.Options{},
							)),
					),
					expectPush(
						func() eventstore.Command {
							event := instance.NewOIDCIDPMigratedAzureADEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
								"id1",
								"name",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								nil,
								"",
								false,
								idp.Options{},
							)
							return event
						}(),
					),
				),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: AzureADProvider{
					Name:         "name",
					ClientID:     "clientID",
					ClientSecret: "clientSecret",
				},
			},
			res: res{
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
		{
			name: "migrate ok full",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							instance.NewOIDCIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
								"id1",
								"name",
								"issuer",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								nil,
								false,
								false,
								idp.Options{},
							)),
					),
					expectPush(
						func() eventstore.Command {
							event := instance.NewOIDCIDPMigratedAzureADEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
								"id1",
								"name",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								[]string{"openid"},
								"tenant",
								true,
								idp.Options{
									IsCreationAllowed: true,
									IsLinkingAllowed:  true,
									IsAutoCreation:    true,
									IsAutoUpdate:      true,
								},
							)
							return event
						}(),
					),
				),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: AzureADProvider{
					Name:          "name",
					ClientID:      "clientID",
					ClientSecret:  "clientSecret",
					Scopes:        []string{"openid"},
					Tenant:        "tenant",
					EmailVerified: true,
					IDPOptions: idp.Options{
						IsCreationAllowed: true,
						IsLinkingAllowed:  true,
						IsAutoCreation:    true,
						IsAutoUpdate:      true,
					},
				},
			},
			res: res{
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore(t),
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			got, err := c.MigrateInstanceGenericOIDCToAzureADProvider(tt.args.ctx, tt.args.id, tt.args.provider)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assertObjectDetails(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_MigrateInstanceOIDCToGoogleIDP(t *testing.T) {
	type fields struct {
		eventstore   func(*testing.T) *eventstore.Eventstore
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx      context.Context
		id       string
		provider GoogleProvider
	}
	type res struct {
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"invalid clientID",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				provider: GoogleProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-D3fvs", ""))
				},
			},
		},
		{
			"invalid clientSecret",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: GoogleProvider{
					ClientID: "clientID",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-W2vqs", ""))
				},
			},
		},
		{
			name: "not found",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: GoogleProvider{
					ClientID:     "clientID",
					ClientSecret: "clientSecret",
				},
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "migrate ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							instance.NewOIDCIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
								"id1",
								"name",
								"issuer",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								nil,
								false,
								false,
								idp.Options{},
							)),
					),
					expectPush(
						instance.NewOIDCIDPMigratedGoogleEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
							"id1",
							"",
							"clientID",
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("clientSecret"),
							},
							nil,
							idp.Options{},
						),
					),
				),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: GoogleProvider{
					ClientID:     "clientID",
					ClientSecret: "clientSecret",
				},
			},
			res: res{
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
		{
			name: "migrate ok full",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							instance.NewOIDCIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
								"id1",
								"name",
								"issuer",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								nil,
								false,
								false,
								idp.Options{},
							)),
					),
					expectPush(
						instance.NewOIDCIDPMigratedGoogleEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
							"id1",
							"",
							"clientID",
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("clientSecret"),
							},
							[]string{"openid"},
							idp.Options{
								IsCreationAllowed: true,
								IsLinkingAllowed:  true,
								IsAutoCreation:    true,
								IsAutoUpdate:      true,
							},
						),
					),
				),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: GoogleProvider{
					ClientID:     "clientID",
					ClientSecret: "clientSecret",
					Scopes:       []string{"openid"},
					IDPOptions: idp.Options{
						IsCreationAllowed: true,
						IsLinkingAllowed:  true,
						IsAutoCreation:    true,
						IsAutoUpdate:      true,
					},
				},
			},
			res: res{
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore(t),
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			got, err := c.MigrateInstanceGenericOIDCToGoogleProvider(tt.args.ctx, tt.args.id, tt.args.provider)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assertObjectDetails(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_AddInstanceAzureADIDP(t *testing.T) {
	type fields struct {
		eventstore   func(*testing.T) *eventstore.Eventstore
		idGenerator  id.Generator
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx      context.Context
		provider AzureADProvider
	}
	type res struct {
		id   string
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"invalid name",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				provider: AzureADProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-sdf3g", ""))
				},
			},
		},
		{
			"invalid client id",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: AzureADProvider{
					Name: "name",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-Fhbr2", ""))
				},
			},
		},
		{
			"invalid client secret",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: AzureADProvider{
					Name:     "name",
					ClientID: "clientID",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-Dzh3g", ""))
				},
			},
		},
		{
			name: "ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectPush(
						instance.NewAzureADIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
							"id1",
							"name",
							"clientID",
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("clientSecret"),
							},
							nil,
							"",
							false,
							idp.Options{},
						),
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: AzureADProvider{
					Name:         "name",
					ClientID:     "clientID",
					ClientSecret: "clientSecret",
				},
			},
			res: res{
				id:   "id1",
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
		{
			name: "ok all set",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectPush(
						instance.NewAzureADIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
							"id1",
							"name",
							"clientID",
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("clientSecret"),
							},
							[]string{"openid"},
							"tenant",
							true,
							idp.Options{
								IsCreationAllowed: true,
								IsLinkingAllowed:  true,
								IsAutoCreation:    true,
								IsAutoUpdate:      true,
							},
						),
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: AzureADProvider{
					Name:          "name",
					ClientID:      "clientID",
					ClientSecret:  "clientSecret",
					Scopes:        []string{"openid"},
					Tenant:        "tenant",
					EmailVerified: true,
					IDPOptions: idp.Options{
						IsCreationAllowed: true,
						IsLinkingAllowed:  true,
						IsAutoCreation:    true,
						IsAutoUpdate:      true,
					},
				},
			},
			res: res{
				id:   "id1",
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore(t),
				idGenerator:         tt.fields.idGenerator,
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			id, got, err := c.AddInstanceAzureADProvider(tt.args.ctx, tt.args.provider)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.id, id)
				assertObjectDetails(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_UpdateInstanceAzureADIDP(t *testing.T) {
	type fields struct {
		eventstore   func(*testing.T) *eventstore.Eventstore
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx      context.Context
		id       string
		provider AzureADProvider
	}
	type res struct {
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"invalid id",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				provider: AzureADProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-SAgh2", ""))
				},
			},
		},
		{
			"invalid name",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				id:       "id1",
				provider: AzureADProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-fh3h1", ""))
				},
			},
		},
		{
			"invalid client id",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: AzureADProvider{
					Name: "name",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-dmitg", ""))
				},
			},
		},
		{
			name: "not found",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: AzureADProvider{
					Name:     "name",
					ClientID: "clientID",
				},
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "no changes",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							instance.NewAzureADIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
								"id1",
								"name",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								nil,
								"",
								false,
								idp.Options{},
							)),
					),
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: AzureADProvider{
					Name:     "name",
					ClientID: "clientID",
				},
			},
			res: res{
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
		{
			name: "change ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							instance.NewAzureADIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
								"id1",
								"name",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								nil,
								"",
								false,
								idp.Options{},
							)),
					),
					expectPush(
						func() eventstore.Command {
							t := true
							event, _ := instance.NewAzureADIDPChangedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
								"id1",
								[]idp.AzureADIDPChanges{
									idp.ChangeAzureADName("new name"),
									idp.ChangeAzureADClientID("new clientID"),
									idp.ChangeAzureADClientSecret(&crypto.CryptoValue{
										CryptoType: crypto.TypeEncryption,
										Algorithm:  "enc",
										KeyID:      "id",
										Crypted:    []byte("new clientSecret"),
									}),
									idp.ChangeAzureADScopes([]string{"openid", "profile"}),
									idp.ChangeAzureADTenant("new tenant"),
									idp.ChangeAzureADIsEmailVerified(true),
									idp.ChangeAzureADOptions(idp.OptionChanges{
										IsCreationAllowed: &t,
										IsLinkingAllowed:  &t,
										IsAutoCreation:    &t,
										IsAutoUpdate:      &t,
									}),
								},
							)
							return event
						}(),
					),
				),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: AzureADProvider{
					Name:          "new name",
					ClientID:      "new clientID",
					ClientSecret:  "new clientSecret",
					Scopes:        []string{"openid", "profile"},
					Tenant:        "new tenant",
					EmailVerified: true,
					IDPOptions: idp.Options{
						IsCreationAllowed: true,
						IsLinkingAllowed:  true,
						IsAutoCreation:    true,
						IsAutoUpdate:      true,
					},
				},
			},
			res: res{
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore(t),
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			got, err := c.UpdateInstanceAzureADProvider(tt.args.ctx, tt.args.id, tt.args.provider)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assertObjectDetails(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_AddInstanceGitHubIDP(t *testing.T) {
	type fields struct {
		eventstore   func(*testing.T) *eventstore.Eventstore
		idGenerator  id.Generator
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx      context.Context
		provider GitHubProvider
	}
	type res struct {
		id   string
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"invalid client id",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				provider: GitHubProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-Jdsgf", ""))
				},
			},
		},
		{
			"invalid client secret",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: GitHubProvider{
					ClientID: "clientID",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-dsgz3", ""))
				},
			},
		},
		{
			name: "ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectPush(
						instance.NewGitHubIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
							"id1",
							"",
							"clientID",
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("clientSecret"),
							},
							nil,
							idp.Options{},
						),
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: GitHubProvider{
					ClientID:     "clientID",
					ClientSecret: "clientSecret",
				},
			},
			res: res{
				id:   "id1",
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
		{
			name: "ok all set",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectPush(
						instance.NewGitHubIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
							"id1",
							"name",
							"clientID",
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("clientSecret"),
							},
							[]string{"openid"},
							idp.Options{
								IsCreationAllowed: true,
								IsLinkingAllowed:  true,
								IsAutoCreation:    true,
								IsAutoUpdate:      true,
							},
						),
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: GitHubProvider{
					Name:         "name",
					ClientID:     "clientID",
					ClientSecret: "clientSecret",
					Scopes:       []string{"openid"},
					IDPOptions: idp.Options{
						IsCreationAllowed: true,
						IsLinkingAllowed:  true,
						IsAutoCreation:    true,
						IsAutoUpdate:      true,
					},
				},
			},
			res: res{
				id:   "id1",
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore(t),
				idGenerator:         tt.fields.idGenerator,
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			id, got, err := c.AddInstanceGitHubProvider(tt.args.ctx, tt.args.provider)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.id, id)
				assertObjectDetails(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_UpdateInstanceGitHubIDP(t *testing.T) {
	type fields struct {
		eventstore   func(*testing.T) *eventstore.Eventstore
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx      context.Context
		id       string
		provider GitHubProvider
	}
	type res struct {
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"invalid id",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				provider: GitHubProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-sdf4h", ""))
				},
			},
		},
		{
			"invalid client id",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				id:       "id1",
				provider: GitHubProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-fdh5z", ""))
				},
			},
		},
		{
			name: "not found",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: GitHubProvider{
					ClientID: "clientID",
				},
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "no changes",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							instance.NewGitHubIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
								"id1",
								"",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								nil,
								idp.Options{},
							)),
					),
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: GitHubProvider{
					ClientID: "clientID",
				},
			},
			res: res{
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
		{
			name: "change ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							instance.NewGitHubIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
								"id1",
								"name",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								nil,
								idp.Options{},
							)),
					),
					expectPush(
						func() eventstore.Command {
							t := true
							event, _ := instance.NewGitHubIDPChangedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
								"id1",
								[]idp.GitHubIDPChanges{
									idp.ChangeGitHubName("new name"),
									idp.ChangeGitHubClientID("new clientID"),
									idp.ChangeGitHubClientSecret(&crypto.CryptoValue{
										CryptoType: crypto.TypeEncryption,
										Algorithm:  "enc",
										KeyID:      "id",
										Crypted:    []byte("new clientSecret"),
									}),
									idp.ChangeGitHubScopes([]string{"openid", "profile"}),
									idp.ChangeGitHubOptions(idp.OptionChanges{
										IsCreationAllowed: &t,
										IsLinkingAllowed:  &t,
										IsAutoCreation:    &t,
										IsAutoUpdate:      &t,
									}),
								},
							)
							return event
						}(),
					),
				),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: GitHubProvider{
					Name:         "new name",
					ClientID:     "new clientID",
					ClientSecret: "new clientSecret",
					Scopes:       []string{"openid", "profile"},
					IDPOptions: idp.Options{
						IsCreationAllowed: true,
						IsLinkingAllowed:  true,
						IsAutoCreation:    true,
						IsAutoUpdate:      true,
					},
				},
			},
			res: res{
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore(t),
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			got, err := c.UpdateInstanceGitHubProvider(tt.args.ctx, tt.args.id, tt.args.provider)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assertObjectDetails(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_AddInstanceGitHubEnterpriseIDP(t *testing.T) {
	type fields struct {
		eventstore   func(*testing.T) *eventstore.Eventstore
		idGenerator  id.Generator
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx      context.Context
		provider GitHubEnterpriseProvider
	}
	type res struct {
		id   string
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"invalid name",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				provider: GitHubEnterpriseProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-Dg4td", ""))
				},
			},
		},
		{
			"invalid clientID",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: GitHubEnterpriseProvider{
					Name: "name",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-dgj53", ""))
				},
			},
		},
		{
			"invalid clientSecret",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: GitHubEnterpriseProvider{
					Name:     "name",
					ClientID: "clientID",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-Ghjjs", ""))
				},
			},
		},
		{
			"invalid auth endpoint",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: GitHubEnterpriseProvider{
					Name:         "name",
					ClientID:     "clientID",
					ClientSecret: "clientSecret",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-sani2", ""))
				},
			},
		},
		{
			"invalid token endpoint",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: GitHubEnterpriseProvider{
					Name:                  "name",
					ClientID:              "clientID",
					ClientSecret:          "clientSecret",
					AuthorizationEndpoint: "auth",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-agj42", ""))
				},
			},
		},
		{
			"invalid user endpoint",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: GitHubEnterpriseProvider{
					Name:                  "name",
					ClientID:              "clientID",
					ClientSecret:          "clientSecret",
					AuthorizationEndpoint: "auth",
					TokenEndpoint:         "token",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-sd5hn", ""))
				},
			},
		},
		{
			name: "ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectPush(
						instance.NewGitHubEnterpriseIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
							"id1",
							"name",
							"clientID",
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("clientSecret"),
							},
							"auth",
							"token",
							"user",
							nil,
							idp.Options{},
						),
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: GitHubEnterpriseProvider{
					Name:                  "name",
					ClientID:              "clientID",
					ClientSecret:          "clientSecret",
					AuthorizationEndpoint: "auth",
					TokenEndpoint:         "token",
					UserEndpoint:          "user",
				},
			},
			res: res{
				id:   "id1",
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
		{
			name: "ok all set",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectPush(
						instance.NewGitHubEnterpriseIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
							"id1",
							"name",
							"clientID",
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("clientSecret"),
							},
							"auth",
							"token",
							"user",
							[]string{"user"},
							idp.Options{
								IsCreationAllowed: true,
								IsLinkingAllowed:  true,
								IsAutoCreation:    true,
								IsAutoUpdate:      true,
							},
						),
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: GitHubEnterpriseProvider{
					Name:                  "name",
					ClientID:              "clientID",
					ClientSecret:          "clientSecret",
					AuthorizationEndpoint: "auth",
					TokenEndpoint:         "token",
					UserEndpoint:          "user",
					Scopes:                []string{"user"},
					IDPOptions: idp.Options{
						IsCreationAllowed: true,
						IsLinkingAllowed:  true,
						IsAutoCreation:    true,
						IsAutoUpdate:      true,
					},
				},
			},
			res: res{
				id:   "id1",
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore(t),
				idGenerator:         tt.fields.idGenerator,
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			id, got, err := c.AddInstanceGitHubEnterpriseProvider(tt.args.ctx, tt.args.provider)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.id, id)
				assertObjectDetails(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_UpdateInstanceGitHubEnterpriseIDP(t *testing.T) {
	type fields struct {
		eventstore   func(*testing.T) *eventstore.Eventstore
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx      context.Context
		id       string
		provider GitHubEnterpriseProvider
	}
	type res struct {
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"invalid id",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				provider: GitHubEnterpriseProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-sdfh3", ""))
				},
			},
		},
		{
			"invalid name",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				id:       "id1",
				provider: GitHubEnterpriseProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-shj42", ""))
				},
			},
		},
		{
			"invalid clientID",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: GitHubEnterpriseProvider{
					Name: "name",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-sdh73", ""))
				},
			},
		},
		{
			"invalid auth endpoint",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: GitHubEnterpriseProvider{
					Name:     "name",
					ClientID: "clientID",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-acx2w", ""))
				},
			},
		},
		{
			"invalid token endpoint",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: GitHubEnterpriseProvider{
					Name:                  "name",
					ClientID:              "clientID",
					AuthorizationEndpoint: "auth",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-dgj6q", ""))
				},
			},
		},
		{
			"invalid user endpoint",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: GitHubEnterpriseProvider{
					Name:                  "name",
					ClientID:              "clientID",
					AuthorizationEndpoint: "auth",
					TokenEndpoint:         "token",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-ybj62", ""))
				},
			},
		},
		{
			name: "not found",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: GitHubEnterpriseProvider{
					Name:                  "name",
					ClientID:              "clientID",
					AuthorizationEndpoint: "auth",
					TokenEndpoint:         "token",
					UserEndpoint:          "user",
				},
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "no changes",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							instance.NewGitHubEnterpriseIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
								"id1",
								"name",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								"auth",
								"token",
								"user",
								nil,
								idp.Options{},
							)),
					),
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: GitHubEnterpriseProvider{
					Name:                  "name",
					ClientID:              "clientID",
					AuthorizationEndpoint: "auth",
					TokenEndpoint:         "token",
					UserEndpoint:          "user",
				},
			},
			res: res{
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
		{
			name: "change ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							instance.NewGitHubEnterpriseIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
								"id1",
								"name",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								"auth",
								"token",
								"user",
								nil,
								idp.Options{},
							)),
					),
					expectPush(
						func() eventstore.Command {
							t := true
							event, _ := instance.NewGitHubEnterpriseIDPChangedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
								"id1",
								[]idp.GitHubEnterpriseIDPChanges{
									idp.ChangeGitHubEnterpriseName("new name"),
									idp.ChangeGitHubEnterpriseClientID("clientID2"),
									idp.ChangeGitHubEnterpriseClientSecret(&crypto.CryptoValue{
										CryptoType: crypto.TypeEncryption,
										Algorithm:  "enc",
										KeyID:      "id",
										Crypted:    []byte("newSecret"),
									}),
									idp.ChangeGitHubEnterpriseAuthorizationEndpoint("new auth"),
									idp.ChangeGitHubEnterpriseTokenEndpoint("new token"),
									idp.ChangeGitHubEnterpriseUserEndpoint("new user"),
									idp.ChangeGitHubEnterpriseScopes([]string{"openid", "profile"}),
									idp.ChangeGitHubEnterpriseOptions(idp.OptionChanges{
										IsCreationAllowed: &t,
										IsLinkingAllowed:  &t,
										IsAutoCreation:    &t,
										IsAutoUpdate:      &t,
									}),
								},
							)
							return event
						}(),
					),
				),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: GitHubEnterpriseProvider{
					Name:                  "new name",
					ClientID:              "clientID2",
					ClientSecret:          "newSecret",
					AuthorizationEndpoint: "new auth",
					TokenEndpoint:         "new token",
					UserEndpoint:          "new user",
					Scopes:                []string{"openid", "profile"},
					IDPOptions: idp.Options{
						IsCreationAllowed: true,
						IsLinkingAllowed:  true,
						IsAutoCreation:    true,
						IsAutoUpdate:      true,
					},
				},
			},
			res: res{
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore(t),
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			got, err := c.UpdateInstanceGitHubEnterpriseProvider(tt.args.ctx, tt.args.id, tt.args.provider)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assertObjectDetails(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_AddInstanceGitLabIDP(t *testing.T) {
	type fields struct {
		eventstore   func(*testing.T) *eventstore.Eventstore
		idGenerator  id.Generator
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx      context.Context
		provider GitLabProvider
	}
	type res struct {
		id   string
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"invalid clientID",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				provider: GitLabProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-adsg2", ""))
				},
			},
		},
		{
			"invalid clientSecret",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: GitLabProvider{
					ClientID: "clientID",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-GD1j2", ""))
				},
			},
		},
		{
			name: "ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectPush(
						instance.NewGitLabIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
							"id1",
							"",
							"clientID",
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("clientSecret"),
							},
							nil,
							idp.Options{},
						),
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: GitLabProvider{
					ClientID:     "clientID",
					ClientSecret: "clientSecret",
				},
			},
			res: res{
				id:   "id1",
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
		{
			name: "ok all set",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectPush(
						instance.NewGitLabIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
							"id1",
							"",
							"clientID",
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("clientSecret"),
							},
							[]string{"openid"},
							idp.Options{
								IsCreationAllowed: true,
								IsLinkingAllowed:  true,
								IsAutoCreation:    true,
								IsAutoUpdate:      true,
							},
						),
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: GitLabProvider{
					ClientID:     "clientID",
					ClientSecret: "clientSecret",
					Scopes:       []string{"openid"},
					IDPOptions: idp.Options{
						IsCreationAllowed: true,
						IsLinkingAllowed:  true,
						IsAutoCreation:    true,
						IsAutoUpdate:      true,
					},
				},
			},
			res: res{
				id:   "id1",
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore(t),
				idGenerator:         tt.fields.idGenerator,
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			id, got, err := c.AddInstanceGitLabProvider(tt.args.ctx, tt.args.provider)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.id, id)
				assertObjectDetails(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_UpdateInstanceGitLabIDP(t *testing.T) {
	type fields struct {
		eventstore   func(*testing.T) *eventstore.Eventstore
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx      context.Context
		id       string
		provider GitLabProvider
	}
	type res struct {
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"invalid id",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				provider: GitLabProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-HJK91", ""))
				},
			},
		},
		{
			"invalid clientID",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				id:       "id1",
				provider: GitLabProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-D12t6", ""))
				},
			},
		},
		{
			name: "not found",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: GitLabProvider{
					ClientID: "clientID",
				},
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "no changes",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							instance.NewGitLabIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
								"id1",
								"",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								nil,
								idp.Options{},
							)),
					),
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: GitLabProvider{
					ClientID: "clientID",
				},
			},
			res: res{
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
		{
			name: "change ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							instance.NewGitLabIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
								"id1",
								"",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								nil,
								idp.Options{},
							)),
					),
					expectPush(
						func() eventstore.Command {
							t := true
							event, _ := instance.NewGitLabIDPChangedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
								"id1",
								[]idp.GitLabIDPChanges{
									idp.ChangeGitLabClientID("clientID2"),
									idp.ChangeGitLabClientSecret(&crypto.CryptoValue{
										CryptoType: crypto.TypeEncryption,
										Algorithm:  "enc",
										KeyID:      "id",
										Crypted:    []byte("newSecret"),
									}),
									idp.ChangeGitLabScopes([]string{"openid", "profile"}),
									idp.ChangeGitLabOptions(idp.OptionChanges{
										IsCreationAllowed: &t,
										IsLinkingAllowed:  &t,
										IsAutoCreation:    &t,
										IsAutoUpdate:      &t,
									}),
								},
							)
							return event
						}(),
					),
				),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: GitLabProvider{
					ClientID:     "clientID2",
					ClientSecret: "newSecret",
					Scopes:       []string{"openid", "profile"},
					IDPOptions: idp.Options{
						IsCreationAllowed: true,
						IsLinkingAllowed:  true,
						IsAutoCreation:    true,
						IsAutoUpdate:      true,
					},
				},
			},
			res: res{
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore(t),
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			got, err := c.UpdateInstanceGitLabProvider(tt.args.ctx, tt.args.id, tt.args.provider)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assertObjectDetails(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_AddInstanceGitLabSelfHostedIDP(t *testing.T) {
	type fields struct {
		eventstore   func(*testing.T) *eventstore.Eventstore
		idGenerator  id.Generator
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx      context.Context
		provider GitLabSelfHostedProvider
	}
	type res struct {
		id   string
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"invalid name",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				provider: GitLabSelfHostedProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-jw4ZT", ""))
				},
			},
		},
		{
			"invalid issuer",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: GitLabSelfHostedProvider{
					Name: "name",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-AST4S", ""))
				},
			},
		},
		{
			"invalid clientID",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: GitLabSelfHostedProvider{
					Name:   "name",
					Issuer: "issuer",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-DBZHJ", ""))
				},
			},
		},
		{
			"invalid clientSecret",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: GitLabSelfHostedProvider{
					Name:     "name",
					Issuer:   "issuer",
					ClientID: "clientID",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-SDGJ4", ""))
				},
			},
		},
		{
			name: "ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectPush(
						instance.NewGitLabSelfHostedIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
							"id1",
							"name",
							"issuer",
							"clientID",
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("clientSecret"),
							},
							nil,
							idp.Options{},
						),
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: GitLabSelfHostedProvider{
					Name:         "name",
					Issuer:       "issuer",
					ClientID:     "clientID",
					ClientSecret: "clientSecret",
				},
			},
			res: res{
				id:   "id1",
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
		{
			name: "ok all set",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectPush(
						instance.NewGitLabSelfHostedIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
							"id1",
							"name",
							"issuer",
							"clientID",
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("clientSecret"),
							},
							[]string{"openid"},
							idp.Options{
								IsCreationAllowed: true,
								IsLinkingAllowed:  true,
								IsAutoCreation:    true,
								IsAutoUpdate:      true,
							},
						),
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: GitLabSelfHostedProvider{
					Name:         "name",
					Issuer:       "issuer",
					ClientID:     "clientID",
					ClientSecret: "clientSecret",
					Scopes:       []string{"openid"},
					IDPOptions: idp.Options{
						IsCreationAllowed: true,
						IsLinkingAllowed:  true,
						IsAutoCreation:    true,
						IsAutoUpdate:      true,
					},
				},
			},
			res: res{
				id:   "id1",
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore(t),
				idGenerator:         tt.fields.idGenerator,
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			id, got, err := c.AddInstanceGitLabSelfHostedProvider(tt.args.ctx, tt.args.provider)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.id, id)
				assertObjectDetails(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_UpdateInstanceGitLabSelfHostedIDP(t *testing.T) {
	type fields struct {
		eventstore   func(*testing.T) *eventstore.Eventstore
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx      context.Context
		id       string
		provider GitLabSelfHostedProvider
	}
	type res struct {
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"invalid id",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				provider: GitLabSelfHostedProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-SAFG4", ""))
				},
			},
		},
		{
			"invalid name",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				id:       "id1",
				provider: GitLabSelfHostedProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-DG4H", ""))
				},
			},
		},
		{
			"invalid issuer",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: GitLabSelfHostedProvider{
					Name: "name",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-SD4eb", ""))
				},
			},
		},
		{
			"invalid clientID",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: GitLabSelfHostedProvider{
					Name:   "name",
					Issuer: "issuer",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-GHWE3", ""))
				},
			},
		},
		{
			name: "not found",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: GitLabSelfHostedProvider{
					Name:     "name",
					Issuer:   "issuer",
					ClientID: "clientID",
				},
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "no changes",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							instance.NewGitLabSelfHostedIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
								"id1",
								"name",
								"issuer",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								nil,
								idp.Options{},
							)),
					),
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: GitLabSelfHostedProvider{
					Name:     "name",
					Issuer:   "issuer",
					ClientID: "clientID",
				},
			},
			res: res{
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
		{
			name: "change ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							instance.NewGitLabSelfHostedIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
								"id1",
								"name",
								"issuer",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								nil,
								idp.Options{},
							)),
					),
					expectPush(
						func() eventstore.Command {
							t := true
							event, _ := instance.NewGitLabSelfHostedIDPChangedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
								"id1",
								[]idp.GitLabSelfHostedIDPChanges{
									idp.ChangeGitLabSelfHostedClientID("clientID2"),
									idp.ChangeGitLabSelfHostedIssuer("newIssuer"),
									idp.ChangeGitLabSelfHostedName("newName"),
									idp.ChangeGitLabSelfHostedClientSecret(&crypto.CryptoValue{
										CryptoType: crypto.TypeEncryption,
										Algorithm:  "enc",
										KeyID:      "id",
										Crypted:    []byte("newSecret"),
									}),
									idp.ChangeGitLabSelfHostedScopes([]string{"openid", "profile"}),
									idp.ChangeGitLabSelfHostedOptions(idp.OptionChanges{
										IsCreationAllowed: &t,
										IsLinkingAllowed:  &t,
										IsAutoCreation:    &t,
										IsAutoUpdate:      &t,
									}),
								},
							)
							return event
						}(),
					),
				),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: GitLabSelfHostedProvider{
					Issuer:       "newIssuer",
					Name:         "newName",
					ClientID:     "clientID2",
					ClientSecret: "newSecret",
					Scopes:       []string{"openid", "profile"},
					IDPOptions: idp.Options{
						IsCreationAllowed: true,
						IsLinkingAllowed:  true,
						IsAutoCreation:    true,
						IsAutoUpdate:      true,
					},
				},
			},
			res: res{
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore(t),
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			got, err := c.UpdateInstanceGitLabSelfHostedProvider(tt.args.ctx, tt.args.id, tt.args.provider)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assertObjectDetails(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_AddInstanceGoogleIDP(t *testing.T) {
	type fields struct {
		eventstore   func(*testing.T) *eventstore.Eventstore
		idGenerator  id.Generator
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx      context.Context
		provider GoogleProvider
	}
	type res struct {
		id   string
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"invalid clientID",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				provider: GoogleProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-D3fvs", ""))
				},
			},
		},
		{
			"invalid clientSecret",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: GoogleProvider{
					ClientID: "clientID",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-W2vqs", ""))
				},
			},
		},
		{
			name: "ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectPush(
						instance.NewGoogleIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
							"id1",
							"",
							"clientID",
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("clientSecret"),
							},
							nil,
							idp.Options{},
						),
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: GoogleProvider{
					ClientID:     "clientID",
					ClientSecret: "clientSecret",
				},
			},
			res: res{
				id:   "id1",
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
		{
			name: "ok all set",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectPush(
						instance.NewGoogleIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
							"id1",
							"",
							"clientID",
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("clientSecret"),
							},
							[]string{"openid"},
							idp.Options{
								IsCreationAllowed: true,
								IsLinkingAllowed:  true,
								IsAutoCreation:    true,
								IsAutoUpdate:      true,
							},
						),
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: GoogleProvider{
					ClientID:     "clientID",
					ClientSecret: "clientSecret",
					Scopes:       []string{"openid"},
					IDPOptions: idp.Options{
						IsCreationAllowed: true,
						IsLinkingAllowed:  true,
						IsAutoCreation:    true,
						IsAutoUpdate:      true,
					},
				},
			},
			res: res{
				id:   "id1",
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore(t),
				idGenerator:         tt.fields.idGenerator,
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			id, got, err := c.AddInstanceGoogleProvider(tt.args.ctx, tt.args.provider)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.id, id)
				assertObjectDetails(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_UpdateInstanceGoogleIDP(t *testing.T) {
	type fields struct {
		eventstore   func(*testing.T) *eventstore.Eventstore
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx      context.Context
		id       string
		provider GoogleProvider
	}
	type res struct {
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"invalid id",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				provider: GoogleProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-S32t1", ""))
				},
			},
		},
		{
			"invalid clientID",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				id:       "id1",
				provider: GoogleProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-ds432", ""))
				},
			},
		},
		{
			name: "not found",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: GoogleProvider{
					ClientID: "clientID",
				},
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "no changes",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							instance.NewGoogleIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
								"id1",
								"",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								nil,
								idp.Options{},
							)),
					),
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: GoogleProvider{
					ClientID: "clientID",
				},
			},
			res: res{
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
		{
			name: "change ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							instance.NewGoogleIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
								"id1",
								"",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								nil,
								idp.Options{},
							)),
					),
					expectPush(
						func() eventstore.Command {
							t := true
							event, _ := instance.NewGoogleIDPChangedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
								"id1",
								[]idp.GoogleIDPChanges{
									idp.ChangeGoogleClientID("clientID2"),
									idp.ChangeGoogleClientSecret(&crypto.CryptoValue{
										CryptoType: crypto.TypeEncryption,
										Algorithm:  "enc",
										KeyID:      "id",
										Crypted:    []byte("newSecret"),
									}),
									idp.ChangeGoogleScopes([]string{"openid", "profile"}),
									idp.ChangeGoogleOptions(idp.OptionChanges{
										IsCreationAllowed: &t,
										IsLinkingAllowed:  &t,
										IsAutoCreation:    &t,
										IsAutoUpdate:      &t,
									}),
								},
							)
							return event
						}(),
					),
				),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: GoogleProvider{
					ClientID:     "clientID2",
					ClientSecret: "newSecret",
					Scopes:       []string{"openid", "profile"},
					IDPOptions: idp.Options{
						IsCreationAllowed: true,
						IsLinkingAllowed:  true,
						IsAutoCreation:    true,
						IsAutoUpdate:      true,
					},
				},
			},
			res: res{
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore(t),
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			got, err := c.UpdateInstanceGoogleProvider(tt.args.ctx, tt.args.id, tt.args.provider)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assertObjectDetails(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_AddInstanceLDAPIDP(t *testing.T) {
	type fields struct {
		eventstore   func(*testing.T) *eventstore.Eventstore
		idGenerator  id.Generator
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx      context.Context
		provider LDAPProvider
	}
	type res struct {
		id   string
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"invalid name",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				provider: LDAPProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-SAfdd", ""))
				},
			},
		},
		{
			"invalid baseDN",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: LDAPProvider{
					Name: "name",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-sv31s", ""))
				},
			},
		},
		{
			"invalid bindDN",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: LDAPProvider{
					Name:   "name",
					BaseDN: "baseDN",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-sdgf4", ""))
				},
			},
		},
		{
			"invalid bindPassword",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: LDAPProvider{
					Name:   "name",
					BindDN: "binddn",
					BaseDN: "baseDN",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-AEG2w", ""))
				},
			},
		},
		{
			"invalid userBase",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: LDAPProvider{
					Name:         "name",
					BindDN:       "binddn",
					BaseDN:       "baseDN",
					BindPassword: "password",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-SAD5n", ""))
				},
			},
		},
		{
			"invalid servers",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: LDAPProvider{
					Name:         "name",
					BindDN:       "binddn",
					BaseDN:       "baseDN",
					BindPassword: "password",
					UserBase:     "user",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-SAx905n", ""))
				},
			},
		},
		{
			"invalid userObjectClasses",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: LDAPProvider{
					Name:         "name",
					Servers:      []string{"server"},
					BindDN:       "binddn",
					BaseDN:       "baseDN",
					BindPassword: "password",
					UserBase:     "user",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-S1x905n", ""))
				},
			},
		},
		{
			"invalid userFilters",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: LDAPProvider{
					Name:              "name",
					Servers:           []string{"server"},
					BindDN:            "binddn",
					BaseDN:            "baseDN",
					BindPassword:      "password",
					UserBase:          "user",
					UserObjectClasses: []string{"object"},
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-aAx905n", ""))
				},
			},
		},
		{
			"invalid rootCA",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: LDAPProvider{
					Name:              "name",
					Servers:           []string{"server"},
					StartTLS:          false,
					BaseDN:            "baseDN",
					BindDN:            "dn",
					BindPassword:      "password",
					UserBase:          "user",
					UserObjectClasses: []string{"object"},
					UserFilters:       []string{"filter"},
					Timeout:           time.Second * 30,
					RootCA:            []byte("certificate"),
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-cwqVVdBwKt", "Errors.Invalid.Argument"))
				},
			},
		},
		{
			name: "ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectPush(
						instance.NewLDAPIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
							"id1",
							"name",
							[]string{"server"},
							false,
							"baseDN",
							"dn",
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("password"),
							},
							"user",
							[]string{"object"},
							[]string{"filter"},
							time.Second*30,
							nil,
							idp.LDAPAttributes{},
							idp.Options{},
						),
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: LDAPProvider{
					Name:              "name",
					Servers:           []string{"server"},
					StartTLS:          false,
					BaseDN:            "baseDN",
					BindDN:            "dn",
					BindPassword:      "password",
					UserBase:          "user",
					UserObjectClasses: []string{"object"},
					UserFilters:       []string{"filter"},
					Timeout:           time.Second * 30,
				},
			},
			res: res{
				id:   "id1",
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
		{
			name: "ok all set",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectPush(
						instance.NewLDAPIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
							"id1",
							"name",
							[]string{"server"},
							false,
							"baseDN",
							"dn",
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("password"),
							},
							"user",
							[]string{"object"},
							[]string{"filter"},
							time.Second*30,
							validLDAPRootCA,
							idp.LDAPAttributes{
								IDAttribute:                "id",
								FirstNameAttribute:         "firstName",
								LastNameAttribute:          "lastName",
								DisplayNameAttribute:       "displayName",
								NickNameAttribute:          "nickName",
								PreferredUsernameAttribute: "preferredUsername",
								EmailAttribute:             "email",
								EmailVerifiedAttribute:     "emailVerified",
								PhoneAttribute:             "phone",
								PhoneVerifiedAttribute:     "phoneVerified",
								PreferredLanguageAttribute: "preferredLanguage",
								AvatarURLAttribute:         "avatarURL",
								ProfileAttribute:           "profile",
							},
							idp.Options{
								IsCreationAllowed: true,
								IsLinkingAllowed:  true,
								IsAutoCreation:    true,
								IsAutoUpdate:      true,
							},
						),
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: LDAPProvider{
					Name:              "name",
					Servers:           []string{"server"},
					StartTLS:          false,
					BaseDN:            "baseDN",
					BindDN:            "dn",
					BindPassword:      "password",
					UserBase:          "user",
					UserObjectClasses: []string{"object"},
					UserFilters:       []string{"filter"},
					Timeout:           time.Second * 30,
					RootCA:            validLDAPRootCA,
					LDAPAttributes: idp.LDAPAttributes{
						IDAttribute:                "id",
						FirstNameAttribute:         "firstName",
						LastNameAttribute:          "lastName",
						DisplayNameAttribute:       "displayName",
						NickNameAttribute:          "nickName",
						PreferredUsernameAttribute: "preferredUsername",
						EmailAttribute:             "email",
						EmailVerifiedAttribute:     "emailVerified",
						PhoneAttribute:             "phone",
						PhoneVerifiedAttribute:     "phoneVerified",
						PreferredLanguageAttribute: "preferredLanguage",
						AvatarURLAttribute:         "avatarURL",
						ProfileAttribute:           "profile",
					},
					IDPOptions: idp.Options{
						IsCreationAllowed: true,
						IsLinkingAllowed:  true,
						IsAutoCreation:    true,
						IsAutoUpdate:      true,
					},
				},
			},
			res: res{
				id:   "id1",
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore(t),
				idGenerator:         tt.fields.idGenerator,
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			id, got, err := c.AddInstanceLDAPProvider(tt.args.ctx, tt.args.provider)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.id, id)
				assertObjectDetails(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_UpdateInstanceLDAPIDP(t *testing.T) {
	type fields struct {
		eventstore   func(*testing.T) *eventstore.Eventstore
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx      context.Context
		id       string
		provider LDAPProvider
	}
	type res struct {
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"invalid id",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				provider: LDAPProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-Dgdbs", ""))
				},
			},
		},
		{
			"invalid name",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				id:       "id1",
				provider: LDAPProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-Sffgd", ""))
				},
			},
		},
		{
			"invalid baseDN",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: LDAPProvider{
					Name: "name",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-vb3ss", ""))
				},
			},
		},
		{
			"invalid bindDN",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: LDAPProvider{
					Name:   "name",
					BaseDN: "baseDN",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-hbere", ""))
				},
			},
		},
		{
			"invalid userbase",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: LDAPProvider{
					Name:   "name",
					BaseDN: "baseDN",
					BindDN: "bindDN",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-DG45z", ""))
				},
			},
		},
		{
			"invalid servers",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: LDAPProvider{
					Name:     "name",
					BaseDN:   "baseDN",
					BindDN:   "bindDN",
					UserBase: "user",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-SAx945n", ""))
				},
			},
		},
		{
			"invalid userObjectClasses",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: LDAPProvider{
					Name:     "name",
					Servers:  []string{"server"},
					BaseDN:   "baseDN",
					BindDN:   "bindDN",
					UserBase: "user",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-S1x605n", ""))
				},
			},
		},
		{
			"invalid userFilters",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: LDAPProvider{
					Name:              "name",
					Servers:           []string{"server"},
					BaseDN:            "baseDN",
					BindDN:            "bindDN",
					UserBase:          "user",
					UserObjectClasses: []string{"object"},
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-aAx901n", ""))
				},
			},
		},
		{
			"invalid rootCA",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: LDAPProvider{
					Name:              "name",
					Servers:           []string{"server"},
					BaseDN:            "baseDN",
					BindDN:            "binddn",
					BindPassword:      "password",
					UserBase:          "user",
					UserObjectClasses: []string{"object"},
					UserFilters:       []string{"filter"},
					RootCA:            []byte("certificate"),
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-cwqVVdBwKt", "Errors.Invalid.Argument"))
				},
			},
		},
		{
			name: "not found",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: LDAPProvider{
					Name:              "name",
					Servers:           []string{"server"},
					BaseDN:            "baseDN",
					BindDN:            "binddn",
					BindPassword:      "password",
					UserBase:          "user",
					UserObjectClasses: []string{"object"},
					UserFilters:       []string{"filter"},
				},
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowNotFound(nil, "INST-ASF3F", ""))
				},
			},
		},
		{
			name: "no changes",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							instance.NewLDAPIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
								"id1",
								"name",
								[]string{"server"},
								false,
								"basedn",
								"binddn",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("password"),
								},
								"user",
								[]string{"object"},
								[]string{"filter"},
								time.Second*30,
								validLDAPRootCA,
								idp.LDAPAttributes{},
								idp.Options{},
							)),
					),
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: LDAPProvider{
					Name:              "name",
					Servers:           []string{"server"},
					StartTLS:          false,
					BaseDN:            "basedn",
					BindDN:            "binddn",
					UserBase:          "user",
					UserObjectClasses: []string{"object"},
					UserFilters:       []string{"filter"},
					Timeout:           time.Second * 30,
					RootCA:            validLDAPRootCA,
				},
			},
			res: res{
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
		{
			name: "change ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							instance.NewLDAPIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
								"id1",
								"name",
								[]string{"server"},
								false,
								"basedn",
								"binddn",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("password"),
								},
								"user",
								[]string{"object"},
								[]string{"filter"},
								time.Second*30,
								nil,
								idp.LDAPAttributes{},
								idp.Options{},
							)),
					),
					expectPush(
						func() eventstore.Command {
							t := true
							event, _ := instance.NewLDAPIDPChangedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
								"id1",
								[]idp.LDAPIDPChanges{
									idp.ChangeLDAPName("new name"),
									idp.ChangeLDAPServers([]string{"new server"}),
									idp.ChangeLDAPStartTLS(true),
									idp.ChangeLDAPBaseDN("new basedn"),
									idp.ChangeLDAPBindDN("new binddn"),
									idp.ChangeLDAPBindPassword(&crypto.CryptoValue{
										CryptoType: crypto.TypeEncryption,
										Algorithm:  "enc",
										KeyID:      "id",
										Crypted:    []byte("new password"),
									}),
									idp.ChangeLDAPUserBase("new user"),
									idp.ChangeLDAPUserObjectClasses([]string{"new object"}),
									idp.ChangeLDAPUserFilters([]string{"new filter"}),
									idp.ChangeLDAPTimeout(time.Second * 20),
									idp.ChangeLDAPAttributes(idp.LDAPAttributeChanges{
										IDAttribute:                stringPointer("new id"),
										FirstNameAttribute:         stringPointer("new firstName"),
										LastNameAttribute:          stringPointer("new lastName"),
										DisplayNameAttribute:       stringPointer("new displayName"),
										NickNameAttribute:          stringPointer("new nickName"),
										PreferredUsernameAttribute: stringPointer("new preferredUsername"),
										EmailAttribute:             stringPointer("new email"),
										EmailVerifiedAttribute:     stringPointer("new emailVerified"),
										PhoneAttribute:             stringPointer("new phone"),
										PhoneVerifiedAttribute:     stringPointer("new phoneVerified"),
										PreferredLanguageAttribute: stringPointer("new preferredLanguage"),
										AvatarURLAttribute:         stringPointer("new avatarURL"),
										ProfileAttribute:           stringPointer("new profile"),
									}),
									idp.ChangeLDAPOptions(idp.OptionChanges{
										IsCreationAllowed: &t,
										IsLinkingAllowed:  &t,
										IsAutoCreation:    &t,
										IsAutoUpdate:      &t,
									}),
									idp.ChangeLDAPRootCA(validLDAPRootCA),
								},
							)
							return event
						}(),
					),
				),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: LDAPProvider{
					Name:              "new name",
					Servers:           []string{"new server"},
					StartTLS:          true,
					BaseDN:            "new basedn",
					BindDN:            "new binddn",
					BindPassword:      "new password",
					UserBase:          "new user",
					UserObjectClasses: []string{"new object"},
					UserFilters:       []string{"new filter"},
					Timeout:           time.Second * 20,
					RootCA:            validLDAPRootCA,
					LDAPAttributes: idp.LDAPAttributes{
						IDAttribute:                "new id",
						FirstNameAttribute:         "new firstName",
						LastNameAttribute:          "new lastName",
						DisplayNameAttribute:       "new displayName",
						NickNameAttribute:          "new nickName",
						PreferredUsernameAttribute: "new preferredUsername",
						EmailAttribute:             "new email",
						EmailVerifiedAttribute:     "new emailVerified",
						PhoneAttribute:             "new phone",
						PhoneVerifiedAttribute:     "new phoneVerified",
						PreferredLanguageAttribute: "new preferredLanguage",
						AvatarURLAttribute:         "new avatarURL",
						ProfileAttribute:           "new profile",
					},
					IDPOptions: idp.Options{
						IsCreationAllowed: true,
						IsLinkingAllowed:  true,
						IsAutoCreation:    true,
						IsAutoUpdate:      true,
					},
				},
			},
			res: res{
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore(t),
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			got, err := c.UpdateInstanceLDAPProvider(tt.args.ctx, tt.args.id, tt.args.provider)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assertObjectDetails(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_AddInstanceAppleIDP(t *testing.T) {
	type fields struct {
		eventstore   func(*testing.T) *eventstore.Eventstore
		idGenerator  id.Generator
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx      context.Context
		provider AppleProvider
	}
	type res struct {
		id   string
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"invalid clientID",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				provider: AppleProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-jkn3w", "Errors.IDP.ClientIDMissing"))
				},
			},
		},
		{
			"invalid teamID",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: AppleProvider{
					ClientID: "clientID",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-Ffg32", "Errors.IDP.TeamIDMissing"))
				},
			},
		},
		{
			"invalid keyID",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: AppleProvider{
					ClientID: "clientID",
					TeamID:   "teamID",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-GDjm5", "Errors.IDP.KeyIDMissing"))
				},
			},
		},
		{
			"invalid privateKey",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: AppleProvider{
					ClientID: "clientID",
					TeamID:   "teamID",
					KeyID:    "keyID",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-GVD4n", "Errors.IDP.PrivateKeyMissing"))
				},
			},
		},
		{
			name: "ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectPush(
						instance.NewAppleIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
							"id1",
							"",
							"clientID",
							"teamID",
							"keyID",
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("privateKey"),
							},
							nil,
							idp.Options{},
						),
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: AppleProvider{
					ClientID:   "clientID",
					TeamID:     "teamID",
					KeyID:      "keyID",
					PrivateKey: []byte("privateKey"),
				},
			},
			res: res{
				id:   "id1",
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
		{
			name: "ok all set",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectPush(
						instance.NewAppleIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
							"id1",
							"",
							"clientID",
							"teamID",
							"keyID",
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("privateKey"),
							},
							[]string{"name", "email"},
							idp.Options{
								IsCreationAllowed: true,
								IsLinkingAllowed:  true,
								IsAutoCreation:    true,
								IsAutoUpdate:      true,
							},
						),
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: AppleProvider{
					ClientID:   "clientID",
					TeamID:     "teamID",
					KeyID:      "keyID",
					PrivateKey: []byte("privateKey"),
					Scopes:     []string{"name", "email"},
					IDPOptions: idp.Options{
						IsCreationAllowed: true,
						IsLinkingAllowed:  true,
						IsAutoCreation:    true,
						IsAutoUpdate:      true,
					},
				},
			},
			res: res{
				id:   "id1",
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore(t),
				idGenerator:         tt.fields.idGenerator,
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			id, got, err := c.AddInstanceAppleProvider(tt.args.ctx, tt.args.provider)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.id, id)
				assertObjectDetails(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_UpdateInstanceAppleIDP(t *testing.T) {
	type fields struct {
		eventstore   func(*testing.T) *eventstore.Eventstore
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx      context.Context
		id       string
		provider AppleProvider
	}
	type res struct {
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"invalid id",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				provider: AppleProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-FRHBH", "Errors.IDMissing"))
				},
			},
		},
		{
			"invalid clientID",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				id:       "id1",
				provider: AppleProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-SFm4l", "Errors.IDP.ClientIDMissing"))
				},
			},
		},
		{
			"invalid teamID",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: AppleProvider{
					ClientID: "clientID",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-SG34t", "Errors.IDP.TeamIDMissing"))
				},
			},
		},
		{
			"invalid keyID",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: AppleProvider{
					ClientID: "clientID",
					TeamID:   "teamID",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-Gh4z2", "Errors.IDP.KeyIDMissing"))
				},
			},
		},
		{
			name: "not found",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: AppleProvider{
					ClientID: "clientID",
					TeamID:   "teamID",
					KeyID:    "keyID",
				},
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "no changes",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							instance.NewAppleIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
								"id1",
								"",
								"clientID",
								"teamID",
								"keyID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("privateKey"),
								},
								nil,
								idp.Options{},
							)),
					),
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: AppleProvider{
					ClientID: "clientID",
					TeamID:   "teamID",
					KeyID:    "keyID",
				},
			},
			res: res{
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
		{
			name: "change ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							instance.NewAppleIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
								"id1",
								"",
								"clientID",
								"teamID",
								"keyID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("privateKey"),
								},
								nil,
								idp.Options{},
							)),
					),
					expectPush(
						func() eventstore.Command {
							t := true
							event, _ := instance.NewAppleIDPChangedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
								"id1",
								[]idp.AppleIDPChanges{
									idp.ChangeAppleClientID("clientID2"),
									idp.ChangeAppleTeamID("teamID2"),
									idp.ChangeAppleKeyID("keyID2"),
									idp.ChangeApplePrivateKey(&crypto.CryptoValue{
										CryptoType: crypto.TypeEncryption,
										Algorithm:  "enc",
										KeyID:      "id",
										Crypted:    []byte("newPrivateKey"),
									}),
									idp.ChangeAppleScopes([]string{"name", "email"}),
									idp.ChangeAppleOptions(idp.OptionChanges{
										IsCreationAllowed: &t,
										IsLinkingAllowed:  &t,
										IsAutoCreation:    &t,
										IsAutoUpdate:      &t,
									}),
								},
							)
							return event
						}(),
					),
				),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: AppleProvider{
					ClientID:   "clientID2",
					TeamID:     "teamID2",
					KeyID:      "keyID2",
					PrivateKey: []byte("newPrivateKey"),
					Scopes:     []string{"name", "email"},
					IDPOptions: idp.Options{
						IsCreationAllowed: true,
						IsLinkingAllowed:  true,
						IsAutoCreation:    true,
						IsAutoUpdate:      true,
					},
				},
			},
			res: res{
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore(t),
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			got, err := c.UpdateInstanceAppleProvider(tt.args.ctx, tt.args.id, tt.args.provider)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assertObjectDetails(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_AddInstanceSAMLIDP(t *testing.T) {
	type fields struct {
		eventstore                 func(*testing.T) *eventstore.Eventstore
		idGenerator                id.Generator
		secretCrypto               crypto.EncryptionAlgorithm
		certificateAndKeyGenerator func(id string) ([]byte, []byte, error)
	}
	type args struct {
		ctx      context.Context
		provider *SAMLProvider
	}
	type res struct {
		id   string
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"invalid name",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				provider: &SAMLProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-o07zjotgnd", ""))
				},
			},
		},
		{
			"no metadata",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: &SAMLProvider{
					Name: "name",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-3bi3esi16t", "Errors.Invalid.Argument"))
				},
			},
		},
		{
			"invalid metadata, error",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: &SAMLProvider{
					Name:     "name",
					Metadata: []byte("metadata"),
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-SF3rwhgh", "Errors.Project.App.SAMLMetadataFormat"))
				},
			},
		},
		{
			name: "ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectPush(
						instance.NewSAMLIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
							"id1",
							"name",
							validSAMLMetadata,
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("key"),
							},
							[]byte("certificate"),
							"",
							false,
							nil,
							"",
							idp.Options{},
						),
					),
				),
				idGenerator:                id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto:               crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
				certificateAndKeyGenerator: func(id string) ([]byte, []byte, error) { return []byte("key"), []byte("certificate"), nil },
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: &SAMLProvider{
					Name:     "name",
					Metadata: validSAMLMetadata,
				},
			},
			res: res{
				id:   "id1",
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
		{
			name: "ok all set",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectPush(
						instance.NewSAMLIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
							"id1",
							"name",
							validSAMLMetadata,
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("key"),
							},
							[]byte("certificate"),
							"binding",
							true,
							gu.Ptr(domain.SAMLNameIDFormatTransient),
							"customAttribute",
							idp.Options{
								IsCreationAllowed: true,
								IsLinkingAllowed:  true,
								IsAutoCreation:    true,
								IsAutoUpdate:      true,
							},
						),
					),
				),
				idGenerator:                id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto:               crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
				certificateAndKeyGenerator: func(id string) ([]byte, []byte, error) { return []byte("key"), []byte("certificate"), nil },
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: &SAMLProvider{
					Name:                          "name",
					Metadata:                      validSAMLMetadata,
					Binding:                       "binding",
					WithSignedRequest:             true,
					NameIDFormat:                  gu.Ptr(domain.SAMLNameIDFormatTransient),
					TransientMappingAttributeName: "customAttribute",
					IDPOptions: idp.Options{
						IsCreationAllowed: true,
						IsLinkingAllowed:  true,
						IsAutoCreation:    true,
						IsAutoUpdate:      true,
					},
				},
			},
			res: res{
				id:   "id1",
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:                     tt.fields.eventstore(t),
				idGenerator:                    tt.fields.idGenerator,
				idpConfigEncryption:            tt.fields.secretCrypto,
				samlCertificateAndKeyGenerator: tt.fields.certificateAndKeyGenerator,
			}
			id, got, err := c.AddInstanceSAMLProvider(tt.args.ctx, tt.args.provider)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.id, id)
				assertObjectDetails(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_UpdateInstanceGenericSAMLIDP(t *testing.T) {
	type fields struct {
		eventstore   func(*testing.T) *eventstore.Eventstore
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx      context.Context
		id       string
		provider *SAMLProvider
	}
	type res struct {
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"invalid id",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				provider: &SAMLProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-7o3rq1owpm", ""))
				},
			},
		},
		{
			"invalid name",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				id:       "id1",
				provider: &SAMLProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-q2s9rak7o9", ""))
				},
			},
		},
		{
			"no metadata",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: &SAMLProvider{
					Name: "name",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-iw1rxnf4sf", ""))
				},
			},
		},
		{
			"invalid metadata, error",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: &SAMLProvider{
					Name:     "name",
					Metadata: []byte("metadata"),
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-dsfj3kl2", "Errors.Project.App.SAMLMetadataFormat"))
				},
			},
		},
		{
			name: "not found",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: &SAMLProvider{
					Name:     "name",
					Metadata: validSAMLMetadata,
				},
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "no changes",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							instance.NewSAMLIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
								"id1",
								"name",
								validSAMLMetadata,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("key"),
								},
								[]byte("certificate"),
								"",
								false,
								nil,
								"",
								idp.Options{},
							)),
					),
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: &SAMLProvider{
					Name:     "name",
					Metadata: validSAMLMetadata,
				},
			},
			res: res{
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
		{
			name: "change ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							instance.NewSAMLIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
								"id1",
								"name",
								[]byte("metadata"),
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("key"),
								},
								[]byte("certificate"),
								"binding",
								false,
								gu.Ptr(domain.SAMLNameIDFormatUnspecified),
								"",
								idp.Options{},
							)),
					),
					expectPush(
						func() eventstore.Command {
							t := true
							event, _ := instance.NewSAMLIDPChangedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
								"id1",
								[]idp.SAMLIDPChanges{
									idp.ChangeSAMLName("new name"),
									idp.ChangeSAMLMetadata(validSAMLMetadata),
									idp.ChangeSAMLBinding("new binding"),
									idp.ChangeSAMLWithSignedRequest(true),
									idp.ChangeSAMLNameIDFormat(gu.Ptr(domain.SAMLNameIDFormatTransient)),
									idp.ChangeSAMLTransientMappingAttributeName("customAttribute"),
									idp.ChangeSAMLOptions(idp.OptionChanges{
										IsCreationAllowed: &t,
										IsLinkingAllowed:  &t,
										IsAutoCreation:    &t,
										IsAutoUpdate:      &t,
									}),
								},
							)
							return event
						}(),
					),
				),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: &SAMLProvider{
					Name:                          "new name",
					Metadata:                      validSAMLMetadata,
					Binding:                       "new binding",
					WithSignedRequest:             true,
					NameIDFormat:                  gu.Ptr(domain.SAMLNameIDFormatTransient),
					TransientMappingAttributeName: "customAttribute",
					IDPOptions: idp.Options{
						IsCreationAllowed: true,
						IsLinkingAllowed:  true,
						IsAutoCreation:    true,
						IsAutoUpdate:      true,
					},
				},
			},
			res: res{
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore(t),
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			got, err := c.UpdateInstanceSAMLProvider(tt.args.ctx, tt.args.id, tt.args.provider)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assertObjectDetails(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_RegenerateInstanceSAMLProviderCertificate(t *testing.T) {
	type fields struct {
		eventstore                 func(*testing.T) *eventstore.Eventstore
		secretCrypto               crypto.EncryptionAlgorithm
		certificateAndKeyGenerator func(id string) ([]byte, []byte, error)
	}
	type args struct {
		ctx context.Context
		id  string
	}
	type res struct {
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"invalid id",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-7de108gqya", ""))
				},
			},
		},
		{
			name: "not found",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "change ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							instance.NewSAMLIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
								"id1",
								"name",
								[]byte("metadata"),
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("key"),
								},
								[]byte("certificate"),
								"binding",
								false,
								gu.Ptr(domain.SAMLNameIDFormatUnspecified),
								"",
								idp.Options{},
							)),
					),
					expectPush(
						func() eventstore.Command {
							event, _ := instance.NewSAMLIDPChangedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
								"id1",
								[]idp.SAMLIDPChanges{
									idp.ChangeSAMLKey(&crypto.CryptoValue{
										CryptoType: crypto.TypeEncryption,
										Algorithm:  "enc",
										KeyID:      "id",
										Crypted:    []byte("new key"),
									}),
									idp.ChangeSAMLCertificate([]byte("new certificate")),
								},
							)
							return event
						}(),
					),
				),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
				certificateAndKeyGenerator: func(id string) ([]byte, []byte, error) {
					return []byte("new key"), []byte("new certificate"), nil
				},
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
			},
			res: res{
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:                     tt.fields.eventstore(t),
				idpConfigEncryption:            tt.fields.secretCrypto,
				samlCertificateAndKeyGenerator: tt.fields.certificateAndKeyGenerator,
			}
			got, err := c.RegenerateInstanceSAMLProviderCertificate(tt.args.ctx, tt.args.id)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assertObjectDetails(t, tt.res.want, got)
			}
		})
	}
}
