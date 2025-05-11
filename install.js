const https = require('https');
const fs = require('fs');
const path = require('path');
const os = require('os');

// Determine platform and architecture
const platform = os.platform();
const arch = os.arch();

// Determine which binary to download
let binaryName;
if (platform === 'win32') {
  binaryName = arch === 'x64' ? 'term-rex-windows-amd64.exe' : 'term-rex-windows-arm64.exe';
} else if (platform === 'darwin') {
  binaryName = arch === 'x64' ? 'term-rex-macos-amd64' : 'term-rex-macos-arm64';
} else if (platform === 'linux') {
  binaryName = arch === 'x64' ? 'term-rex-linux-amd64' : 'term-rex-linux-arm64';
} else {
  console.error('Unsupported platform:', platform);
  process.exit(1);
}

// Download URL (assuming binaries are uploaded to GitHub Releases)
const downloadUrl = `https://github.com/jianongHe/Term-Rex/releases/download/v1.0.0/${binaryName}`;

// Local binary path
const binPath = path.join(__dirname, 'bin');
const binaryPath = path.join(binPath, platform === 'win32' ? 'term-rex.exe' : 'term-rex');

// Create bin directory
if (!fs.existsSync(binPath)) {
  fs.mkdirSync(binPath);
}

// Download binary file
console.log(`Downloading ${binaryName}...`);
const file = fs.createWriteStream(binaryPath);
https.get(downloadUrl, (response) => {
  response.pipe(file);
  file.on('finish', () => {
    file.close(() => {
      console.log('Download complete');
      // Set execute permissions
      if (platform !== 'win32') {
        fs.chmodSync(binaryPath, '755');
      }
      console.log('Installation complete! You can now run the game using the term-rex command.');
    });
  });
}).on('error', (err) => {
  fs.unlink(binaryPath);
  console.error('Download failed:', err.message);
});
