import React, { Fragment, useContext, useEffect, useState } from "react";
import { AuthRequestContext } from "../utils/authrequest";
import styles from "../css/authrequest.module.css";
import CodeBlock from "@theme/CodeBlock";
import { Listbox } from "@headlessui/react";
import { Transition } from "@headlessui/react";
import { ChevronUpDownIcon, CheckIcon } from "@heroicons/react/24/solid";
import clsx from "clsx";

export function SetAuthRequest() {
  const {
    instance: [instance, setInstance],
    clientId: [clientId, setClientId],
    redirectUri: [redirectUri, setRedirectUri],
    responseType: [responseType, setResponseType],
    scope: [scope, setScope],
    prompt: [prompt, setPrompt],
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

  const allPrompts = ["", "login", "select_account", "create"];

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
    "urn:zitadel:iam:user:resourceowner",
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

  useEffect(() => {
    const newScopeState = allScopes.map((s) => scope.includes(s));
    if (scopeState !== newScopeState) {
      setScopeState(newScopeState);
    }
  }, [scope]);

  return (
    <div className="bg-white/5 rounded-md p-6 shadow">
      <h5 className="text-lg mt-0 mb-4 font-semibold">Your Domain</h5>
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
        <span className={hintClasses}>The domain of your zitadel instance</span>
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
                This is the resource id of an application. It's the application
                where you want your users to login.
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

      <h5 className="text-lg mt-6 mb-2 font-semibold">Additional Parameters</h5>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        <div>
          <div className="flex flex-col">
            <label className={`${labelClasses} text-emerald-500`}>Prompt</label>
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
                  Define if and how the user should be prompted on login.
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
                              <span className="absolute inset-y-0 left-0 flex items-center pl-3 text-emerald-500 dark:text-emerald-400">
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
              <span className="text-emerald-500">prompt</span> the login will
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
              allScopes[8] = `urn:zitadel:iam:org:id:${
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
            When requesting the <span className="text-purple-500">scope</span>{" "}
            <code>urn:zitadel:iam:org:id:{organizationId}`</code>, ZITADEL will
            enforce that the user is a member of the selected organization
          </span>
        </div>
      </div>

      <div className="py-4">
        <p className="text-sm mt-0 mb-0 text-purple-500">Scopes</p>
        <span className={`${hintClasses} flex mb-2`}>
          Request additional information about the user with scopes. The claims
          (results of scopes) will be returned on the userinfo_endpoint or in
          the token (when configured).
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

      <br />

      <div className="rounded-md bg-gray-700 shadow dark:bg-black/10 p-2 flex flex-col items-center">
        <code className="text-sm w-full mb-4 bg-transparent border-none">
          <span className="text-yellow-500">{instance}</span>
          <span className="text-white">oauth/v2/authorize</span>
          <CodeSnipped cname="text-green-500">{`?client_id=${encodeURIComponent(
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
            <CodeSnipped cname="text-emerald-500">{`&prompt=${encodeURIComponent(
              prompt
            )}`}</CodeSnipped>
          )}
          {loginHint && (
            <CodeSnipped cname="text-rose-500">{`&login_hint=${encodeURIComponent(
              loginHint
            )}`}</CodeSnipped>
          )}
        </code>

        <a
          target="_blank"
          className="flex flex-row items-center py-2 px-4 text-white bg-green-500 dark:bg-green-600 hover:dark:bg-green-500 hover:text-white rounded-md hover:no-underline font-semibold text-sm"
          href={`${instance}oauth/v2/authorize?client_id=${encodeURIComponent(
            clientId
          )}&redirect_uri=${encodeURIComponent(
            redirectUri
          )}&response_type=${encodeURIComponent(
            responseType
          )}&scope=${encodeURIComponent(scope)}${
            prompt ? `&prompt=${encodeURIComponent(prompt)}` : ""
          }${loginHint ? `&login_hint=${encodeURIComponent(loginHint)}` : ""}`}
        >
          <span>Try it out</span>
          <i className="text-md ml-2 las la-external-link-alt"></i>
        </a>
      </div>
    </div>
  );
}
