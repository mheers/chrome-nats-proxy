'use strict';

// With background scripts you can communicate with popup
// and contentScript files.
// For more information on background script,
// See https://developer.chrome.com/extensions/background_pages

console.error('Hello from Background');

import './mqProxy.js';
import './httpProxy.js';

import { stopMQProxy, proxyMQ } from './mqProxy.js';
import { stopHTTPProxy, proxyHTTP } from './httpProxy.js';

function startProxies() {
  chrome.storage.sync.get(['localMQ', 'remoteMQ'], (data) => {
    const { localMQ = {}, remoteMQ = {} } = data;

    // MQ Proxy
    console.error('starting MQ proxy');
    try {
      proxyMQ(
        {
          servers: localMQ.url,
          creds: localMQ.credentials,
        },
        {
          servers: remoteMQ.url,
          creds: remoteMQ.credentials,
        },
        ['instance.clipboard.*']
      );
    } catch (err) {
      console.error(`mqProxy: error connecting to server`);
      console.error(err);
    }

    // HTTP Proxy
    console.error('starting HTTP proxy');
    try {
      proxyHTTP({
        servers: localMQ.url,
        creds: localMQ.credentials,
      });
    } catch (err) {
      console.error(`httpProxy: error connecting to server`);
      console.error(err);
    }
  });
}

function stopProxies() {
  stopMQProxy();
  stopHTTPProxy();
}

// Listen for messages from the popup
chrome.runtime.onMessage.addListener((message, sender, sendResponse) => {
  if (message.type === 'startProxies') {
    startProxies();

    // Keep the message channel open for async response
    return true;
  } else if (message.type === 'stopProxies') {
    stopProxies();
  }
});
