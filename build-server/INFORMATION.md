# Build Server

This contains a simple and efficient build server setup using Google Cloud Run Jobs. The build server triggers a job to build and deploy code to an S3 bucket when a request is made. This setup is particularly useful for automated deployments and CI/CD pipelines.

## Overview
The build server utilizes Google Cloud Run Jobs to handle build and deployment tasks. When a request is received, a Cloud Run Job is triggered to build the code and deploy the artifacts to an S3 bucket. This ensures a scalable and efficient build process.

## Features
- **Automated Builds**: Automatically triggers build jobs based on requests.
- **Scalability**: Leverages Cloud Run Jobs for scalable build and deployment processes.
- **Integration**: Can be integrated into CI/CD pipelines for continuous deployment.

## Setup and Running
1. **Create a Cloud Run Job**:
   ```sh
   gcloud beta run jobs create JOB_NAME \
     --image IMAGE_URL \
     --region REGION \
     --args "pid=12,repo=https://github.com/your-repo.git"
   ```

2. **Execute the Job**:
   ```sh
   gcloud beta run jobs execute JOB_NAME --region REGION
   ```

3. **Monitor Job Execution**:
   You can monitor the job execution in the Google Cloud Console under the Cloud Run Jobs section.

## Efficiency
- **On-Demand Execution**: Jobs are triggered on-demand, ensuring resources are used only when needed.
- **Parallel Execution**: Supports parallel execution of build tasks, improving build times.
- **Resource Management**: Cloud Run automatically manages the underlying infrastructure, scaling based on the workload.

### How It Works
1. **Request Trigger**: A request triggers the Cloud Run Job.
2. **Build Process**: The job runs a container that performs the build process using `./main image build-image`.
3. **Deployment**: The built artifacts are deployed to the specified S3 bucket.

## Usage
- **Trigger a Job**: Use the provided commands to create and execute a Cloud Run Job.
- **Monitor Progress**: Check the status and logs of the job execution in the Google Cloud Console.

## Contributing
Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.

## License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.