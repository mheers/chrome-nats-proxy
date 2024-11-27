'use strict';

import './popup.css';

(function () {
  // We will make use of Storage API to get and store `count` value
  // More information on Storage API can we found at
  // https://developer.chrome.com/extensions/storage

  // To get storage access, we have to mention it in `permissions` property of manifest.json file
  // More information on Permissions can we found at
  // https://developer.chrome.com/extensions/declare_permissions

  // Initialize the popup with stored values or defaults
  document.addEventListener('DOMContentLoaded', () => {
    // Get saved configurations from chrome.storage
    chrome.storage.sync.get(['localMQ', 'remoteMQ'], (data) => {
      const { localMQ = {}, remoteMQ = {} } = data;

      document.getElementById('localUrl').value = localMQ.url || '';
      document.getElementById('localCredentials').value =
        localMQ.credentials || '';
      document.getElementById('remoteUrl').value = remoteMQ.url || '';
      document.getElementById('remoteCredentials').value =
        remoteMQ.credentials || '';
    });

    // Attach save functionality
    document
      .getElementById('saveBtn')
      .addEventListener('click', saveConfigurations);

    // Attach start functionality
    document.getElementById('startBtn').addEventListener('click', startProxies);

    // Attach stop functionality
    document.getElementById('stopBtn').addEventListener('click', stopProxies);

    // Attach import functionality
    document
      .getElementById('importBtn')
      .addEventListener('click', importConfigurations);

    // Attach export functionality
    document
      .getElementById('exportBtn')
      .addEventListener('click', exportConfigurations);
  });

  function startProxies() {
    chrome.runtime.sendMessage({ type: 'startProxies' });
  }

  function stopProxies() {
    chrome.runtime.sendMessage({ type: 'stopProxies' });
  }

  // Save configurations to chrome.storage
  function saveConfigurations() {
    const localUrl = document.getElementById('localUrl').value.trim();
    const localCredentials =
      document.getElementById('localCredentials').value.trim() + `\n`;
    const remoteUrl = document.getElementById('remoteUrl').value.trim();
    const remoteCredentials =
      document.getElementById('remoteCredentials').value.trim() + `\n`;

    chrome.storage.sync.set(
      {
        localMQ: { url: localUrl, credentials: localCredentials },
        remoteMQ: { url: remoteUrl, credentials: remoteCredentials },
      },
      () => {
        console.error('MQ configurations saved!');
        alert('MQ configurations saved successfully.');
      }
    );
  }
})();

// Export configurations to a JSON file
function exportConfigurations() {
  chrome.storage.sync.get(['localMQ', 'remoteMQ'], (data) => {
    const blob = new Blob([JSON.stringify(data, null, 2)], {
      type: 'application/json',
    });
    const url = URL.createObjectURL(blob);

    // Create a temporary download link
    const a = document.createElement('a');
    a.href = url;
    a.download = 'mq-configurations.json';
    a.click();

    URL.revokeObjectURL(url); // Clean up
    console.log('MQ configurations exported');
  });
}

// Import configurations from a JSON file
function importConfigurations() {
  const fileInput = document.createElement('input');
  fileInput.type = 'file';
  fileInput.accept = 'application/json';

  fileInput.addEventListener('change', (event) => {
    const file = event.target.files[0];
    if (!file) return;

    const reader = new FileReader();
    reader.onload = (e) => {
      try {
        const config = JSON.parse(e.target.result);
        chrome.storage.sync.set(config, () => {
          alert('MQ configurations imported successfully.');
          console.log('MQ configurations imported');
          location.reload(); // Reload popup to reflect changes
        });
      } catch (error) {
        alert(
          'Failed to import configurations. Please ensure the file is a valid JSON.'
        );
        console.error('Import error:', error);
      }
    };

    reader.readAsText(file);
  });

  fileInput.click(); // Trigger file picker
}
