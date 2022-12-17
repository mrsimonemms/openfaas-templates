ARG WATCHDOG_VERSION="0.9.10"
ARG NODE_VERSION=18-alpine
FROM --platform=${TARGETPLATFORM:-linux/amd64} ghcr.io/openfaas/of-watchdog:${WATCHDOG_VERSION} as watchdog

FROM node:${NODE_VERSION}
COPY --from=watchdog /fwatchdog /usr/bin/fwatchdog
WORKDIR /home/node
COPY . .
WORKDIR /home/node/function
WORKDIR /home/node
RUN npm ci --production \
  && chmod 777 /tmp \
  && chown node:node -Rf ./
ENV NPM_CONFIG_LOGLEVEL warn
ENV fprocess="npm run start"
ENV mode="http"
ENV upstream_url="http://127.0.0.1:3000"
ENV exec_timeout="10s"
ENV write_timeout="15s"
ENV read_timeout="15s"
ENV prefix_logs="false"
ENV LOGGER_LEVEL="trace"
ENV LOGGER_REDACT=""
ENV MONGODB_URL=""
ENV SERVER_ENABLE_SHUTDOWN_HOOKS="true"
ENV SERVER_PORT="3000"
USER node
HEALTHCHECK --interval=3s CMD [ -e /tmp/.lock ] || exit 1
CMD ["fwatchdog"]
