import React, { Fragment, useContext, useEffect, useState } from "react";
import { AuthRequestContext } from "../utils/authrequest";
import { Listbox } from "@headlessui/react";
import { Transition } from "@headlessui/react";
import { ChevronUpDownIcon, CheckIcon } from "@heroicons/react/24/solid";
import clsx from "clsx";
import { Buffer } from "buffer";
import { CopyToClipboard } from "react-copy-to-clipboard";
import BrowserOnly from "@docusaurus/BrowserOnly";

const LinkButton = ({
  instance,
  clientId,
  redirectUri,
  responseType,
  prompt,
  organizationId,
  authMethod,
  codeVerifier,
  scope,
  loginHint,
  idTokenHint,
}) => {
  const [copied, setCopied] = useState(false);

  return (
    <CopyToClipboard
      text={`https://zitadel.com/docs/apis/openidoauth/authrequest?instance=${encodeURIComponent(
        instance
      )}&client_id=${encodeURIComponent(
        clientId
      )}&redirect_uri=${encodeURIComponent(
        redirectUri
      )}&response_type=${encodeURIComponent(
        responseType
      )}&scope=${encodeURIComponent(scope)}&prompt=${encodeURIComponent(
        prompt
      )}&auth_method=${encodeURIComponent(
        authMethod
      )}&code_verifier=${encodeURIComponent(
        codeVerifier
      )}&login_hint=${encodeURIComponent(
        loginHint
      )}&id_token_hint=${encodeURIComponent(
        idTokenHint
      )}&organization_id=${encodeURIComponent(organizationId)}
  `}
      onCopy={() => {
        setCopied(true);
        setTimeout(() => {
          setCopied(false);
        }, 2000);
      }}
    >
      <button className="cursor-pointer border-none h-10 flex flex-row items-center py-2 px-4 text-white bg-gray-500 dark:bg-gray-600 hover:dark:bg-gray-500 hover:text-white rounded-md hover:no-underline font-semibold text-sm plausible-event-name=OIDC+Playground plausible-event-method=Save">
        Copy link
        {copied ? (
          <i className="text-[20px] ml-2 las la-clipboard-check"></i>
        ) : (
          <i className="text-[20px] ml-2 las la-clipboard"></i>
        )}
      </button>
    </CopyToClipboard>
  );
};

export function SetAuthRequest() {
  const {
    instance: [instance, setInstance],
    clientId: [clientId, setClientId],
    redirectUri: [redirectUri, setRedirectUri],
    responseType: [responseType, setResponseType],
    scope: [scope, setScope],
    prompt: [prompt, setPrompt],
    authMethod: [authMethod, setAuthMethod],
    codeVerifier: [codeVerifier, setCodeVerifier],
    codeChallenge: [codeChallenge, setCodeChallenge],
    loginHint: [loginHint, setLoginHint],
    idTokenHint: [idTokenHint, setIdTokenHint],
    organizationId: [organizationId, setOrganizationId],
  } = useContext(AuthRequestContext);

  const inputClasses = (error) =>
    clsx({
      "w-full sm:text-sm h-10 mb-2px rounded-md p-2 bg-input-light-background dark:bg-input-dark-background transition-colors duration-300": true,
      "border border-solid border-input-light-border dark:border-input-dark-border hover:border-black hover:dark:border-white focus:border-primary-light-500 focus:dark:border-primary-dark-500": true,
      "focus:outline-none focus:ring-0 text-base text-black dark:text-white placeholder:italic placeholder-gray-700 dark:placeholder-gray-700": true,
      "border border-warn-light-500 dark:border-warn-dark-500 hover:border-warn-light-500 hover:dark:border-warn-dark-500 focus:border-warn-light-500 focus:dark:border-warn-dark-500":
        error,
    });

  const labelClasses = "text-sm";
  const hintClasses = "mt-1 text-xs text-black/50 dark:text-white/50";

  const allResponseTypes = ["code", "id_token", "id_token token"];

  const allPrompts = ["", "login", "select_account", "create", "none"];

  const allAuthMethods = ["(none) PKCE", "Client Secret Basic"];

  const CodeSnipped = ({ cname, children }) => {
    return <span className={cname}>{children}</span>;
  };

  const allScopes = [
    "openid",
    "email",
    "profile",
    "address",
    "offline_access",
    "urn:zitadel:iam:org:project:id:zitadel:aud",
    "urn:zitadel:iam:user:metadata",
    `urn:zitadel:iam:org:id:${
      organizationId ? organizationId : "[organizationId]"
    }`,
  ];

  const [scopeState, setScopeState] = useState(
    [true, true, true, false, false, false, false, false]
    // new Array(allScopes.length).fill(false)
  );

  function toggleScope(position, forceChecked = false) {
    const updatedCheckedState = scopeState.map((item, index) =>
      index === position ? !item : item
    );

    if (forceChecked) {
      updatedCheckedState[position] = true;
    }

    setScopeState(updatedCheckedState);

    setScope(
      updatedCheckedState
        .map((checked, i) => (checked ? allScopes[i] : ""))
        .filter((s) => !!s)
        .join(" ")
    );
  }

  // Encoding functions for code_challenge

  async function string_to_sha256(message) {
    // encode as UTF-8
    const msgBuffer = new TextEncoder().encode(message);
    // hash the message
    const hashBuffer = await crypto.subtle.digest("SHA-256", msgBuffer);
    // return ArrayBuffer
    return hashBuffer;
  }
  async function encodeCodeChallenge(codeChallenge) {
    let arrayBuffer = await string_to_sha256(codeChallenge);
    let buffer = Buffer.from(arrayBuffer);
    let base64 = buffer.toString("base64");
    let base54url = base64_to_base64url(base64);
    return base54url;
  }
  var base64_to_base64url = function (input) {
    input = input.replace(/\+/g, "-").replace(/\//g, "_").replace(/=+$/g, "");
    return input;
  };

  useEffect(async () => {
    setCodeChallenge(await encodeCodeChallenge(codeVerifier));
  }, [codeVerifier]);

  useEffect(() => {
    const newScopeState = allScopes.map((s) => scope.includes(s));
    if (scopeState !== newScopeState) {
      setScopeState(newScopeState);
    }
  }, [scope]);

  return (
    <div className="bg-white/5 rounded-md p-6 shadow">
      <div className="flex flex-row justify-between">
        <h5 className="text-lg mt-0 mb-4 font-semibold">Your Domain</h5>
        <BrowserOnly>
          {() => (
            <LinkButton
              instance={instance}
              clientId={clientId}
              redirectUri={redirectUri}
              responseType={responseType}
              prompt={prompt}
              scope={scope}
              organizationId={organizationId}
              authMethod={authMethod}
              codeVerifier={codeVerifier}
              loginHint={loginHint}
              idTokenHint={idTokenHint}
            />
          )}
        </BrowserOnly>
      </div>
      <div className="flex flex-col">
        <label className={`${labelClasses} text-yellow-500`}>
          Instance Domain
        </label>
        <input
          className={inputClasses(false)}
          id="instance"
          value={instance}
          onChange={(event) => {
            const value = event.target.value;
            setInstance(value);
          }}
        />
        <span className={hintClasses}>
          The domain of your zitadel instance.
        </span>
      </div>

      <h5 className="text-lg mt-6 mb-2 font-semibold">Required Parameters</h5>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        <div className="flex flex-col">
          <label className={`${labelClasses} text-green-500`}>Client ID</label>
          <input
            className={inputClasses(false)}
            id="client_id"
            value={clientId}
            onChange={(event) => {
              const value = event.target.value;
              setClientId(value);
            }}
          />
          <span className={hintClasses}>
            This is the resource id of an application. It's the application
            where you want your users to login.
          </span>
        </div>

        <div className="flex flex-col">
          <label className={`${labelClasses} text-blue-500`}>
            Redirect URI
          </label>
          <input
            className={inputClasses(false)}
            id="redirect_uri"
            value={redirectUri}
            onChange={(event) => {
              const value = event.target.value;
              setRedirectUri(value);
            }}
          />
          <span className={hintClasses}>
            Must be one of the pre-configured redirect uris for your
            application.
          </span>
        </div>

        <div className="flex flex-col">
          <label className={`${labelClasses} text-orange-500`}>
            ResponseType
          </label>
          <Listbox value={responseType} onChange={setResponseType}>
            <div className="relative">
              <Listbox.Button className="transition-colors duration-300 text-black dark:text-white h-10 relative w-full cursor-default rounded-md bg-white dark:bg-input-dark-background py-2 pl-3 pr-10 text-left focus:outline-none focus-visible:border-indigo-500 focus-visible:ring-2 focus-visible:ring-white focus-visible:ring-opacity-75 focus-visible:ring-offset-2 focus-visible:ring-offset-orange-300 sm:text-sm border border-solid border-input-light-border dark:border-input-dark-border hover:border-black hover:dark:border-white focus:border-primary-light-500 focus:dark:border-primary-dark-500">
                <span className="block truncate">{responseType}</span>
                <span className="pointer-events-none absolute inset-y-0 right-0 flex items-center pr-2">
                  <ChevronUpDownIcon
                    className="h-5 w-5 text-gray-400"
                    aria-hidden="true"
                  />
                </span>
              </Listbox.Button>
              <span className={`${hintClasses} flex`}>
                Determines whether a code, id_token token or just id_token will
                be returned. Most use cases will need code.
              </span>
              <Transition
                as={Fragment}
                leave="transition ease-in duration-100"
                leaveFrom="opacity-100"
                leaveTo="opacity-0"
              >
                <Listbox.Options className="pl-0 list-none z-10 top-10 absolute mt-1 max-h-60 w-full overflow-auto rounded-md bg-white dark:bg-background-dark-300 text-black dark:text-white py-1 text-base ring-1 ring-black ring-opacity-5 focus:outline-none sm:text-sm">
                  {allResponseTypes.map((type, typeIdx) => (
                    <Listbox.Option
                      key={typeIdx}
                      className={({ active }) =>
                        `relative cursor-default select-none py-2 pl-10 pr-4 ${
                          active ? "bg-black/20 dark:bg-white/20" : ""
                        }`
                      }
                      value={type}
                    >
                      {({ selected }) => (
                        <>
                          <span
                            className={`block truncate ${
                              selected ? "font-medium" : "font-normal"
                            }`}
                          >
                            {type}
                          </span>
                          {selected ? (
                            <span className="absolute inset-y-0 left-0 flex items-center pl-3 text-orange-500 dark:text-orange-400">
                              <CheckIcon
                                className="h-5 w-5"
                                aria-hidden="true"
                              />
                            </span>
                          ) : null}
                        </>
                      )}
                    </Listbox.Option>
                  ))}
                </Listbox.Options>
              </Transition>
            </div>
          </Listbox>
        </div>
      </div>

      <div className="grid grid-cols-2 md:grid-cols-2 lg:grid-cols-3 gap-4 mt-6">
        <div className="flex flex-col">
          <label className={`${labelClasses} text-teal-600`}>
            Authentication method
          </label>
          <Listbox value={authMethod} onChange={setAuthMethod}>
            <div className="relative">
              <Listbox.Button className="transition-colors duration-300 text-black dark:text-white h-10 relative w-full cursor-default rounded-md bg-white dark:bg-input-dark-background py-2 pl-3 pr-10 text-left focus:outline-none focus-visible:border-indigo-500 focus-visible:ring-2 focus-visible:ring-white focus-visible:ring-opacity-75 focus-visible:ring-offset-2 focus-visible:ring-offset-orange-300 sm:text-sm border border-solid border-input-light-border dark:border-input-dark-border hover:border-black hover:dark:border-white focus:border-primary-light-500 focus:dark:border-primary-dark-500">
                <span className="block truncate">{authMethod}</span>
                <span className="pointer-events-none absolute inset-y-0 right-0 flex items-center pr-2">
                  <ChevronUpDownIcon
                    className="h-5 w-5 text-gray-400"
                    aria-hidden="true"
                  />
                </span>
              </Listbox.Button>
              <span className={`${hintClasses} flex`}>
                Authentication method
              </span>
              <Transition
                as={Fragment}
                leave="transition ease-in duration-100"
                leaveFrom="opacity-100"
                leaveTo="opacity-0"
              >
                <Listbox.Options className="pl-0 list-none z-10 absolute top-10 mt-1 max-h-60 w-full overflow-auto rounded-md bg-white dark:bg-background-dark-300 text-black dark:text-white py-1 text-base ring-1 ring-black ring-opacity-5 focus:outline-none sm:text-sm">
                  {allAuthMethods.map((type, typeIdx) => (
                    <Listbox.Option
                      key={typeIdx}
                      className={({ active }) =>
                        `h-10 relative cursor-default select-none py-2 pl-10 pr-4 ${
                          active ? "bg-black/20 dark:bg-white/20" : ""
                        }`
                      }
                      value={type}
                    >
                      {({ selected }) => (
                        <>
                          <span
                            className={`block truncate ${
                              selected ? "font-medium" : "font-normal"
                            }`}
                          >
                            {type}
                          </span>
                          {selected ? (
                            <span className="absolute inset-y-0 left-0 flex items-center pl-3 text-cyan-500 dark:text-cyan-400">
                              <CheckIcon
                                className="h-5 w-5"
                                aria-hidden="true"
                              />
                            </span>
                          ) : null}
                        </>
                      )}
                    </Listbox.Option>
                  ))}
                </Listbox.Options>
              </Transition>
            </div>
          </Listbox>
        </div>
        {authMethod === "(none) PKCE" && (
          <div className="flex flex-col">
            <label className={`${labelClasses} text-teal-600`}>
              Code Verifier
            </label>
            <input
              className={inputClasses(false)}
              id="code_verifier"
              value={codeVerifier}
              onChange={(event) => {
                const value = event.target.value;
                setCodeVerifier(value);
              }}
            />
            <span className={hintClasses}>
              <span className="text-teal-600">Authentication method</span> PKCE
              requires a random string used to generate a{" "}
              <code>code_challenge</code>
            </span>
          </div>
        )}
      </div>

      <h5 className="text-lg mt-6 mb-2 font-semibold">Additional Parameters</h5>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        <div>
          <div className="flex flex-col">
            <label className={`${labelClasses} text-cyan-500`}>Prompt</label>
            <Listbox value={prompt} onChange={setPrompt}>
              <div className="relative">
                <Listbox.Button className="transition-colors duration-300 text-black dark:text-white h-10 relative w-full cursor-default rounded-md bg-white dark:bg-input-dark-background py-2 pl-3 pr-10 text-left focus:outline-none focus-visible:border-indigo-500 focus-visible:ring-2 focus-visible:ring-white focus-visible:ring-opacity-75 focus-visible:ring-offset-2 focus-visible:ring-offset-orange-300 sm:text-sm border border-solid border-input-light-border dark:border-input-dark-border hover:border-black hover:dark:border-white focus:border-primary-light-500 focus:dark:border-primary-dark-500">
                  <span className="block truncate">{prompt}</span>
                  <span className="pointer-events-none absolute inset-y-0 right-0 flex items-center pr-2">
                    <ChevronUpDownIcon
                      className="h-5 w-5 text-gray-400"
                      aria-hidden="true"
                    />
                  </span>
                </Listbox.Button>
                <span className={`${hintClasses} flex`}>
                  Define how the user should be prompted on login and register.
                </span>
                <Transition
                  as={Fragment}
                  leave="transition ease-in duration-100"
                  leaveFrom="opacity-100"
                  leaveTo="opacity-0"
                >
                  <Listbox.Options className="pl-0 list-none z-10 absolute top-10 mt-1 max-h-60 w-full overflow-auto rounded-md bg-white dark:bg-background-dark-300 text-black dark:text-white py-1 text-base ring-1 ring-black ring-opacity-5 focus:outline-none sm:text-sm">
                    {allPrompts.map((type, typeIdx) => (
                      <Listbox.Option
                        key={typeIdx}
                        className={({ active }) =>
                          `h-10 relative cursor-default select-none py-2 pl-10 pr-4 ${
                            active ? "bg-black/20 dark:bg-white/20" : ""
                          }`
                        }
                        value={type}
                      >
                        {({ selected }) => (
                          <>
                            <span
                              className={`block truncate ${
                                selected ? "font-medium" : "font-normal"
                              }`}
                            >
                              {type}
                            </span>
                            {selected ? (
                              <span className="absolute inset-y-0 left-0 flex items-center pl-3 text-cyan-500 dark:text-cyan-400">
                                <CheckIcon
                                  className="h-5 w-5"
                                  aria-hidden="true"
                                />
                              </span>
                            ) : null}
                          </>
                        )}
                      </Listbox.Option>
                    ))}
                  </Listbox.Options>
                </Transition>
              </div>
            </Listbox>
          </div>
        </div>

        {prompt === "select_account" && (
          <div className="flex flex-col">
            <label className={`${labelClasses} text-rose-500`}>
              Login hint
            </label>
            <input
              className={inputClasses(false)}
              id="login_hint"
              value={loginHint}
              onChange={(event) => {
                const value = event.target.value;
                setLoginHint(value);
              }}
            />
            <span className={hintClasses}>
              This in combination with a{" "}
              <span className="text-black dark:text-white">select_account</span>{" "}
              <span className="text-cyan-500">prompt</span> the login will
              preselect a user.
            </span>
          </div>
        )}

        {/* <div className="flex flex-col">
          <label className={`${labelClasses} text-blue-500`}>
            ID Token hint
          </label>
          <input
            className={inputClasses(false)}
            id="id_token_hint"
            value={idTokenHint}
            onChange={(event) => {
              const value = event.target.value;
              setIdTokenHint(value);
            }}
          />
          <span className={hintClasses}>
            This in combination with a{" "}
            <span className="text-black dark:text-white">select_account</span>{" "}
            <span className="text-emerald-500">prompt</span> the login will
            preselect a user.
          </span>
        </div> */}
      </div>

      <h5 className="text-lg mt-6 mb-2 font-semibold">Scopes</h5>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-2 gap-4">
        <div className="flex flex-col">
          <label className={`${labelClasses} text-purple-500`}>
            Organization ID
          </label>
          <input
            className={inputClasses(false)}
            id="organization_id"
            value={organizationId}
            onChange={(event) => {
              const value = event.target.value;
              setOrganizationId(value);
              allScopes[7] = `urn:zitadel:iam:org:id:${
                value ? value : "[organizationId]"
              }`;
              toggleScope(8, true);
              setScope(
                scopeState
                  .map((checked, i) => (checked ? allScopes[i] : ""))
                  .filter((s) => !!s)
                  .join(" ")
              );
            }}
          />
          <span className={hintClasses}>
            Enforce organization policies and user membership by requesting the{" "}
            <span className="text-purple-500">scope</span>{" "}
            <code>urn:zitadel:iam:org:id:{organizationId}</code>
          </span>
        </div>
      </div>
      <div className="py-4">
        <p className="text-sm mt-0 mb-0 text-purple-500">Scopes</p>
        <span className={`${hintClasses} flex mb-2`}>
          Request additional information about the user with scopes. The claims
          will be returned on the userinfo_endpoint or in the token (when
          configured).
        </span>
        {allScopes.map((scope, scopeIndex) => {
          return (
            <div key={`scope-${scope}`} className="flex flex-row items-center">
              <input
                type="checkbox"
                id={`scope_${scope}`}
                name="scopes"
                value={`${scope}`}
                checked={scopeState[scopeIndex]}
                onChange={() => {
                  toggleScope(scopeIndex);
                }}
              />
              <label className="ml-4" htmlFor={`scope_${scope}`}>
                {scope}{" "}
                {scopeIndex === 8 && scopeState[8] && !organizationId ? (
                  <strong className="text-red-500">
                    Organization ID missing!
                  </strong>
                ) : null}
              </label>
            </div>
          );
        })}
      </div>

      {/* <h5>Optional Parameters</h5>

      <div className={styles.grid}>
        <div className={styles.inputwrapper}>
          <label className={styles.label}>Id Token Hint</label>
          <input
            className={styles.input}
            id="id_token_hint"
            value={idTokenHint}
            onChange={(event) => {
              const value = event.target.value;
              setIdTokenHint(value);
            }}
          />
        </div>
      </div> */}

      <h5 className="text-lg mt-6 mb-2 font-semibold">
        Your authorization request
      </h5>

      <div className="rounded-md bg-gray-700 shadow dark:bg-black/10 p-2 flex flex-col items-center">
        <code className="text-sm w-full mb-4 bg-transparent border-none">
          <span className="text-yellow-500">
            {instance.endsWith("/") ? instance : instance + "/"}
          </span>
          <span className="text-white">oauth/v2/authorize?</span>
          <CodeSnipped cname="text-green-500">{`client_id=${encodeURIComponent(
            clientId
          )}`}</CodeSnipped>
          <CodeSnipped cname="text-blue-500">{`&redirect_uri=${encodeURIComponent(
            redirectUri
          )}`}</CodeSnipped>
          <CodeSnipped cname="text-orange-500">
            {`&response_type=${encodeURIComponent(responseType)}`}
          </CodeSnipped>
          <CodeSnipped cname="text-purple-500">{`&scope=${encodeURIComponent(
            scope
          )}`}</CodeSnipped>
          {prompt && (
            <CodeSnipped cname="text-cyan-500">{`&prompt=${encodeURIComponent(
              prompt
            )}`}</CodeSnipped>
          )}
          {loginHint && prompt === "select_account" && (
            <CodeSnipped cname="text-rose-500">{`&login_hint=${encodeURIComponent(
              loginHint
            )}`}</CodeSnipped>
          )}
          {authMethod === "(none) PKCE" && (
            <CodeSnipped cname="text-teal-600">{`&code_challenge=${encodeURIComponent(
              codeChallenge
            )}&code_challenge_method=S256`}</CodeSnipped>
          )}
        </code>

        <a
          onClick={() => {
            window.plausible("OIDC Playground", {
              props: { method: "Try it out", pageloc: "Authorize" },
            });
          }}
          target="_blank"
          className="mt-2 flex flex-row items-center py-2 px-4 text-white bg-green-500 dark:bg-green-600 hover:dark:bg-green-500 hover:text-white rounded-md hover:no-underline font-semibold text-sm plausible-event-name=OIDC+Playground plausible-event-method=Try+it+out"
          href={`${
            instance.endsWith("/") ? instance : instance + "/"
          }oauth/v2/authorize?client_id=${encodeURIComponent(
            clientId
          )}&redirect_uri=${encodeURIComponent(
            redirectUri
          )}&response_type=${encodeURIComponent(
            responseType
          )}&scope=${encodeURIComponent(scope)}${
            prompt ? `&prompt=${encodeURIComponent(prompt)}` : ""
          }${
            loginHint && prompt === "select_account"
              ? `&login_hint=${encodeURIComponent(loginHint)}`
              : ""
          }${
            authMethod === "(none) PKCE"
              ? `&code_challenge=${encodeURIComponent(
                  codeChallenge
                )}&code_challenge_method=S256`
              : ""
          }`}
        >
          <span>Try it out</span>
          <i className="text-white text-md ml-2 las la-external-link-alt"></i>
        </a>
      </div>
    </div>
  );
}
