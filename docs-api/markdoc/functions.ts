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
    const instance = parameters[0];
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
