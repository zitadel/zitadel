---
title: ZITADEL with Java Spring Boot
sidebar_label: Java Spring Boot
---

This integration guide shows you how to integrate **ZITADEL** into your Java Spring Boot API. It demonstrates how to secure your API using
OAuth 2 Token Introspection.

At the end of the guide you should have an API with a protected endpoint.

:::info
This documentation references our [example](https://github.com/zitadel/zitadel-java) on GitHub.
You can either create your own application or directly run the example by providing the necessary arguments.
:::

## Set up application

Before we begin developing our API, we need to perform a few configuration steps in the ZITADEL Console.
You'll need to provide some information about your app. We recommend creating a new app to start from scratch. Navigate to your Project, then add a new application at the top of the page.
Select the **API** application type and continue.

![Create app in console](/img/java-spring/api-create.png)

Select Basic Auth for authenticating at the Introspection Endpoint.

![Create app in console](/img/java-spring/api-create-auth.png)

After successful creation of the app, a pop-up will appear displaying the app's client ID. Copy the client ID and secret, as you will need it to configure your Java client.

![Create api key in console](/img/java-spring/api-create-clientid-secret.png)

## Spring Setup

Now that you have configured your web application on the ZITADEL side, you can proceed with the integration of your Spring client.
This guide will reference the [example repository](https://github.com/zitadel/zitadel-java) and explain the necessary steps taken in there.
If your starting from scratch, you can use the Spring Initializer with the [following setup](https://start.spring.io/#!type=maven-project&language=java&platformVersion=3.2.1&packaging=jar&jvmVersion=17&dependencies=web,lombok,oauth2-resource-server) as a base.

### Support class

To be able to take the most out of ZITADELs RBAC, we first need to create a CustomAuthorityOpaqueTokenIntrospector, that will
customize the introspection behavior and map the role claims (`urn:zitadel:iam:org:project:roles`)
into Spring Security `authiorities`, which can be used later on to determine the granted permissions.

So in your application, create a `support/zitadel` package and in there the `CustomAuthorityOpaqueTokenIntrospector.java`:

```java reference
https://github.com/zitadel/zitadel-java/blob/main/api/src/main/java/demo/app/support/zitadel/CustomAuthorityOpaqueTokenIntrospector.java
```

### Application server configuration

As we have now our support class, we can now create and configure the application server itself.

In a new `config` package, create the `WebSecurityConfig.java`.
This class will take care of the authorization by require the calls on `/api/tasks` to be authorized. Any other endpoint will be public by default.
It will also use the just created CustomAuthorityOpaqueTokenIntrospector for the introspection call:

```java reference
https://github.com/zitadel/zitadel-java/blob/main/api/src/main/java/demo/app/config/WebSecurityConfig.java
```

For the authorization (and the server in general) to work, the application needs some configuration, so please provide the following to your `application.yml` (resources folder):

```yaml reference
https://github.com/zitadel/zitadel-java/blob/main/api/src/main/resources/application.yml
```

Note that the `introspection-uri`, `client-id` and `client-secret` are only placeholders. You can either change them in here using the values provided by ZITADEL
or pass them later on as arguments when starting the application.

### Create example API

Create a `api` package with a `ExampleController.java` file with the content below. This will create an API with three endpoints / methods:
- `/api/healthz`: can be called by anyone and always returns `OK`
- `/api/tasks (GET)`: requires authorization and returns the available tasks
- `/api/tasks (POST)`: requires authorization with granted `admin` role and adds the task to the list

If authorization is required, the token must not be expired and the API has to be part of the audience (either client_id or project_id).

For tests we will use a Personal Access Token or the [Java Spring web example](../login/java-spring).

```java reference
https://github.com/zitadel/zitadel-java/blob/main/api/src/main/java/demo/app/api/ExampleController.java
```

## Test API

In case you've created your own application and depending on your development setup you might need to build the application first:

```bash
mvn clean package -DskipTests
```

You will need to provide the `introspection-uri` (your ZITADEL domain> /oauth/v2/introspect), the `client-id` and `client-secret` previously created:

```bash
java \
  -Dspring.security.oauth2.resourceserver.opaquetoken.introspection-uri=<see configuration above> \
  -Dspring.security.oauth2.resourceserver.opaquetoken.client-id=<see configuration above> \
  -Dspring.security.oauth2.resourceserver.opaquetoken.client-secret=<see configuration above> \
  -jar api/target/api-0.0.2-SNAPSHOT.jar
```

This could look like:

```bash
java \
  -Dspring.security.oauth2.resourceserver.opaquetoken.introspection-uri=https://my-domain.zitadel.cloud/oauth/v2/introspect  \
  -Dspring.security.oauth2.resourceserver.opaquetoken.client-id=243861220627644836@example \
  -Dspring.security.oauth2.resourceserver.opaquetoken.client-secret=WJKLF3kfPOi3optkg9vi3jmfjv8oj32nfiäohj!FSC09RWUSR \
  -jar web/target/web-0.0.2-SNAPSHOT.jar
```

### Public endpoint

Now you can call the API by browser or curl. Try the healthz endpoint first:

```bash
curl -i http://localhost:18090/api/healthz
```

it should return something like: 

```
HTTP/1.1 200 
Vary: Origin
Vary: Access-Control-Request-Method
Vary: Access-Control-Request-Headers
X-Content-Type-Options: nosniff
X-XSS-Protection: 0
Cache-Control: no-cache, no-store, max-age=0, must-revalidate
Pragma: no-cache
Expires: 0
X-Frame-Options: DENY
Content-Type: text/plain;charset=UTF-8
Content-Length: 2
Date: Mon, 15 Jan 2024 09:07:21 GMT

OK
```

### Task list

and the task list endpoint:

```bash
curl -i http://localhost:18090/api/tasks
```

it will return:

```
HTTP/1.1 401 
Vary: Origin
Vary: Access-Control-Request-Method
Vary: Access-Control-Request-Headers
WWW-Authenticate: Bearer
X-Content-Type-Options: nosniff
X-XSS-Protection: 0
Cache-Control: no-cache, no-store, max-age=0, must-revalidate
Pragma: no-cache
Expires: 0
X-Frame-Options: DENY
Content-Length: 0
Date: Mon, 15 Jan 2024 09:07:55 GMT
```

Get a valid access_token for the API. You can either achieve this by getting an access token with the project_id in the audience
(e.g. by using the [Spring Boot web example](../login/java-spring)) use a PAT of a service account.

If you provide a valid Bearer Token:

```bash
curl -i -H "Authorization: Bearer ${token}" http://localhost:18090/api/tasks
```

it will return an empty list:
```
HTTP/1.1 200 
Vary: Origin
Vary: Access-Control-Request-Method
Vary: Access-Control-Request-Headers
X-Content-Type-Options: nosniff
X-XSS-Protection: 0
Cache-Control: no-cache, no-store, max-age=0, must-revalidate
Pragma: no-cache
Expires: 0
X-Frame-Options: DENY
Content-Type: application/json
Transfer-Encoding: chunked
Date: Mon, 15 Jan 2024 09:15:10 GMT

[]
```

### Try to add a new task

Let's see what happens if you call the tasks endpoint:

```bash
 curl -i -X POST -H "Authorization: Bearer ${token}" -H "Content-Type: application/json" --data 'my new task' http://localhost:18090/api/tasks
```

it will complain with a permission denied (missing `admin` role):
```
HTTP/1.1 403 
Vary: Origin
Vary: Access-Control-Request-Method
Vary: Access-Control-Request-Headers
WWW-Authenticate: Bearer error="insufficient_scope", error_description="The request requires higher privileges than provided by the access token.", error_uri="https://tools.ietf.org/html/rfc6750#section-3.1"
X-Content-Type-Options: nosniff
X-XSS-Protection: 0
Cache-Control: no-cache, no-store, max-age=0, must-revalidate
Pragma: no-cache
Expires: 0
X-Frame-Options: DENY
Content-Length: 0
Date: Mon, 15 Jan 2024 09:24:39 GMT
```

### Add admin role

So let's create the role and grant it to the user. To do so, go to your project in ZITADEL Console
and create the role by selecting `Roles` in the navigation and then clicking on the `New Role` button.
Finally, create the role as shown below:

![Create project role in console](/img/java-spring/api-project-role.png)

After you have created the role, let's grant it the user, who requested the tasks.
Click on `Authorization` in the navigation and create a new one by selecting the user and the `admin` role.
After successful creation, it should look like:

![Created authorization in console](/img/java-spring/api-project-auth.png)

So you should now be able to add a new task:

```bash
curl -i -X POST -H "Authorization: Bearer ${token}" -H "Content-Type: application/json" --data 'my new task' http://localhost:18090/api/tasks
```

which will report back the successful addition:
```
HTTP/1.1 200 
Vary: Origin
Vary: Access-Control-Request-Method
Vary: Access-Control-Request-Headers
X-Content-Type-Options: nosniff
X-XSS-Protection: 0
Cache-Control: no-cache, no-store, max-age=0, must-revalidate
Pragma: no-cache
Expires: 0
X-Frame-Options: DENY
Content-Type: application/json
Content-Length: 10
Date: Mon, 15 Jan 2024 09:26:11 GMT

task added
```

Let's now retrieve the task list again:

```bash
curl -i -H "Authorization: Bearer ${token}" http://localhost:18090/api/tasks
```

As you can see your new task is listed:
```
HTTP/1.1 200 
Vary: Origin
Vary: Access-Control-Request-Method
Vary: Access-Control-Request-Headers
X-Content-Type-Options: nosniff
X-XSS-Protection: 0
Cache-Control: no-cache, no-store, max-age=0, must-revalidate
Pragma: no-cache
Expires: 0
X-Frame-Options: DENY
Content-Type: application/json
Transfer-Encoding: chunked
Date: Mon, 15 Jan 2024 09:26:48 GMT

["my new task"]
```
