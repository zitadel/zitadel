import { Fragment, useReducer, useState } from "react";
import { Listbox, Transition } from "@headlessui/react";
import { CheckIcon, ChevronUpDownIcon } from "@heroicons/react/20/solid";
import { useRouter } from "next/router";

const protocols = [
  { name: "REST", code: "rest" },
  { name: "GRPC", code: "grpc" },
];

type Props = {
  defaultProtocol: string;
};

export default function ProtocolSwitcher({ defaultProtocol }: Props) {
  const router = useRouter();

  const { protocol } = router.query;

  console.log(protocol);
  const initial = protocol
    ? protocols.find((p) => p.code === protocol) ?? protocols[0]
    : defaultProtocol
    ? protocols.find((p) => p.code === defaultProtocol) ?? protocols[0]
    : protocols[0];

  console.log(initial);
  const [selected, setSelected] = useState(initial);

  function select(value) {
    setSelected(value);
    console.log(value.code);
    router.push({
      query: { protocol: value.code },
    });
  }

  return (
    <div className="w-28">
      <Listbox value={selected} onChange={select}>
        <div className="relative">
          <Listbox.Button className="relative w-full cursor-default rounded-lg bg-white dark:bg-background-dark-400 py-2 pl-3 pr-10 text-left shadow-md focus:outline-none focus-visible:border-indigo-500 focus-visible:ring-2 focus-visible:ring-white focus-visible:ring-opacity-75 focus-visible:ring-offset-2 focus-visible:ring-offset-orange-300 sm:text-sm">
            <span className="block truncate text-xs">{selected.name}</span>
            <span className="pointer-events-none absolute inset-y-0 right-0 flex items-center pr-2">
              <ChevronUpDownIcon
                className="h-5 w-5 text-gray-400"
                aria-hidden="true"
              />
            </span>
          </Listbox.Button>
          <Transition
            as={Fragment}
            leave="transition ease-in duration-100"
            leaveFrom="opacity-100"
            leaveTo="opacity-0"
          >
            <Listbox.Options className="absolute right-0 mt-1 max-h-60 w-fit overflow-auto rounded-md bg-white dark:bg-background-dark-400 dark:text-white py-1 text-base shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none sm:text-sm">
              {protocols.map((protocol, protocolIdx) => (
                <Listbox.Option
                  key={protocolIdx}
                  className={({ active }) =>
                    `relative cursor-default select-none py-2 pl-10 pr-4 ${
                      active
                        ? "bg-amber-100 dark:bg-background-dark-300 text-amber-900 dark:text-white"
                        : "text-gray-900 dark:text-gray-200"
                    }`
                  }
                  value={protocol}
                >
                  {({ selected }) => (
                    <>
                      <span
                        className={`block truncate ${
                          selected ? "font-medium" : "font-normal"
                        }`}
                      >
                        {protocol.name}
                      </span>
                      {selected ? (
                        <span className="absolute inset-y-0 left-0 flex items-center pl-3 text-primary-light-500 dark:text-primary-dark-500">
                          <CheckIcon className="h-5 w-5" aria-hidden="true" />
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
  );
}
