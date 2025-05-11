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
  binaryName = arch === 'x64' ? 'term-rex_0.1.2_windows_amd64.zip' : null;
} else if (platform === 'darwin') {
  binaryName = arch === 'x64' ? 'term-rex_0.1.2_darwin_amd64.tar.gz' : 'term-rex_0.1.2_darwin_arm64.tar.gz';
} else if (platform === 'linux') {
  binaryName = arch === 'x64' ? 'term-rex_0.1.2_linux_amd64.tar.gz' : 'term-rex_0.1.2_linux_arm64.tar.gz';
} else {
  console.error('Unsupported platform:', platform);
  process.exit(1);
}

// Download URL (assuming binaries are uploaded to GitHub Releases)
const downloadUrl = `https://github.com/jianongHe/Term-Rex/releases/download/v0.1.2/${binaryName}`;

// Local binary path
const binPath = path.join(__dirname, 'bin');
const binaryPath = path.join(binPath, platform === 'win32' ? 'term-rex.exe' : 'term-rex');

// Create bin directory
if (!fs.existsSync(binPath)) {
  fs.mkdirSync(binPath);
}

// Download binary file
console.log(`Downloading ${binaryName}...`);

// Create a temporary directory for extraction
const tempDir = path.join(os.tmpdir(), `term-rex-${Math.random().toString(36).substring(7)}`);
fs.mkdirSync(tempDir, { recursive: true });

const tempFile = path.join(tempDir, binaryName);
const file = fs.createWriteStream(tempFile);

https.get(downloadUrl, (response) => {
  response.pipe(file);
  file.on('finish', () => {
    file.close(() => {
      console.log('Download complete');
      
      // Extract the archive
      console.log('Extracting...');
      if (binaryName.endsWith('.zip')) {
        // For Windows, we need to extract the zip file
        const extract = require('extract-zip');
        extract(tempFile, { dir: tempDir })
          .then(() => {
            // Find the executable in the extracted files
            const files = fs.readdirSync(tempDir);
            const exeFile = files.find(f => f.endsWith('.exe'));
            if (!exeFile) {
              throw new Error('Could not find executable in zip file');
            }
            
            // Copy to bin directory
            fs.copyFileSync(path.join(tempDir, exeFile), binaryPath);
            console.log('Installation complete! You can now run the game using the term-rex command.');
          })
          .catch(err => {
            console.error('Extraction failed:', err.message);
          });
      } else {
        // For Unix systems, extract the tar.gz file
        const tar = require('tar');
        tar.extract({
          file: tempFile,
          cwd: tempDir
        }).then(() => {
          // Find the executable in the extracted files
          const files = fs.readdirSync(tempDir);
          const exeFile = files.find(f => f === 'term-rex');
          if (!exeFile) {
            throw new Error('Could not find executable in tar.gz file');
          }
          
          // Copy to bin directory
          fs.copyFileSync(path.join(tempDir, exeFile), binaryPath);
          fs.chmodSync(binaryPath, '755');
          console.log('Installation complete! You can now run the game using the term-rex command.');
        }).catch(err => {
          console.error('Extraction failed:', err.message);
        });
      }
    });
  });
}).on('error', (err) => {
  fs.unlink(tempFile, () => {});
  console.error('Download failed:', err.message);
});
