'use strict';

import { connect, credsAuthenticator, headers, Empty } from 'nats.ws';

const UNIQUE_SRC_ID = 'src-server'; // Unique identifier for source server
const UNIQUE_DST_ID = 'dst-server'; // Unique identifier for destination server

function proxyConnections(mqSrcConnection, mqDstConnection, subjects) {
  for (const subject of subjects) {
    const srcSub = mqSrcConnection.subscribe(subject);
    (async (srcSub) => {
      for await (const m of srcSub) {
        const data = new TextDecoder().decode(m.data);

        if (m.headers?.get('source') === UNIQUE_DST_ID) {
          // Skip messages originating from the destination server
          continue;
        }

        const h = headers();
        h.append('source', UNIQUE_SRC_ID);

        console.error('from src', m.subject, data);
        mqDstConnection.publish(m.subject, data, {
          headers: h,
        });
      }
    })(srcSub);

    const dstSub = mqDstConnection.subscribe(subject);
    (async (dstSub) => {
      for await (const m of dstSub) {
        const data = new TextDecoder().decode(m.data);

        if (m.headers?.get('source') === UNIQUE_SRC_ID) {
          // Skip messages originating from the source server
          continue;
        }

        const h = headers();
        h.append('source', UNIQUE_SRC_ID);

        console.error('from dst', m.subject, data);
        mqSrcConnection.publish(m.subject, data, {
          headers: h,
        });
      }
    })(dstSub);
  }
}

let mqSrcConnection;
let mqDstConnection;

export async function proxyMQ(mqSrc, mqDst, subjects) {
  mqSrcConnection = await connect({
    servers: mqSrc.servers,
    waitOnFirstConnect: true,
    authenticator: credsAuthenticator(new TextEncoder().encode(mqSrc.creds)),
  });

  console.error(`mqProxy connected src to ${mqSrcConnection.getServer()}`);

  keepAlive(mqSrcConnection);

  mqDstConnection = await connect({
    servers: mqDst.servers,
    waitOnFirstConnect: true,
    authenticator: credsAuthenticator(new TextEncoder().encode(mqDst.creds)),
  });

  console.error(`mqProxy connected dst to ${mqDstConnection.getServer()}`);

  proxyConnections(mqSrcConnection, mqDstConnection, subjects);
}

export function stopMQProxy() {
  if (mqSrcConnection) {
    mqSrcConnection.close();
    console.error('mqProxy src connection closed');
  }

  if (mqDstConnection) {
    mqDstConnection.close();
    console.error('mqProxy dst connection closed');
  }
}

function keepAlive(nc) {
  const keepAliveIntervalId = setInterval(
    async () => {
      if (nc) {
        await nc
          .request('instance.chrome-nats-proxy.keepalive', Empty, {
            timeout: 1000,
          })
          .then((m) => {
            console.log(`got response: ${sc.decode(m.data)}`);
          })
          .catch((err) => {
            console.log(`problem with request: ${err.message}`);
          });
      } else {
        clearInterval(keepAliveIntervalId);
      }
    },
    // Set the interval to 20 seconds to prevent the service worker from becoming inactive.
    20 * 1000
  );
}
