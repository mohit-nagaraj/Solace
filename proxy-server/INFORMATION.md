# Go Proxy Server

This contains a simple and efficient reverse proxy server implemented in Go using the Gin framework. The server forwards incoming requests to a specified S3 bucket based on the subdomain of the request.

## Overview
The Go proxy server efficiently handles incoming requests by extracting the subdomain from the request hostname and forwarding the request to a corresponding directory in an S3 bucket. This setup is particularly useful for serving static files from S3 based on subdomains.

## Features
- **Dynamic Routing**: Routes requests to different S3 directories based on the subdomain.
- **Path Adjustment**: Adjusts request paths, defaulting to `/index.html` for root requests.
- **Reverse Proxy**: Forwards requests directly to the S3 bucket without storing any data locally.

## Setup and Running
1. **Move in this folder**:
   ```sh
   cd proxy-server
   ```

2. **Install dependencies**:
   Ensure you have Go installed

3. **Run the server**:
   ```sh
   go run main.go
   ```

4. **Server Running**:
   The server will be running on `http://localhost:8000`.

## Efficiency
- **Direct Proxying**: Forwards requests directly to the S3 bucket, minimizing latency.
- **Stateless Operation**: Does not store request or response data, keeping memory usage low.
- **Low Overhead**: Minimal processing is done to adjust the request URL, ensuring high performance.

### How It Works
1. **Incoming Request**: The server receives a request to `http://myapp.example.com/`.
2. **Subdomain Extraction**: Extracts `myapp` from the request hostname.
3. **Target URL Construction**: Constructs the target URL `https://vercel-clone-outputs.s3.ap-south-1.amazonaws.com/__outputs/myapp`.
4. **Request Forwarding**: Forwards the request to the target URL. If the request path is `/`, it changes it to `/index.html`.
5. **Response Serving**: Sends the response from the S3 bucket back to the client.

## Usage
Simply access your proxy server through the browser or any HTTP client. The server will handle the routing and proxying based on the subdomain in the request URL.

## Contributing
Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.

## License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
