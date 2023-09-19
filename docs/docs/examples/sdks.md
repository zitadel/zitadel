---
title: ZITADEL SDKs
sidebar_label: SDKs
---

On this page you find our official SDKs, links to supporting frameworks and providers, and resources to help with SDKs.
The SDKs wrap either our [gRPC or REST APIs](/docs/apis/introduction) to provide the client with User Authentication and Management for resources.

## ZITADEL SDKs

| Language / Framework | Link Github                                                   | User Authentication | Manage resources | Notes |
|----------------------|---------------------------------------------------------------| --- | --- | --- |
| .NET                 | [zitadel-net](https://github.com/smartive/zitadel-net)        | ✔️ | ✔️ | `community` |
| Elixir               | [zitadel_api](https://github.com/jshmrtn/zitadel_api)         | ✔️ | ✔️ | `community` |
| Go                   | [zitadel-go](https://github.com/zitadel/zitadel-go)           | ❌ | ✔️ | `official` |
| JVM                  | 🚧 [WIP](https://github.com/zitadel/zitadel/discussions/3650) | ❓ | ❓ | TBD |
| Python               | 🚧 [WIP](https://github.com/zitadel/zitadel/issues/3675)      | ❓ | ❓ | TBD |
| NodeJS               | [@zitadel/node](https://www.npmjs.com/package/@zitadel/node)  | ❌ | ✔️ | `community` |
| Dart                 | [zitadel-dart](https://github.com/smartive/zitadel-dart)      | ❌ | ✔️ | `community` |
| Rust                 | [zitadel-rust](https://github.com/smartive/zitadel-rust)      | ✔️ | ✔️ | `community` |

## Missing SDK

Is your language/framework missing? Fear not, you can generate your gRPC API Client with ease.

1. Make sure to install [buf](https://buf.build/docs/installation/)
2. Create a `buf.gen.yaml` and configure the [plugins](https://buf.build/plugins) you need
3. Run `buf generate https://github.com/zitadel/zitadel#format=git,tag=v2.23.1` (change the versions to your needs)

Let us make an example with Ruby. Any other supported language by buf will work as well. Consult the [buf plugin registry](https://buf.build/plugins) for more ideas.

### Example with Ruby

With gRPC we usually need to generate the client stub and the messages/types. This is why we need two plugins.
The plugin `grpc/ruby` generates the client stub and the plugin `protocolbuffers/ruby` takes care of the messages/types.

```yaml
version: v1
plugins:
  - plugin: buf.build/grpc/ruby
    out: gen
  - plugin: buf.build/protocolbuffers/ruby
    out: gen
```

If you now run `buf generate https://github.com/zitadel/zitadel#format=git,tag=v2.23.1` in the folder where your `buf.gen.yaml` is located you should see the folder `gen` appear.

If you run `ls -la gen/zitadel/` you should see something like this:

```bash
ffo@ffo-pc:~/git/zitadel/ruby$ ls -la gen/zitadel/
total 704
drwxr-xr-x 2 ffo ffo   4096 Apr 11 16:49 .
drwxr-xr-x 3 ffo ffo   4096 Apr 11 16:49 ..
-rw-r--r-- 1 ffo ffo   4397 Apr 11 16:49 action_pb.rb
-rw-r--r-- 1 ffo ffo 141097 Apr 11 16:49 admin_pb.rb
-rw-r--r-- 1 ffo ffo  25151 Apr 11 16:49 admin_services_pb.rb
-rw-r--r-- 1 ffo ffo   6537 Apr 11 16:49 app_pb.rb
-rw-r--r-- 1 ffo ffo   1134 Apr 11 16:49 auth_n_key_pb.rb
-rw-r--r-- 1 ffo ffo  32881 Apr 11 16:49 auth_pb.rb
-rw-r--r-- 1 ffo ffo   6896 Apr 11 16:49 auth_services_pb.rb
-rw-r--r-- 1 ffo ffo   1571 Apr 11 16:49 change_pb.rb
-rw-r--r-- 1 ffo ffo   2488 Apr 11 16:49 event_pb.rb
-rw-r--r-- 1 ffo ffo  14782 Apr 11 16:49 idp_pb.rb
-rw-r--r-- 1 ffo ffo   5031 Apr 11 16:49 instance_pb.rb
-rw-r--r-- 1 ffo ffo 223348 Apr 11 16:49 management_pb.rb
-rw-r--r-- 1 ffo ffo  44402 Apr 11 16:49 management_services_pb.rb
-rw-r--r-- 1 ffo ffo   3020 Apr 11 16:49 member_pb.rb
-rw-r--r-- 1 ffo ffo    855 Apr 11 16:49 message_pb.rb
-rw-r--r-- 1 ffo ffo   1445 Apr 11 16:49 metadata_pb.rb
-rw-r--r-- 1 ffo ffo   2370 Apr 11 16:49 object_pb.rb
-rw-r--r-- 1 ffo ffo    621 Apr 11 16:49 options_pb.rb
-rw-r--r-- 1 ffo ffo   4425 Apr 11 16:49 org_pb.rb
-rw-r--r-- 1 ffo ffo   8538 Apr 11 16:49 policy_pb.rb
-rw-r--r-- 1 ffo ffo   8223 Apr 11 16:49 project_pb.rb
-rw-r--r-- 1 ffo ffo   1022 Apr 11 16:49 quota_pb.rb
-rw-r--r-- 1 ffo ffo   5872 Apr 11 16:49 settings_pb.rb
-rw-r--r-- 1 ffo ffo  20985 Apr 11 16:49 system_pb.rb
-rw-r--r-- 1 ffo ffo   4784 Apr 11 16:49 system_services_pb.rb
-rw-r--r-- 1 ffo ffo  28759 Apr 11 16:49 text_pb.rb
-rw-r--r-- 1 ffo ffo  24170 Apr 11 16:49 user_pb.rb
-rw-r--r-- 1 ffo ffo  13568 Apr 11 16:49 v1_pb.rb
```

Import these files into your project to start interacting with ZITADEL's APIs.

## More

While we are not actively maintaining the following projects, it is worth checking out if you're interested in exploring ZITADEL in different programming languages or frameworks.

- [NodeJS passport](https://github.com/buehler/node-passport-zitadel) authentication helper
- [NextAuth Provider for ZITADEL](https://next-auth.js.org/providers/zitadel)

If we do not provide an example, SDK or guide, we strongly recommend using existing authentication libraries for your language or framework instead of building your own.
Certified libraries have undergone rigorous testing and validation to ensure high security and reliability.
There are many recommended libraries available, this saves time and ensures that users' data is well-protected.

You might want to check out the following links to find a good library:

- [awesome-auth](https://github.com/casbin/awesome-auth)
- [OpenID General References](https://openid.net/developers/libraries/)
- [OpenID certified developer tools](https://openid.net/certified-open-id-developer-tools/)