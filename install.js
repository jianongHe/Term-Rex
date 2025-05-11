const https = require('https');
const fs = require('fs');
const path = require('path');
const os = require('os');

const version= '0.1.4';

// Determine platform and architecture
const platform = os.platform();
const arch = os.arch();

// Determine which binary to download
let binaryName;
if (platform === 'win32') {
  binaryName = arch === 'x64' ? `term-rex_${version}_windows_amd64.zip` : null;
} else if (platform === 'darwin') {
  binaryName = arch === 'x64' ? `term-rex_${version}_darwin_amd64.tar.gz` : `term-rex_${version}_darwin_arm64.tar.gz`;
} else if (platform === 'linux') {
  binaryName = arch === 'x64' ? `term-rex_${version}_linux_amd64.tar.gz` : `term-rex_${version}_linux_arm64.tar.gz`;
} else {
  console.error('Unsupported platform:', platform);
  process.exit(1);
}

if (!binaryName) {
  console.error(`Unsupported platform/architecture combination: ${platform}/${arch}`);
  process.exit(1);
}

// Download URL (assuming binaries are uploaded to GitHub Releases)
const downloadUrl = `https://github.com/jianongHe/Term-Rex/releases/download/v${version}/${binaryName}`;

// Local binary path
const binPath = path.join(__dirname, 'bin');
const binaryPath = path.join(binPath, platform === 'win32' ? 'term-rex.exe' : 'term-rex');

// Create bin directory
if (!fs.existsSync(binPath)) {
  fs.mkdirSync(binPath, { recursive: true });
}

// Download binary file
console.log(`Downloading ${binaryName} from ${downloadUrl}...`);

// Create a temporary directory for extraction
const tempDir = path.join(os.tmpdir(), `term-rex-${Math.random().toString(36).substring(7)}`);
fs.mkdirSync(tempDir, { recursive: true });

const tempFile = path.join(tempDir, binaryName);
const file = fs.createWriteStream(tempFile);

https.get(downloadUrl, (response) => {
  if (response.statusCode === 302 || response.statusCode === 301) {
    // Handle redirects
    https.get(response.headers.location, (redirectResponse) => {
      redirectResponse.pipe(file);
      handleFileCompletion(file, tempFile, tempDir, binaryPath, binaryName);
    }).on('error', handleError(tempFile));
  } else if (response.statusCode === 200) {
    response.pipe(file);
    handleFileCompletion(file, tempFile, tempDir, binaryPath, binaryName);
  } else {
    console.error(`Failed to download: Server returned status code ${response.statusCode}`);
    console.error(`URL: ${downloadUrl}`);
    fs.unlinkSync(tempFile);
    process.exit(1);
  }
}).on('error', handleError(tempFile));

function handleFileCompletion(file, tempFile, tempDir, binaryPath, binaryName) {
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
            const exeFile = files.find(f => f === 'term-rex.exe');
            if (!exeFile) {
              throw new Error('Could not find executable in zip file');
            }
            
            // Copy to bin directory
            fs.copyFileSync(path.join(tempDir, exeFile), binaryPath);
            console.log(`Installation complete! Binary installed at: ${binaryPath}`);
            console.log('You can now run the game using the term-rex command.');
          })
          .catch(err => {
            console.error('Extraction failed:', err.message);
            process.exit(1);
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
          console.log('Extracted files:', files);
          
          const exeFile = files.find(f => f === 'term-rex');
          if (!exeFile) {
            throw new Error('Could not find executable in tar.gz file');
          }
          
          // Copy to bin directory
          fs.copyFileSync(path.join(tempDir, exeFile), binaryPath);
          fs.chmodSync(binaryPath, '755');
          console.log(`Installation complete! Binary installed at: ${binaryPath}`);
          console.log('You can now run the game using the term-rex command.');
        }).catch(err => {
          console.error('Extraction failed:', err.message);
          console.error('Contents of temp directory:', fs.readdirSync(tempDir));
          process.exit(1);
        });
      }
    });
  });
}

function handleError(tempFile) {
  return (err) => {
    if (fs.existsSync(tempFile)) {
      fs.unlinkSync(tempFile);
    }
    console.error('Download failed:', err.message);
    process.exit(1);
  };
}
