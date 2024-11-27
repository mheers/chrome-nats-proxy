'use strict';

import { connect, StringCodec, credsAuthenticator } from 'nats.ws';

let nc;

export async function proxyHTTP(mqConfig) {
  try {
    nc = await connect({
      servers: mqConfig.servers,
      waitOnFirstConnect: true,
      authenticator: credsAuthenticator(
        new TextEncoder().encode(mqConfig.creds)
      ),
    });

    console.error(`httpProxy connected to ${nc.getServer()}`);

    // create a codec
    const sc = StringCodec();

    // This subscription listens for `http` requests and returns the response
    const httpProxy = nc.subscribe('instance.chrome-nats-proxy.http');

    (async (httpProxy) => {
      console.error(`Listening for ${httpProxy.getSubject()} requests...`);

      for await (const m of httpProxy) {
        try {
          // Decode the request payload (assuming it contains serialized JSON)
          const request = JSON.parse(sc.decode(m.data));

          const { url, method = 'GET', headers = {}, body = null } = request;
          console.error(
            `[httpProxy] Received request: ${method} ${url} ${JSON.stringify(
              headers
            )}`
          );

          // Construct fetch options
          const fetchOptions = {
            method,
            headers,
            body: body ? JSON.stringify(body) : undefined, // Add body only if provided
          };

          // Fetch the response from the target URL
          const response = await fetch(url, fetchOptions);

          // Collect the response headers and body
          const responseHeaders = {};
          response.headers.forEach((value, key) => {
            responseHeaders[key] = value;
          });

          const responseBody = await response.text();

          // Construct a structured response object
          const proxyResponse = {
            status: response.status,
            statusText: response.statusText,
            headers: responseHeaders,
            body: responseBody,
          };

          console.error(
            `[httpProxy] Fetched ${url} with status ${response.status}`
          );

          // Send the response back to the client
          if (m.respond(sc.encode(JSON.stringify(proxyResponse)))) {
            console.error(`[httpProxy] Handled #${httpProxy.getProcessed()}`);
          } else {
            console.error(
              `[httpProxy] #${httpProxy.getProcessed()} ignored - no reply subject`
            );
          }
        } catch (error) {
          console.error(`[httpProxy] Error processing request:`, error);

          // Send an error response back to the client
          const errorResponse = {
            status: 500,
            statusText: 'Internal Server Error',
            body: error.message,
          };

          if (m.respond(sc.encode(JSON.stringify(errorResponse)))) {
            console.error(
              `[httpProxy] Sent error response for #${httpProxy.getProcessed()}`
            );
          }
        }
      }

      console.error(`Subscription ${httpProxy.getSubject()} drained.`);
    })(httpProxy).catch((err) => {
      console.error(`[httpProxy] Subscription error:`, err);
    });
  } catch (error) {
    console.error(`Error connecting to NATS: ${error}`);
  }
}

export function stopHTTPProxy() {
  if (nc) {
    nc.close();
  }
}
