# Api-server with Google Cloud Run Jobs

This repository contains a simple build server implemented in Go. It uses Google Cloud Run Jobs to build and deploy code based on incoming HTTP requests. It triggers a Cloud Run Job with the for build-image with specified overrides.

## Overview
This build server listens for HTTP POST requests to trigger Google Cloud Run Jobs. It supports two endpoints: one for running a job with default parameters and another for running a job with specified overrides.

## Routes

### POST `/run-job-override`
Triggers a Cloud Run Job with overridden container arguments.

#### Request Body
```json
{
    "pid": "123",
    "repo_link": "<ur_url>"
}
```

#### Response
- `200 OK`: Job started successfully with overrides
- `500 Internal Server Error`: Failed to start the job

## Setup and Running
1. **Go into the folder**:
   ```sh
   cd build-server
   ```

2. **Update your configuration**:
    Go to main.go and update the following variables with your project details:
   ```sh
   GOOGLE_CLOUD_PROJECT=your-project-id
   CLOUD_RUN_LOCATION=your-cloud-run-location
   CLOUD_RUN_JOB_NAME=your-job-name
   PORT=9000
   ```

3. **Get Google Cloud Service Account Key**:
   Follow the steps in the next section to obtain the `gcp.json` file and place it in the root of this project directory.

4. **Run the server**:
   ```sh
   go run main.go
   ```

## Obtaining Google Cloud Service Account Key
1. **Go to the Google Cloud Console**:
   Visit the [Google Cloud Console](https://console.cloud.google.com/).

2. **Select Your Project**:
   Ensure you are working in the correct project.

3. **Navigate to the Service Accounts Page**:
   Go to `IAM & Admin` -> `Service Accounts` or use the [direct link](https://console.cloud.google.com/iam-admin/serviceaccounts).

4. **Create a Service Account**:
   - Click `Create Service Account`.
   - Enter a name and description for the service account.
   - Click `Create and Continue`.

5. **Grant Permissions**:
   - Assign roles like `Cloud Run Admin` and `Storage Admin`.
   - Click `Continue`.

6. **Create Key**:
   - Click `Done` on the `Grant users access to this service account` step.
   - Find your service account in the list.
   - Click on the `Actions` (three dots) next to it and select `Manage Keys`.
   - Click `Add Key`, then `Create New Key`.
   - Select `JSON` and click `Create`.
   - A JSON file will be downloaded to your computer.

7. **Place the Key in the Root Directory**:
   Move the downloaded JSON file to the root of your project directory and rename it to `credentials.json`.

## Contributing
Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.

## License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
