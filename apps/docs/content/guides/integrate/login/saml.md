---
title: Authenticate users with SAML
sidebar_label: SAML
---

SAML stands for Security Assertion Markup Language. It is a standard commonly used for identity federation and single sign-on (SSO). It is one of the original and most popular standards for SSO. Although it is prone to certain security flaws and exploits if not implemented correctly, it remains relevant and widely used. 


## Why use SAML?

Here are some reasons why organizations might choose SAML:

**Legacy systems compatibility**

SAML has been in use since 2002 and is deeply integrated into many legacy systems and enterprise environments. Organizations with existing SAML infrastructure may prefer to continue using it to avoid costly and complex migrations.

**Enterprise use cases**

SAML is often favored in enterprise settings where detailed user attributes and complex authorization requirements are necessary. Its support for rich metadata and customizable assertions makes it suitable for intricate access control scenarios.

**Mature ecosystem**

The SAML ecosystem is mature, with extensive support from a wide range of enterprise applications and identity providers. This broad compatibility ensures that SAML can be used seamlessly across various platforms and services.

## Common SAML terms
- **Service Provider (SP)**: The application the user is trying to sign into.
- **Identity Provider (IdP)**: The centralized point of authentication.
- **SAML Request**: A communication from the SP to the IdP.
- **SAML Response**: A communication from the IdP to the SP, containing assertions about the user.
- **Assertions**: Statements within the SAML response about the user, signed using XML signatures.
- **Assertion Consumer Service (ACS)**: The endpoint at the SP responsible for processing the SAML response.
- **Attributes**: User information within the SAML response.
- **Relay State**: A way for the IdP to remember where a user was before authentication.
- **SAML Trust**: The configuration between the identity provider and the service provider
- **Metadata**: Trust information exchanged between the identity provider and service provider 


## SAML explained

The **Service Provider (SP)** is the application that the user is trying to sign into. When the service provider sends a communication to the identity provider, it is called a **SAML request**. When the **Identity Provider (IdP)** responds or sends a communication to the service provider, it is called a SAML response. 

Within the SAML response, there are multiple statements about the user, known as **assertions**. These assertions are all signed using XML signatures (also known as DSig) and are sent to be processed at the service provider's specific endpoint called an **Assertion Consumer Service (ACS)**. The ACS is responsible for receiving the SAML response from the identity provider, checking the assertions' signatures, and validating the entire document. This is a crucial part of implementing SAML from the service provider's perspective.

Additionally, the SAML response contains other pieces of information about the user, such as their first name, last name, and other profile information, referred to as **attributes**.

Another important concept in SAML is the **relay state**. The relay state allows the identity provider to remember where a user was before authentication. If a user is browsing anonymously through the service provider and triggers authentication, they will be redirected to the identity provider. After validating the user's identity, the identity provider will redirect them back to the service provider's ACS. The relay state makes sure that the user returns to their original location instead of being dropped on a generic home page.

**SAML trust** is the configuration between the identity provider and the service provider, involving shared information such as a signing certificate and an entity ID (also known as an issuer) from the identity provider. This shared information establishes a trust that allows both parties to validate communications, requests, and responses.

**Metadata** is another crucial term in SAML. Metadata allows for self-configuration between the identity provider and the service provider. Instead of manually exchanging certificates, endpoint URLs, and issuer information, metadata enables the sharing of an XML configuration file or URLs to these files. This allows the service provider and identity provider to self-configure based on the information within these configuration files, making the process less manual and more convenient.


## SAML workflow

One important aspect of SAML is that the user can initiate the authentication process in two primary workflows:

**1. IdP-Initiated: The user starts at the IdP, which sends a SAML response to the SP.**

Commonly, with workforce identity providers, there is a centralized user portal where a user can see a list of applications. Clicking on one of these applications initiates the authentication process without the user first going to the service provider. This method is called IdP-initiated (Identity Provider-initiated) authentication. In this flow, the user goes to the identity provider, which automatically kicks off the SAML response, directing them to the desired application.

**2. SP-Initiated: The user starts at the SP, which sends a SAML request to the IdP, followed by a SAML response from the IdP.**

The other method is SP-initiated (Service Provider-initiated) authentication. Here, the user starts at the service provider, which then sends a SAML request to the identity provider. The identity provider processes this request and sends back a SAML response.
In IdP-initiated flows, there is no SAML request, while in SP-initiated flows, the SAML request and response are both involved. The identity provider must understand how to receive a SAML request and create a SAML response. The service provider, particularly its Assertion Consumer Service (ACS), needs to validate the SAML response and potentially generate a SAML request if SP-initiated flow is supported.

It's important to note that not all service providers support both methods. Some only support IdP-initiated, while others only support SP-initiated. The choice between supporting IdP-initiated or SP-initiated authentication depends on the specific requirements of the product and the preferences of the developer implementing SAML. ZITADEL, for instance, supports only the SP-initiated flow. 

## SAML requests and responses

SAML uses XML for both requests and responses. A typical SAML request from an SP to an IdP includes an ID, timestamp, destination URL, and issuer information. The IdP processes this request and returns a SAML response containing user assertions, which the SP validates.

Let's delve into what a SAML request and response look like with the following shortened examples.

**Sample SAML request**

```xml

<?xml version="1.0" encoding="utf-8"?>
<ns0:AuthnRequest
	xmlns:ns0="urn:oasis:names:tc:SAML:2.0:protocol" ID="id-8LjuzBEUQFYFWjL55" Version="2.0" IssueInstant="2024-06-11T04:13:52Z" Destination="https://my-instance-xtzfbc.zitadel.cloud/saml/v2/SSO" AssertionConsumerServiceURL="http://127.0.0.1:5000/acs">
	<ns1:Issuer
		xmlns:ns1="urn:oasis:names:tc:SAML:2.0:assertion" Format="urn:oasis:names:tc:SAML:2.0:nameid-format:entity">https://zitadel-test.sp/metadata
	</ns1:Issuer>
	<ns2:Signature
		xmlns:ns2="http://www.w3.org/2000/09/xmldsig#">
		<ns2:SignedInfo>
			<ns2:CanonicalizationMethod Algorithm="http://www.w3.org/2001/10/xml-exc-c14n#"></ns2:CanonicalizationMethod>
			<ns2:SignatureMethod Algorithm="http://www.w3.org/2000/09/xmldsig#rsa-sha1"></ns2:SignatureMethod>
			<ns2:Reference URI="#id-8LjuzBEUQFYFWjL55">
				<ns2:Transforms>
					<ns2:Transform Algorithm="http://www.w3.org/2000/09/xmldsig#enveloped-signature"></ns2:Transform>
					<ns2:Transform Algorithm="http://www.w3.org/2001/10/xml-exc-c14n#"></ns2:Transform>
				</ns2:Transforms>
				<ns2:DigestMethod Algorithm="http://www.w3.org/2000/09/xmldsig#sha1"></ns2:DigestMethod>
				<ns2:DigestValue>WSz/EQ72RZ0DTh3DBSRCElpITqM=</ns2:DigestValue>
			</ns2:Reference>
		</ns2:SignedInfo>
		<ns2:SignatureValue>HHjsNh0OLj7...</ns2:SignatureValue>
		<ns2:KeyInfo>
			<ns2:X509Data>
				<ns2:X509Certificate>MIIDtzCCAp+gAwIBAgIUfITRQGue...</ns2:X509Certificate>
			</ns2:X509Data>
		</ns2:KeyInfo>
	</ns2:Signature>
</ns0:AuthnRequest>
```

In this SAML request:

- The **ID** (id-8LjuzBEUQFYFWjL55) is generated by the service provider, and the identity provider will respond to this ID in the SAML response.
- **IssueInstant** (2024-05-11T04:13:52Z) is the timestamp for this request, used by the identity provider for verification to ensure the request is within an acceptable time frame.
- **Destination** (https://my-instance-xtzfbc.zitadel.cloud/saml/v2/SSO) points to the identity provider's URL, ensuring the request is sent to the correct recipient.
- **AssertionConsumerServiceURL** (http://127.0.0.1:5000/acs) indicates where the response should be sent after user authentication.
- **Issuer** (https://zitadel-test.sp/metadata) is a predefined string formatted as a URL, matching what the identity provider expects.
- **Signature** ensures the integrity and authenticity of the SAML request. It includes:
	- **SignedInfo** contains details about the canonicalization and signature methods.
  	- **Reference** points to the signed data and includes a transform and digest method.
  	- **SignatureValue** provides the actual digital signature.
  	- **KeyInfo** containing the X.509 certificate, which holds the public key used to verify the signature.

These XMLs are stringified, encoded, and transmitted according to the SAML binding used. For instance, in HTTP Redirect Binding, the request is sent as URL parameters, while in HTTP POST Binding, it is included in the request body. The signature's transmission method also depends on the binding used, ensuring the integrity and authenticity of the SAML message across different communication channels.


**Sample SAML response**

```xml
<?xml version="1.0" encoding="utf-8"?>
<Response
	xmlns="urn:oasis:names:tc:SAML:2.0:protocol" ID="_164ba12b-6711-40e0-8ddb-55aa810f1c92" InResponseTo="id-8LjuzBEUQFYFWjL55" Version="2.0" IssueInstant="2024-06-11T04:17:41Z" Destination="http://127.0.0.1:5000/acs">
	<Issuer
		xmlns="urn:oasis:names:tc:SAML:2.0:assertion" Format="urn:oasis:names:tc:SAML:2.0:nameid-format:entity">https://my-instance-xtzfbc.zitadel.cloud/saml/v2/metadata
	</Issuer>
	<Status>
		<StatusCode Value="urn:oasis:names:tc:SAML:2.0:status:Success"></StatusCode>
	</Status>
	<Assertion
		xmlns="urn:oasis:names:tc:SAML:2.0:assertion" Version="2.0" ID="_6fbdb616-b77f-46af-9554-989c8b89eeda" IssueInstant="2024-06-11T04:17:41Z">
		<Issuer Format="urn:oasis:names:tc:SAML:2.0:nameid-format:entity">https://my-instance-xtzfbc.zitadel.cloud/saml/v2/metadata</Issuer>
		<Signature
			xmlns="http://www.w3.org/2000/09/xmldsig#">
			<SignedInfo>
				<CanonicalizationMethod Algorithm="http://www.w3.org/2001/10/xml-exc-c14n#"></CanonicalizationMethod>
				<SignatureMethod Algorithm="http://www.w3.org/2001/04/xmldsig-more#rsa-sha256"></SignatureMethod>
				<Reference URI="#_6fbdb616-b77f-46af-9554-989c8b89eeda">
					<Transforms>
						<Transform Algorithm="http://www.w3.org/2000/09/xmldsig#enveloped-signature"></Transform>
						<Transform Algorithm="http://www.w3.org/2001/10/xml-exc-c14n#"></Transform>
					</Transforms>
					<DigestMethod Algorithm="http://www.w3.org/2001/04/xmlenc#sha256"></DigestMethod>
					<DigestValue>/b1R9LJJSeNX...</DigestValue>
				</Reference>
			</SignedInfo>
			<SignatureValue>n5GbV4xhkXV...</SignatureValue>
			<KeyInfo>
				<X509Data>
					<X509Certificate>MIIFITCCAwmgAwIBAgIBUTANBgkqh...</X509Certificate>
				</X509Data>
			</KeyInfo>
		</Signature>
		<Subject>
			<NameID Format="urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress">dakshitha.devrel@gmail.com</NameID>
			<SubjectConfirmation Method="urn:oasis:names:tc:SAML:2.0:cm:bearer">
				<SubjectConfirmationData NotOnOrAfter="2024-06-11T04:22:41Z" Recipient="http://127.0.0.1:5000/acs" InResponseTo="id-8LjuzBEUQFYFWjL55"></SubjectConfirmationData>
			</SubjectConfirmation>
		</Subject>
		<Conditions NotBefore="2024-06-11T04:17:41Z" NotOnOrAfter="2024-06-11T04:22:41Z">
			<AudienceRestriction>
				<Audience>https://zitadel-test.sp/metadata</Audience>
			</AudienceRestriction>
		</Conditions>
		<AuthnStatement AuthnInstant="2024-06-11T04:17:41Z" SessionIndex="_6fbdb616-b77f-46af-9554-989c8b89eeda">
			<AuthnContext>
				<AuthnContextClassRef>urn:oasis:names:tc:SAML:2.0:ac:classes:PasswordProtectedTransport</AuthnContextClassRef>
			</AuthnContext>
		</AuthnStatement>
		<AttributeStatement>
			<Attribute Name="Email" NameFormat="urn:oasis:names:tc:SAML:2.0:attrname-format:basic">
				<AttributeValue>tony.stark@gmail.com</AttributeValue>
			</Attribute>
			<Attribute Name="SurName" NameFormat="urn:oasis:names:tc:SAML:2.0:attrname-format:basic">
				<AttributeValue>Stark</AttributeValue>
			</Attribute>
			<Attribute Name="FirstName" NameFormat="urn:oasis:names:tc:SAML:2.0:attrname-format:basic">
				<AttributeValue>Tony</AttributeValue>
			</Attribute>
			<Attribute Name="FullName" NameFormat="urn:oasis:names:tc:SAML:2.0:attrname-format:basic">
				<AttributeValue>Tony Stark</AttributeValue>
			</Attribute>
			<Attribute Name="UserName" NameFormat="urn:oasis:names:tc:SAML:2.0:attrname-format:basic">
				<AttributeValue>tony.stark@gmail.com</AttributeValue>
			</Attribute>
			<Attribute Name="UserID" NameFormat="urn:oasis:names:tc:SAML:2.0:attrname-format:basic">
				<AttributeValue>260242264868201995</AttributeValue>
			</Attribute>
		</AttributeStatement>
	</Assertion>
</Response>
```

In this SAML response:

- The **ID** (_164ba12b-6711-40e0-8ddb-55aa810f1c92) matches the ID from the SAML request (id-8LjuzBEUQFYFWjL55), allowing the service provider to verify it.
- **InResponseTo** references the ID of the SAML request.
- **IssueInstant** (2024-06-11T04:17:41Z) is the timestamp of the response.
- **Destination** (http://127.0.0.1:5000/acs) confirms that the response is intended for the service provider's Assertion Consumer Service URL.
- **Issuer** (https://my-instance-xtzfbc.zitadel.cloud/saml/v2/metadata) is the identity provider's unique string, used to verify the response's origin.
- **StatusCode** indicates the status of the authentication process, with a value of "Success." If the status is not "Success," it means there was a problem with the login, such as incorrect credentials or other authentication issues.
- **Signature** ensures the integrity and authenticity of the SAML response. It includes:
	- **SignedInfo** contains details about the canonicalization and signature methods.
	- **Reference** points to the signed data and including transforms and a digest method.
	- **SignatureValue** provides the actual digital signature.
	- **KeyInfo** contains the X.509 certificate, which holds the public key used to verify the signature.
- **Assertion** includes user information, such as email, surname, first name, full name, username, user ID etc.
- **Subject** contains the user's unique identifier.
- **Conditions** specify the response's validity window.
- **AuthnStatement** includes authentication details.
- **AttributeStatement** contains additional user attributes.

## SAML identity brokering

### How SAML identity brokering works

- **Initial Authentication Request**: A user attempts to access a service (SP1) that is protected by an IdP (IdP1).
- **Redirection to IdP1**: The user is redirected to IdP1 for authentication. If IdP1 trusts another IdP (IdP2) for authentication, it will redirect the user to IdP2.
- **Authentication at IdP2**: The user authenticates with IdP2, which generates a SAML assertion containing the user's identity and attributes.
- **Assertion Processing**: ZITADEL (IdP1) processes the SAML assertion from IdP2. ZITADEL then creates an independent SAML assertion based on the information received and its own policies.
- **Response to SP1**: ZITADEL (IdP1) sends this newly created SAML assertion to SP1, completing the authentication process.
- **Access Granted to SP1**:IdP1 then sends a final SAML assertion to SP1, which grants the user access to the requested service.

See [Let Users Login with Preferred Identity Provider](https://zitadel.com/docs/guides/integrate/identity-providers/introduction) for more information.


## Best practices for SAML implementation

Implementing SAML securely involves several best practices:

- Limit XML Parser Features: Disable unnecessary features to prevent XML external entity (XXE) attacks.
- Use Canonicalized XML: Normalize XML to prevent manipulation.
- Validate XML Schema: Ensure only expected XML formats are accepted.
- Validate Signatures: Check all signatures in the SAML response.
- Use SSL Encryption: Protect against interception.
- Validate Parties: Ensure the destination, audience, recipient, and issuer information is correct.
- Enforce Validation Window: Accept responses only within a valid time frame.
- Use Historical Cache: Track and reject duplicate IDs to prevent replay attacks.
- Minimize Buffer Size: Protect against DDoS attacks.

## Alternatives to SAML

SAML has its flaws; it can be complex and cumbersome to implement. The primary alternatives to consider is OpenID Connect (OIDC). Some illustrate this with the following analogy: “SAML is to OpenID Connect as SOAP is to REST." Just as REST was created to address some inherent flaws in SOAP, OpenID Connect was created to address some of the limitations in specifications like SAML. OpenID Connect is flexible, easy to use, widely adopted, and reliably secure.

For example, SAML is not well-suited for desktop and mobile applications due to its reliance on HTTP redirects, cookie-based session management, and complex certificate handling. OIDC, on the other hand, offers a modern, flexible, and simpler solution for authentication and authorization needs across desktop and mobile applications. It uses authorization codes and tokens that are easier to manage and supports custom URL schemes and deep linking, which simplifies the handling of redirects and improves the overall user experience​. Furthermore, OIDC offers more flexibility in managing sessions without relying solely on cookies. 

If a project requires SAML due to specific requirements or existing infrastructure, it should be used. However, for new projects, it is advisable to consider OpenID Connect because it is a more modern standard and is the more popular choice in the industry. 


## Testing SAML scenarios using ZITADEL

To test SAML scenarios with ZITADEL, follow these steps:

1. Integrate a SAML SP with ZITADEL as the IdP:
    - Sign up for a ZITADEL account if you don't already have one. If you are self-hosting ZITADEL, you can skip this step.
    -  Create an Organization and a Project in ZITADEL.
    - Within your project, create a SAML application.
    - Follow this example on how to create a SAML SP and integrate ZITADEL as the SAML IdP: [ZITADEL Python SAML SP Integration](https://github.com/zitadel/python-saml-sp).

2. Integrate ZITADEL with another SAML IdP for identity brokering:

    - Configure an identity provider that supports SAML.
    - Set up the necessary metadata and endpoints. Here are some guides to help with this setup:
        - [Configure Entra ID as a SAML IdP](https://zitadel.com/docs/guides/integrate/identity-providers/azure-ad-saml)
        - [Configure Okta as a SAML IdP](https://zitadel.com/docs/guides/integrate/identity-providers/okta-saml)
        - [Configure MockSAML as a SAML IdP](https://zitadel.com/docs/guides/integrate/identity-providers/mocksaml)

3. Create test users and simulate authentication requests:
    - Create test users in ZITADEL.
    - Simulate authentication requests to verify that the SAML assertions are correctly generated and transmitted.
    - Ensure that the SAML assertions contain the expected attributes and that these attributes are correctly processed by the service provider.

4. Simulate various scenarios:
    - Successful Login: Verify that a valid user can successfully authenticate and access the service.
    - Failed Login: Test scenarios where authentication fails, such as incorrect credentials or disabled accounts.
    - Attribute Mapping: Check that user attributes (e.g., roles, permissions) are correctly mapped and utilized by the service provider.
    - Logout Requests: Test single logout (SLO) to ensure that logging out from one service logs the user out of all connected services.

For more information, refer to [SAML Endpoints in ZITADEL](https://zitadel.com/docs/apis/saml/endpoints). 
