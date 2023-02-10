/* Use this file to export your Markdoc functions */

export const includes = {
  transform(parameters) {
    const [array, value] = Object.values(parameters);

    return Array.isArray(array) ? array.includes(value) : false;
  },
};

export const upper = {
  transform(parameters) {
    const string = parameters[0];
    return typeof string === "string" ? string.toUpperCase() : string;
  },
};

export const instanceDomain = {
  transform(parameters) {
    let instance = parameters[0];
    // setInterval(() => {
    //   console.log("i");
    //   instance = instance + "a";
    // }, 1000);

    // const instance = parameters[0];
    const servicePath = parameters[1];
    const serviceVersion = parameters[2];

    return `https://${instance}/${servicePath}/${serviceVersion}`;
  },
};

export const endpoint = {
  transform(parameters) {
    const method = parameters[0];
    const path = parameters[1];
    const link = parameters[2];

    return `${method} ${path}`;
  },
};

// export const proto = {
//   transform(parameters) {
//     const protoPath = parameters[0];
//     const text = readFileSync(
//       join(__dirname, `../../proto/zitadel/${protoPath}`),
//       "utf8"
//     );

//     console.log(text);

//     const content: string = "";

//     const protoDocument = t.parse(content) as t.ProtoDocument;
//     console.log(JSON.stringify(protoDocument, null, 2));

//     return `${protoPath}`;
//   },
// };
