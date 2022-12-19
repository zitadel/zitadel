export const instanceDomain = {
  transform(parameters) {
    const instance = parameters[0];
    const servicePath = parameters[1];
    const serviceVersion = parameters[2];

    return `https://${instance}/${servicePath}/${serviceVersion}`;
  },
};
